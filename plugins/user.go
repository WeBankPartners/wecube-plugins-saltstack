package plugins

import (
	"fmt"

	"strings"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

const (
	ADD_USER_DEFALUT_PASSWORD = "Ab888888"
)

var UserPluginActions = make(map[string]Action)

func init() {
	UserPluginActions["add"] = new(AddUserAction)
	UserPluginActions["delete"] = new(DeleteUserAction)
	UserPluginActions["password"] = new(ChangeUserPasswordAction)
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
	RwDir  string `json:"rwDir,omitempty"`
	RwFile    string `json:"rwFile,omitempty"`
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
	Language string
}

func (action *AddUserAction) SetAcceptLanguage(language string) {
	action.Language = language
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
		return getParamEmptyError(action.Language, "target")
	}

	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}

	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}

	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}

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

		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		if strings.Contains(input.UserName, ":") {
			input.UserName = strings.Split(input.UserName, ":")[0]
			if input.UserGroup == "" {
				input.UserGroup = strings.Split(input.UserName, ":")[1]
			}
		}
		password := ""
		execArg := fmt.Sprintf("--action add --user '%s'", input.UserName)
		if input.Password != "" {
			password = input.Password
		} else {
			password = createRandomPassword()
		}
		execArg += " --password '" + password + "'"

		if input.UserGroup != "" {
			execArg += " --group '" + input.UserGroup + "'"
		}
		if input.UserId != "" {
			execArg += " --userId '" + input.UserId + "'"
		}
		if input.GroupId != "" {
			execArg += " --groupId '" + input.GroupId + "'"
		}
		if input.HomeDir != "" {
			execArg += " --home '" + input.HomeDir + "'"
		}
		if input.RwDir != "" {
			input.RwDir = strings.Replace(input.RwDir, "[", "", -1)
			input.RwDir = strings.Replace(input.RwDir, "]", "", -1)
			input.RwDir = strings.Replace(input.RwDir, "&", "", -1)
			execArg += " --makeDir '" + input.RwDir + "'"
		}
		if input.RwFile != "" {
			input.RwFile = strings.Replace(input.RwFile, "[", "", -1)
			input.RwFile = strings.Replace(input.RwFile, "]", "", -1)
			input.RwFile = strings.Replace(input.RwFile, "&", "", -1)
			execArg += " --rwFile '" + input.RwFile + "'"
		}

		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg, action.Language)
		if err != nil {
			log.Logger.Error("Add user action", log.Error(err))
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
			log.Logger.Error("Add user action", log.Error(err))
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
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
			log.Logger.Error("Add user action", log.Error(err))
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		encryptPassword, err := AesEnPassword(input.Guid, input.Seed, password, DEFALT_CIPHER)
		if err != nil {
			err = getPasswordEncodeError(action.Language, err)
			log.Logger.Error("Add user action", log.Error(err))
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
	Language string
}

func (action *DeleteUserAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *DeleteUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *DeleteUserAction) DeleteUserCheckParam(input DeleteUserInput) error {
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}

	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
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
		if strings.Contains(input.UserName, ":") {
			input.UserName = strings.Split(input.UserName, ":")[0]
		}
		if err := action.DeleteUserCheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		execArg := fmt.Sprintf("--action remove --user %s ", input.UserName)
		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg, action.Language)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
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
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
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

type ChangeUserPasswordAction struct {
	Language string
}

func (action *ChangeUserPasswordAction) SetAcceptLanguage(language string) {
	action.Language = language
}

type ChangeUserPasswordInputs struct {
	Inputs []ChangeUserPasswordInput `json:"inputs,omitempty"`
}

type ChangeUserPasswordInput struct {
	CallBackParameter
	Guid      string `json:"guid,omitempty"`
	Seed      string `json:"seed,omitempty"`
	Target    string `json:"target,omitempty"`
	UserName  string `json:"userName,omitempty"`
	Password  string `json:"password,omitempty"`
}

type ChangeUserPasswordOutputs struct {
	Outputs []ChangeUserPasswordOutput `json:"outputs,omitempty"`
}

type ChangeUserPasswordOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

func (action *ChangeUserPasswordAction) CheckParam(input ChangeUserPasswordInput) error {
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}

	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}

	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}

	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}

	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}

	return nil
}

func (action *ChangeUserPasswordAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ChangeUserPasswordInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ChangeUserPasswordAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(ChangeUserPasswordInputs)
	outputs := ChangeUserPasswordOutputs{}
	runAs := ""
	var finalErr error

	for _, input := range inputs.Inputs {
		output := ChangeUserPasswordOutput{
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

		if strings.Contains(input.UserName, ":") {
			input.UserName = strings.Split(input.UserName, ":")[0]
		}
		password := input.Password
		execArg := fmt.Sprintf("--action change_password --user '%s'", input.UserName)

		execArg += " --password '" + password + "'"

		result, err := executeS3Script("user_manage.sh", input.Target, runAs, execArg, action.Language)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
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
			err = fmt.Errorf("Parse SaltApi CallResult meet err=%v ", err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		encryptPassword, err := AesEnPassword(input.Guid, input.Seed, password, DEFALT_CIPHER)
		if err != nil {
			err = getPasswordEncodeError(action.Language, err)
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