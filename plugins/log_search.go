package plugins

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

//LogActions define
var LogActions = make(map[string]Action)

func init() {
	LogActions["search"] = new(LogSearchAction)
	LogActions["searchdetail"] = new(LogSearchDetailAction)
}

//LogPlugin .
type LogPlugin struct {
}

//GetActionByName .
func (plugin *LogPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := LogActions[actionName]

	if !found {
		return nil, fmt.Errorf("Log plugin,action = %s not found", actionName)
	}

	return action, nil
}

//LogSearchAction .
type LogSearchAction struct { Language string }

//SearchInputs .
type SearchInputs struct {
	Inputs []SearchInput `json:"inputs,omitempty"`
}

//SearchInput .
type SearchInput struct {
	CallBackParameter
	Guid       string `json:"guid,omitempty"`
	KeyWord    string `json:"keyWord,omitempty"`
	LineNumber int    `json:"lineNumber,omitempty"`
}

//SearchOutputs .
type SearchOutputs struct {
	Outputs []SearchOutput `json:"outputs,omitempty"`
}

//SearchOutput .
type SearchOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	FileName string `json:"fileName,omitempty"`
	Line     string `json:"lineNumber,omitempty"`
	Log      string `json:"log,omitempty"`
}

func (action *LogSearchAction) SetAcceptLanguage(language string) {
	action.Language = language
}

//ReadParam .
func (action *LogSearchAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogSearchAction) CheckParam(input SearchInput) error {

	if input.KeyWord == "" {
		return getParamEmptyError(action.Language, "KeyWord")
	}

	return nil
}

//Do .
func (action *LogSearchAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(SearchInputs)

	var logoutputs SearchOutputs
	var finalErr error

	for i := 0; i < len(logs.Inputs); i++ {
		outputs, err := action.Search(&logs.Inputs[i])
		if err == nil {
			logoutputs.Outputs = append(logoutputs.Outputs, outputs.Outputs...)
		}else{
			log.Logger.Error("Log search action", log.Error(err))
		}
		finalErr = err
	}

	return &logoutputs, finalErr
}

//Search .
func (action *LogSearchAction) Search(input *SearchInput) (outputs SearchOutputs, err error) {
	defer func() {
		if err != nil {
			output := SearchOutput{}
			output.Guid = input.Guid
			output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
		}
	}()
	err = action.CheckParam(*input)
	if err != nil {
		return outputs, err
	}

	sh := "cd logs && "

	keystring := []string{}
	if strings.Contains(input.KeyWord, ",") {
		keystring = strings.Split(input.KeyWord, ",")

		sh += "grep -rin '" + keystring[0] + "' *.log"

		for i := 1; i < len(keystring); i++ {
			sh += "|grep '" + keystring[i] + "'"
		}

	} else {
		sh += "grep -rin '" + input.KeyWord + "' *.log"
	}

	cmd := exec.Command("/bin/sh", "-c", sh)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("can not obtain stdout pipe for command,%s", err.Error())
		return outputs, err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		err = fmt.Errorf("command start fail,%s", err.Error())
		return outputs, err
	}

	logOutput, err := LogReadLine(cmd, stdout)
	if err != nil {
		return outputs, err
	}

	//get filename and lineinfo
	if len(logOutput) > 0 {
		for k := 0; k < len(logOutput); k++ {
			var info SearchOutput
			info.CallBackParameter.Parameter = input.CallBackParameter.Parameter
			info.Guid = input.Guid
			info.Result.Code = RESULT_CODE_SUCCESS

			if logOutput[k] == "" {
				continue
			}

			if !strings.Contains(logOutput[k], ":time=") {
				continue
			}

			fileline := strings.Split(logOutput[k], ":time=")

			if fileline[1] == "" {
				continue
			}

			//single log file
			if !strings.Contains(fileline[0], ":") {
				info.FileName = "wecube-plugins-saltstack.log"
				info.Line = fileline[0]
			} else {
				f := strings.Split(fileline[0], ":")
				info.FileName = f[0]
				info.Line = f[1]
			}

			if len(fileline) == 2 {
				info.Log = "time=" + fileline[1]
			}

			if len(fileline) > 2 {
				info.Log = "time="
				for j := 1; j < len(fileline); j++ {
					info.Log += fileline[j]
				}
			}

			outputs.Outputs = append(outputs.Outputs, info)
		}
	}

	return outputs, err
}

//LogSearchDetailAction .
type LogSearchDetailAction struct { Language string }

//SearchDetailInputs .
type SearchDetailInputs struct {
	Inputs []SearchDetailInput `json:"inputs,omitempty"`
}

//SearchDetailInput .
type SearchDetailInput struct {
	CallBackParameter
	Guid            string `json:"guid,omitempty"`
	FileName        string `json:"fileName,omitempty"`
	LineNumber      string `json:"lineNumber,omitempty"`
	RelateLineCount int    `json:"relateLineCount,omitempty"`
}

//SearchDetailOutputs .
type SearchDetailOutputs struct {
	Outputs []SearchDetailOutput `json:"outputs,omitempty"`
}

//SearchDetailOutput .
type SearchDetailOutput struct {
	CallBackParameter
	Result
	Guid       string `json:"guid,omitempty"`
	FileName   string `json:"fileName,omitempty"`
	LineNumber string `json:"lineNumber,omitempty"`
	Logs       string `json:"logs,omitempty"`
}

func (action *LogSearchDetailAction) SetAcceptLanguage(language string) {
	action.Language = language
}

//ReadParam .
func (action *LogSearchDetailAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SearchDetailInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *LogSearchDetailAction) CheckParam(input SearchDetailInput) error {
	if input.FileName == "" {
		return getParamEmptyError(action.Language, "fileName")
	}
	if input.LineNumber == "" {
		return getParamEmptyError(action.Language, "lineNumber")
	}

	return nil
}

//Do .
func (action *LogSearchDetailAction) Do(input interface{}) (interface{}, error) {
	logs, _ := input.(SearchDetailInputs)

	var logoutputs SearchDetailOutputs
	var finalErr error
	for i := 0; i < len(logs.Inputs); i++ {
		output, err := action.SearchDetail(&logs.Inputs[i])
		if err != nil {
			log.Logger.Error("Log search detail action", log.Error(err))
			finalErr = err
		}

		logoutputs.Outputs = append(logoutputs.Outputs, output)
	}

	return &logoutputs, finalErr
}

//SearchDetail .
func (action *LogSearchDetailAction) SearchDetail(input *SearchDetailInput) (output SearchDetailOutput, err error) {
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

	if input.RelateLineCount <= 0 {
		input.RelateLineCount = 10
	}

	startLine, _ := strconv.Atoi(input.LineNumber)
	shellCmd := fmt.Sprintf("cd logs && cat -n %s |sed -n \"%d,%dp\" ", input.FileName, startLine, startLine+input.RelateLineCount)
	contextText, err := runCmd(shellCmd)
	if err != nil {
		return output, err
	}

	output.FileName = input.FileName
	output.LineNumber = input.LineNumber
	output.Logs = contextText

	return output, err
}

//CountLineNumber .
func CountLineNumber(wLine int, rLine string) (string, string) {
	rline, _ := strconv.Atoi(rLine)
	var num int

	var startLineNumber int
	if rline <= wLine {
		startLineNumber = 1
		num = wLine + rline
	} else {
		startLineNumber = rline - wLine
		num = 2*wLine + 1
	}

	line1 := strconv.Itoa(startLineNumber)

	line2 := strconv.Itoa(num)

	return line1, line2
}
