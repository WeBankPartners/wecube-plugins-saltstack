package plugins

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"errors"
)

var PasswordPluginActions = make(map[string]Action)

func init() {
	PasswordPluginActions["encode"] = new(PasswordEncodeAction)
	PasswordPluginActions["decode"] = new(PasswordDecodeAction)
}

type PasswordPlugin struct {
}

func (plugin *PasswordPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := PasswordPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("Password plugin,action = %s not found", actionName)
	}

	return action, nil
}

type PasswordEncodeInputs struct {
	Inputs []PasswordEncodeInput `json:"inputs,omitempty"`
}

type PasswordEncodeInput struct {
	CallBackParameter
	Guid      string `json:"guid,omitempty"`
	Seed      string `json:"seed,omitempty"`
	Password  string `json:"password,omitempty"`
}

type PasswordEncodeOutputs struct {
	Outputs []PasswordEncodeOutput `json:"outputs,omitempty"`
}

type PasswordEncodeOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordDecodeInputs struct {
	Inputs []PasswordDecodeInput `json:"inputs,omitempty"`
}

type PasswordDecodeInput struct {
	CallBackParameter
	Guid      string `json:"guid,omitempty"`
	Seed      string `json:"seed,omitempty"`
	Password  string `json:"password,omitempty"`
}

type PasswordDecodeOutputs struct {
	Outputs []PasswordDecodeOutput `json:"outputs,omitempty"`
}

type PasswordDecodeOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordEncodeAction struct {
	Language string
}

type PasswordDecodeAction struct {
	Language string
}

func (action *PasswordEncodeAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *PasswordEncodeAction) CheckParam(input PasswordEncodeInput) error {
	if input.Guid == "" {
		return errors.New("PasswordEncodeAction guid is empty")
	}
	if input.Seed == "" {
		return errors.New("PasswordEncodeAction seed is empty")
	}
	if input.Password == "" {
		return errors.New("PasswordEncodeAction password is empty")
	}

	return nil
}

func (action *PasswordEncodeAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PasswordEncodeInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *PasswordEncodeAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(PasswordEncodeInputs)
	outputs := PasswordEncodeOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := PasswordEncodeOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		encryptPassword,err := AesEnPassword(input.Guid, input.Seed, input.Password, DEFALT_CIPHER)
		if err != nil {
			logrus.Errorf("AesEncodePassword meet error(%v)", err)
			err = fmt.Errorf("passwordEncode meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		output.Password = encryptPassword
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

func (action *PasswordDecodeAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *PasswordDecodeAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PasswordDecodeInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *PasswordDecodeAction) CheckParam(input PasswordDecodeInput) error {
	if input.Guid == "" {
		return errors.New("PasswordDecodeAction guid is empty")
	}
	if input.Seed == "" {
		return errors.New("PasswordDecodeAction seed is empty")
	}
	if input.Password == "" {
		return errors.New("PasswordDecodeAction password is empty")
	}

	return nil
}

func (action *PasswordDecodeAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(PasswordDecodeInputs)
	outputs := PasswordDecodeOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := PasswordDecodeOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		decodePassword,err := AesDePassword(input.Guid, input.Seed, input.Password)
		if err != nil {
			logrus.Errorf("AesDecodePassword meet error(%v)", err)
			err = fmt.Errorf("passwordDecode meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		output.Password = decodePassword
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}