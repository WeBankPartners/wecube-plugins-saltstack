package plugins

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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
	Target   string `json:"target,omitempty"`
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
	SearchPattern string `json:"pattern,omitempty"`
}

type SearchTextOutputs struct {
	Outputs []SearchTextOutput `json:"outputs"`
}

type SearchResult struct {
	LineNum  int    `json:"lineNum"`
	LineText string `json:"lineText"`
}

type SearchTextOutput struct {
	Host    string         `json:"host"`
	Results []SearchResult `json:"result"`
}

type SearchTextAction struct {
}

func (action *SearchTextAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchTextInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *SearchTextAction) CheckParam(input interface{}) error {
	inputs, ok := input.(SearchTextInputs)
	if !ok {
		return fmt.Errorf(" SearchTextInputs:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.EndPoint == "" {
			return errors.New("endpoint is empty")
		}

		if input.Target == "" {
			return errors.New("target is empty")
		}

		if input.EndPoint == "" {
			return errors.New("endpoint is empty")
		}

		// if input.AccessKey == "" {
		// 	return errors.New("accessKey is empty")
		// }

		// if input.SecretKey == "" {
		// 	return errors.New("secretKey is empty")
		// }

		if input.SearchPattern == "" {
			return errors.New("search pattern is empty")
		}
	}

	return nil
}

func runCmd(shellCommand string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", shellCommand)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		logrus.Errorf("runCmd (%s) meet err=%v,stderr=%v", shellCommand, err, stderr.String())
		return stderr.String(), nil
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
			logrus.Errorf("searchText  get lineNum meet error,lineNum=%s", line[0:index])
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

func (action *SearchTextAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(SearchTextInputs)
	outputs := SearchTextOutputs{}

	for _, input := range inputs.Inputs {
		// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
		fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
		if err != nil {
			return &outputs, err
		}

		results, err := searchText(fileName, input.SearchPattern)
		os.Remove(fileName)
		if err != nil {
			return &outputs, err
		}

		output := SearchTextOutput{
			Results: results,
			Host:    input.Target,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

//get context
type GetContextInputs struct {
	Inputs []GetContextInput `json:"inputs,omitempty"`
}

type GetContextInput struct {
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey string  `json:"accessKey,omitempty"`
	// SecretKey string  `json:"secretKey,omitempty"`
	LineNum int `json:"lineNum,omitempty"`
	Offset  int `json:"offset,omitempty"`
}

type GetContextOutputs struct {
	Outputs []GetContextOutput `json:"outputs"`
}

type GetContextOutput struct {
	ContextText string `json:"context"`
}

type GetContextAction struct {
}

func (action *GetContextAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs GetContextInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetContextAction) CheckParam(input interface{}) error {
	inputs, ok := input.(GetContextInputs)
	if !ok {
		return fmt.Errorf(" SearchTextInputs:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.EndPoint == "" {
			return errors.New("endpoint is empty")
		}

		// if input.AccessKey == "" {
		// 	return errors.New("accessKey is empty")
		// }

		// if input.SecretKey == "" {
		// 	return errors.New("secretKey is empty")
		// }

		if input.LineNum <= 0 {
			return errors.New("invalid lineNum")
		}
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

func (action *GetContextAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(GetContextInputs)
	outputs := GetContextOutputs{}

	for _, input := range inputs.Inputs {
		// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
		fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
		if err != nil {
			return &outputs, err
		}

		contextText, err := getTextContext(fileName, input.LineNum, input.Offset)
		os.Remove(fileName)
		if err != nil {
			return &outputs, err
		}

		output := GetContextOutput{
			ContextText: contextText,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
