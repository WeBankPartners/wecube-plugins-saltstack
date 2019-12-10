package plugins

import (
	"errors"
	"fmt"
)

const (
	ADD_USER_DEFALUT_PASSWORD = "Ab888888"
)

var UserPluginActions = make(map[string]Action)

func init() {
	UserPluginActions["add"] = new(AddUserAction)
	UserPluginActions["remove"] = new(RemoveUserAction)
}

type UserPlugin struct {
}

func (plugin *UserPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := UserPluginActions[actionName]

	if !found {
		return nil, fmt.Errorf("User plugin,action = %s not found", actionName)
	}

	return action, nil
}

type AddUserInputs struct {
	Inputs []AddUserInput `json:"inputs,omitempty"`
}

type AddUserInput struct {
	CallBackParameter
	Guid      string `json:"guid,omitempty"`
	Target    string `json:"target,omitempty"`
	UserName  string `json:"userName,omitempty"`
	UserId    string `json:"userId,omitempty"`
	Password  string `json:"password,omitempty"`
	UserGroup string `json:"userGroup,omitempty"`
	GroupId   string `json:"groupId,omitempty"`
	HomeDir   string `json:"homeDir,omitempty"`
}

type AddUserOutputs struct {
	Outputs []AddUserOutput `json:"outputs,omitempty"`
}

type AddUserOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type AddUserAction struct {
}

func (action *AddUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *AddUserAction) CheckParam(input AddUserInput) error {
	if input.Target == "" {
		return errors.New("AddUserAction target is empty")
	}

	if input.UserName == "" {
		return errors.New("AddUserAction userName is empty")
	}

	// if input.Password == "" {
	// 	return errors.New("AddUserAction password is empty")
	// }

	// if input.UserGroup == "" {
	// 	return errors.New("AddUserAction userGroup is empty")
	// }

	return nil
}

func (action *AddUserAction) addUser(input *AddUserInput) (output AddUserOutput, err error) {
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

	execArg := fmt.Sprintf("--action add --user %s", input.UserName)
	if input.Password != "" {
		execArg += " --password " + input.Password
	} else {
		execArg += " --password " + ADD_USER_DEFALUT_PASSWORD
	}
	if input.UserGroup != "" {
		execArg += " --userId " + input.UserId
	}
	if input.UserId != "" {
		execArg += " --group " + input.UserGroup
	}
	if input.GroupId != "" {
		execArg += " --groupId " + input.GroupId
	}
	if input.HomeDir != "" {
		execArg += " --home " + input.HomeDir
	}

	result, err := executeS3Script("user_manage.sh", input.Target, "", execArg)
	if err != nil {
		return output, err
	}

	saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
	if err != nil {
		output.Detail = fmt.Sprintf("parseSaltApiCmdScriptCallResult meet err=%v", err)
		return output, err
	}

	for _, v := range saltApiResult.Results[0] {
		if v.RetCode != 0 {
			output.Detail = v.Stderr
			err = fmt.Errorf("%s", v.Stdout+v.Stderr)
			return output, err
		}
		break
	}
	output.Detail = result

	return output, err
}

func (action *AddUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddUserInputs)
	outputs := AddUserOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.addUser(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

type RemoveUserInputs struct {
	Inputs []RemoveUserInput `json:"inputs,omitempty"`
}

type RemoveUserInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Target   string `json:"target,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type RemoveUserOutputs struct {
	Outputs []RemoveUserOutput `json:"outputs,omitempty"`
}

type RemoveUserOutput struct {
	CallBackParameter
	Result
	Detail string `json:"detail,omitempty"`
	Guid   string `json:"guid,omitempty"`
}

type RemoveUserAction struct {
}

func (action *RemoveUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RemoveUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *RemoveUserAction) CheckParam(input RemoveUserInput) error {
	if input.Target == "" {
		return errors.New("RemoveUserAction target is empty")
	}

	if input.UserName == "" {
		return errors.New("RemoveUserAction userName is empty")
	}

	return nil
}

func (action *RemoveUserAction) removeUser(input *RemoveUserInput) (output RemoveUserOutput, err error) {
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

	execArg := fmt.Sprintf("--action remove --user %s ", input.UserName)

	result, err := executeS3Script("user_manage.sh", input.Target, "", execArg)
	if err != nil {
		return output, err
	}

	saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
	if err != nil {
		output.Detail = fmt.Sprintf("parseSaltApiCmdScriptCallResult meet err=%v", err)
		return output, err
	}

	for _, v := range saltApiResult.Results[0] {
		if v.RetCode != 0 {
			output.Detail = v.Stderr
			err = fmt.Errorf("%s", v.Stdout+v.Stderr)
			return output, err
		}
		break
	}
	output.Detail = result

	return output, err
}

func (action *RemoveUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RemoveUserInputs)
	outputs := RemoveUserOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output, err := action.removeUser(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
