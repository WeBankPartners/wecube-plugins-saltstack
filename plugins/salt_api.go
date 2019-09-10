package plugins

import (
	"fmt"
	"encoding/json"
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

type SaltApiResults struct {
	Results []map[string]SaltApiResult `json:"return,omitempty"`
}

type SaltApiResult struct {
	Pid     int    `json:"pid,omitempty"`
	RetCode int    `json:"retcode,omitempty"`
	Stderr  string `json:"stderr,omitempty"`
	Stdout  string `json:"stdout,omitempty"`
}

func parseSaltApiCallResult(jsonStr string) (*SaltApiResults, error) {
	result := SaltApiResults{}

	if err := json.Unmarshal([]byte(jsonStr), &result);err != nil {
		return &result,err
	}
	
	if len(result.Results) == 0 {
		return &result,fmt.Errorf("parseSaltApiCallResult,get %d result",len(result.Results))
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
