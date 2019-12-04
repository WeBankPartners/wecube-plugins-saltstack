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
	Seed      string `json:"seed,omitempty"`
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
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
	Detail   string `json:"detail,omitempty"`
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

		if input.Guid == "" {
			return errors.New("AddUserAction guid is empty")
		}

		if input.Seed == "" {
			return errors.New("AddUserAction seed is empty")
		}

		// if input.Password == "" {
		// 	return errors.New("AddUserAction password is empty")
		// }

		// if input.UserGroup == "" {
		// 	return errors.New("AddUserAction userGroup is empty")
		// }
	}

	return nil
}

func (action *AddUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddUserInputs)
	outputs := AddUserOutputs{}
	runAs := ""

	for _, input := range inputs.Inputs {
		password := ""
		execArg := fmt.Sprintf("--action add --user %s", input.UserName)
		if input.Password != "" {
			password = input.Password
		} else {
			password = createRandomPassword()
		}
		execArg += " --password " + password

		if input.UserGroup != "" {
			execArg += " --group " + input.UserGroup
		}
		if input.UserId != "" {
			execArg += " --userId " + input.UserId
		}
		if input.GroupId != "" {
			execArg += " --groupId " + input.GroupId
		}
		if input.HomeDir != "" {
			execArg += " --home " + input.HomeDir
		}

		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg)
		if err != nil {
			return nil, err
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			return fmt.Sprintf("parseSaltApiCmdScriptCallResult meet err=%v", err), err
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				return nil, fmt.Errorf("%s", v.Stdout+v.Stderr)
			}
			break
		}

		md5sum := Md5Encode(input.Guid + input.Seed)
		encryptPassword, err := AesEncode(md5sum[0:16], password)
		if err != nil {
			fmt.Printf("AesEncode meet error(%v)\n", err)
			return nil, err
		}

		output := AddUserOutput{
			Detail:   result,
			Guid:     input.Guid,
			Password: encryptPassword,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
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

		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg)
		if err != nil {
			return nil, err
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			return fmt.Sprintf("parseSaltApiCmdScriptCallResult meet err=%v", err), err
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				return v.Stderr, fmt.Errorf("%s", v.Stdout+v.Stderr)
			}
			break
		}

		output := RemoveUserOutput{
			Detail: result,
			Guid:   input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
