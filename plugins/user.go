package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	ADD_USER_DEFALUT_PASSWORD = "Ab888888"
)

var UserPluginActions = make(map[string]Action)

func init() {
	UserPluginActions["addUser"] = new(AddUserAction)
	UserPluginActions["deleteUser"] = new(DeleteUserAction)
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
	Result
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

func (action *AddUserAction) CheckParam(input AddUserInput) error {
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

	return nil
}

func (action *AddUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddUserInputs)
	outputs := AddUserOutputs{}
	runAs := ""
	var finalErr error

	for _, input := range inputs.Inputs {
		output := AddUserOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

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
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			err = fmt.Errorf("parseSaltApiCmdScriptCallResult meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				err = fmt.Errorf("%s", v.Stdout+v.Stderr)
			}
			break
		}
		if err != nil {
			err = fmt.Errorf("parseSaltApiCmdScriptCallResult meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		encryptPassword, err := AesEnPassword(input.Guid, input.Seed, password, DEFALT_CIPHER)
		if err != nil {
			logrus.Errorf("AesEnPassword meet error(%v)", err)
			err = fmt.Errorf("parseSaltApiCmdScriptCallResult meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.Detail = result
		output.Password = encryptPassword
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

type DeleteUserInputs struct {
	Inputs []DeleteUserInput `json:"inputs,omitempty"`
}

type DeleteUserInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Target   string `json:"target,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type DeleteUserOutputs struct {
	Outputs []DeleteUserOutput `json:"outputs,omitempty"`
}

type DeleteUserOutput struct {
	CallBackParameter
	Result
	Detail string `json:"detail,omitempty"`
	Guid   string `json:"guid,omitempty"`
}

type DeleteUserAction struct {
}

func (action *DeleteUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func deleteUserCheckParam(input DeleteUserInput) error {
	if input.Target == "" {
		return errors.New("DeleteUserAction target is empty")
	}

	if input.UserName == "" {
		return errors.New("DeleteUserAction userName is empty")
	}

	return nil
}

func (action *DeleteUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(DeleteUserInputs)
	outputs := DeleteUserOutputs{}
	runAs := ""
	var finalErr error

	for _, input := range inputs.Inputs {
		output := DeleteUserOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := deleteUserCheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		execArg := fmt.Sprintf("--action remove --user %s ", input.UserName)
		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			err = fmt.Errorf("parseSaltApiCmdScriptCallResult meet err=%v", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				err = fmt.Errorf("%s", v.Stdout+v.Stderr)
			}
			break
		}
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.Detail = result
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
