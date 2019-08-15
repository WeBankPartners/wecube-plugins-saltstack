package plugins

import (
	"errors"
	"fmt"
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

func (action *AddUserAction) CheckParam(input interface{}) error {
	inputs, ok := input.(AddUserInputs)
	if !ok {
		return fmt.Errorf("AddUserAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.Target == "" {
			return errors.New("AddUserAction target is empty")
		}

		if input.UserName == "" {
			return errors.New("AddUserAction userName is empty")
		}

		if input.Password == "" {
			return errors.New("AddUserAction password is empty")
		}

		if input.UserGroup == "" {
			return errors.New("AddUserAction userGroup is empty")
		}
	}

	return nil
}

func (action *AddUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddUserInputs)
	outputs := AddUserOutputs{}
	runAs := ""

	for _, input := range inputs.Inputs {
		execArg := fmt.Sprintf("--action add --user %s --password %s --group %s", input.UserName, input.Password, input.UserGroup)
		if input.UserId != "" {
			execArg += " --userId " + input.UserId
		}

		if input.GroupId != "" {
			execArg += " --groupId " + input.GroupId
		}

		if input.HomeDir != "" {
			execArg += " --home " + input.HomeDir
		}

		result, err := executeScript("user_manage.sh", input.Target, runAs, execArg)
		if err != nil {
			return nil, err
		}

		output := AddUserOutput{
			Detail: result,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

type RemoveUserInputs struct {
	Inputs []RemoveUserInput `json:"inputs,omitempty"`
}

type RemoveUserInput struct {
	Target   string `json:"target,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type RemoveUserOutputs struct {
	Outputs []RemoveUserOutput `json:"outputs,omitempty"`
}

type RemoveUserOutput struct {
	Detail string `json:"detail,omitempty"`
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

func (action *RemoveUserAction) CheckParam(input interface{}) error {
	inputs, ok := input.(RemoveUserInputs)
	if !ok {
		return fmt.Errorf("RemoveUserAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.Target == "" {
			return errors.New("RemoveUserAction target is empty")
		}

		if input.UserName == "" {
			return errors.New("RemoveUserAction userName is empty")
		}
	}

	return nil
}

func (action *RemoveUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RemoveUserInputs)
	outputs := RemoveUserOutputs{}
	runAs := ""

	for _, input := range inputs.Inputs {
		execArg := fmt.Sprintf("--action remove --user %s ", input.UserName)

		result, err := executeScript("user_manage.sh", input.Target, runAs, execArg)
		if err != nil {
			return nil, err
		}

		output := RemoveUserOutput{
			Detail: result,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
