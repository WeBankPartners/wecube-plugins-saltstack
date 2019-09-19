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
	Guid     string   `json:"guid,omitempty"`
	Client   string   `json:"client,omitempty"`
	Target   string   `json:"tgt,omitempty"`
	Function string   `json:"fun,omitempty"`
	Args     []string `json:"arg,omitempty"`
}

type SaltApiCallOutputs struct {
	Outputs []SaltApiCallOutput `json:"outputs,omitempty"`
}

type SaltApiCallOutput struct {
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

func (action *SaltApiCallAction) CheckParam(input interface{}) error {
	_, ok := input.(SaltApiCallInputs)
	if !ok {
		return fmt.Errorf("saltApiAction:input type=%T not right", input)
	}

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
	Results []map[string]string `json:"return,omitempty"`
}

func parseSaltApiCmdRunCallResult(jsonStr string) (*SaltApiCmdRunResults, error) {
	result := SaltApiCmdRunResults{}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return &result, err
	}

	if len(result.Results) == 0 {
		return &result, fmt.Errorf("parseSaltApiCmdRunCallResult,get %d result", len(result.Results))
	}

	return &result, nil
}

func (action *SaltApiCallAction) callSaltApiCall(input *SaltApiCallInput) (*SaltApiCallOutput, error) {
	request := SaltApiRequest{}
	request.Client = input.Client
	request.Function = input.Function
	request.Target = input.Target
	request.Args = input.Args

	result, err := CallSaltApi("https://127.0.0.1:8080", request)
	if err != nil {
		return nil, err
	}

	output := SaltApiCallOutput{}
	output.Guid = input.Guid
	output.Detail = result
	return &output, nil
}

func (action *SaltApiCallAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(SaltApiCallInputs)
	outputs := SaltApiCallOutputs{}
	for _, file := range files.Inputs {
		fileOutput, err := action.callSaltApiCall(&file)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *fileOutput)
	}

	logrus.Infof("all salt request = %v have been handled", files)
	return &outputs, nil
}
