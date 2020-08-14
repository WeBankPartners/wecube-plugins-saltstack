package plugins

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

const (
	DEFAULT_GET_TEXT_CONTEXT_LINE_OFFSET = 10
)

var TextProcessorPluginActions = make(map[string]Action)

func init() {
	TextProcessorPluginActions["search"] = new(SearchTextAction)
	TextProcessorPluginActions["getContext"] = new(GetContextAction)
}

type TextProcessorPlugin struct {
}

func (plugin *TextProcessorPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := TextProcessorPluginActions[actionName]

	if !found {
		return nil, fmt.Errorf("TextProcessor plugin,action = %s not found", actionName)
	}

	return action, nil
}

type SearchTextInputs struct {
	Inputs []SearchTextInput `json:"inputs,omitempty"`
}

type SearchTextInput struct {
	CallBackParameter
	Guid          string `json:"guid,omitempty"`
	Target        string `json:"target,omitempty"`
	EndPoint      string `json:"endpoint,omitempty"`
	SearchPattern string `json:"pattern,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
}

type SearchTextOutputs struct {
	Outputs []SearchTextOutput `json:"outputs"`
}

type SearchResult struct {
	LineNum  int    `json:"lineNum"`
	LineText string `json:"lineText"`
}

type SearchTextOutput struct {
	CallBackParameter
	Result
	Guid    string         `json:"guid,omitempty"`
	Host    string         `json:"host,omitempty"`
	Results []SearchResult `json:"result,omitempty"`
}

type SearchTextAction struct {
	Language string
}

func (action *SearchTextAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *SearchTextAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchTextInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *SearchTextAction) CheckParam(input SearchTextInput) error {
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}

	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}

	if input.SearchPattern == "" {
		return getParamEmptyError(action.Language, "search")
	}

	return nil
}

func runCmd(shellCommand string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", shellCommand)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		log.Logger.Error("Run cmd error", log.String("command", shellCommand), log.String("output", stderr.String()), log.Error(err))
		return stderr.String(), err
	}

	return stdout.String(), nil
}

func searchText(fileName string, pattern string) ([]SearchResult, error) {
	results := []SearchResult{}

	shellCmd := "grep -n " + "\"" + pattern + "\" " + fileName
	stdout, err := runCmd(shellCmd)
	if err != nil {
		return results, err
	}

	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		index := strings.IndexAny(line, ":")
		if index == -1 {
			continue
		}

		lineNum, err := strconv.Atoi(line[0:index])
		if err != nil {
			log.Logger.Error("SearchText get line error", log.String("lineNum", line[0:index]))
			continue
		}

		result := SearchResult{
			LineNum:  lineNum,
			LineText: line[index+1:],
		}
		results = append(results, result)
	}

	return results, nil
}

func (action *SearchTextAction) searchText(input *SearchTextInput) (output SearchTextOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Host = input.Target

		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}
	fileName, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false, action.Language)
	if err != nil {
		return output, err
	}

	results, err := searchText(fileName, input.SearchPattern)
	os.Remove(fileName)
	if err != nil {
		return output, err
	}
	output.Results = results

	return output, err
}

func (action *SearchTextAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(SearchTextInputs)
	outputs := SearchTextOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.searchText(&input)
		if err != nil {
			log.Logger.Error("Search text action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

//get context
type GetContextInputs struct {
	Inputs []GetContextInput `json:"inputs,omitempty"`
}

type GetContextInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	EndPoint string `json:"endpoint,omitempty"`
	LineNum  int    `json:"lineNum,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	// AccessKey string  `json:"accessKey,omitempty"`
	// SecretKey string  `json:"secretKey,omitempty"`
}

type GetContextOutputs struct {
	Outputs []GetContextOutput `json:"outputs"`
}

type GetContextOutput struct {
	CallBackParameter
	Result
	Guid        string `json:"guid,omitempty"`
	ContextText string `json:"context"`
}

type GetContextAction struct {
	Language string
}

func (action *GetContextAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *GetContextAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs GetContextInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetContextAction) CheckParam(input GetContextInput) error {
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}

	if input.LineNum <= 0 {
		return getParamEmptyError(action.Language, "lineNum")
	}

	return nil
}

// sed -n '1,3p' filename
func getTextContext(fileName string, lineNum int, offset int) (string, error) {
	if offset <= 0 {
		offset = DEFAULT_GET_TEXT_CONTEXT_LINE_OFFSET
	}

	startLine := lineNum - offset
	if startLine <= 0 {
		startLine = 1
	}

	shellCmd := fmt.Sprintf("cat -n %s |sed -n \"%d,%dp\" ", fileName, startLine, lineNum+offset)
	//shellCmd:=fmt.Sprintf("sed -n \"%d,%dp\" %s",startLine,lineNum+offset,fileName)
	return runCmd(shellCmd)
}

func (action *GetContextAction) getContext(input *GetContextInput) (output GetContextOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	fileName, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false, action.Language)
	if err != nil {
		return output, err
	}

	contextText, err := getTextContext(fileName, input.LineNum, input.Offset)
	os.Remove(fileName)
	if err != nil {
		return output, err
	}
	output.ContextText = contextText

	return output, err
}

func (action *GetContextAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(GetContextInputs)
	outputs := GetContextOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.getContext(&input)
		if err != nil {
			if err != nil {
				log.Logger.Error("Get context action", log.Error(err))
			}
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
