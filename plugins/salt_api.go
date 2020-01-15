package plugins

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

var SaltApiActions = make(map[string]Action)

func init() {
	SaltApiActions["call"] = new(SaltApiCallAction)
}

type SaltApiPlugin struct {
}

func (plugin *SaltApiPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := SaltApiActions[actionName]
	if !found {
		return nil, fmt.Errorf("SaltApi plugin,action = %s not found", actionName)
	}

	return action, nil
}

type SaltApiCallInputs struct {
	Inputs []SaltApiCallInput `json:"inputs,omitempty"`
}

type SaltApiCallInput struct {
	CallBackParameter
	Guid     string   `json:"guid,omitempty"`
	Client   string   `json:"client,omitempty"`
	Target   string   `json:"target,omitempty"`
	Function string   `json:"function,omitempty"`
	Args     []string `json:"args,omitempty"`
}

type SaltApiCallOutputs struct {
	Outputs []SaltApiCallOutput `json:"outputs,omitempty"`
}

type SaltApiCallOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type SaltApiCallAction struct {
}

func (action *SaltApiCallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SaltApiCallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SaltApiCallAction) CheckParam(input SaltApiCallInput) error {
	return nil
}

type SaltApiCmdScriptResults struct {
	Results []map[string]SaltApiCmdScriptResult `json:"return,omitempty"`
}

type SaltApiCmdScriptResult struct {
	Pid     int    `json:"pid,omitempty"`
	RetCode int    `json:"retcode,omitempty"`
	Stderr  string `json:"stderr,omitempty"`
	Stdout  string `json:"stdout,omitempty"`
}

func parseSaltApiCmdScriptCallResult(jsonStr string) (*SaltApiCmdScriptResults, error) {
	result := SaltApiCmdScriptResults{}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &result, err
	}

	if len(result.Results) == 0 {
		return &result, fmt.Errorf("parseSaltApiCmdScriptCallResult,get %d result", len(result.Results))
	}

	return &result, nil
}

type SaltApiCmdRunResults struct {
	Results []map[string]SaltApiCmdRunResult `json:"return,omitempty"`
}

type SaltApiCmdRunResult struct {
	Jid       string `json:"jid,omitempty"`
	RetCode   int    `json:"retcode,omitempty"`
	RetDetail string `json:"ret,omitempty"`
}

func parseSaltApiCmdRunCallResult(jsonStr string) (*SaltApiCmdRunResults, error) {
	result := SaltApiCmdRunResults{}

	logrus.Infof("parseSaltApiCmdRunCallResult jsonStr: %++v", jsonStr)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &result, err
	}

	if len(result.Results) == 0 {
		return &result, fmt.Errorf("parseSaltApiCmdRunCallResult,get %d result", len(result.Results))
	}

	return &result, nil
}

func (action *SaltApiCallAction) callSaltApiCall(input *SaltApiCallInput) (output SaltApiCallOutput, err error) {
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

	request := SaltApiRequest{}
	request.Client = input.Client
	request.Function = input.Function
	request.Target = input.Target
	request.Args = input.Args

	result, err := CallSaltApi("https://127.0.0.1:8080", request)
	if err != nil {
		return output, err
	}
	output.Detail = result

	return output, err
}

func (action *SaltApiCallAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(SaltApiCallInputs)
	outputs := SaltApiCallOutputs{}
	var finalErr error
	for _, file := range files.Inputs {
		fileOutput, err := action.callSaltApiCall(&file)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, fileOutput)
	}

	logrus.Infof("all salt request = %v have been handled", files)
	return &outputs, finalErr
}
