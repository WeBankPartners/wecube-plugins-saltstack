package plugins

import (
	"fmt"
	"regexp"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var RedisUserPluginActions = make(map[string]Action)

func init() {
	RedisUserPluginActions["add"] = new(AddRedisUserAction)
	RedisUserPluginActions["delete"] = new(DeleteRedisUserAction)
}

type RedisUserPlugin struct {
}

func (plugin *RedisUserPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := RedisUserPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("redis user plugin, action = %s not found", actionName)
	}

	return action, nil
}

// AddRedisUserAction add redis user ------------------
type AddRedisUserAction struct{ Language string }

type AddRedisUserInputs struct {
	Inputs []AddRedisUserInput `json:"inputs,omitempty"`
}

type AddRedisUserInput struct {
	CallBackParameter
	Guid          string `json:"guid,omitempty"`
	Seed          string `json:"seed,omitempty"`
	Host          string `json:"host,omitempty"`
	Port          string `json:"port,omitempty"`
	AdminUserName string `json:"adminUserName,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`

	//user info
	UserGuid            string `json:"userGuid,omitempty"`
	UserName            string `json:"userName,omitempty"`
	UserPassword        string `json:"userPassword,omitempty"`
	UserReadKeyPattern  string `json:"userReadKeyPattern,omitempty"`
	UserWriteKeyPattern string `json:"userWriteKeyPattern,omitempty"`
}

type AddRedisUserOutputs struct {
	Outputs []AddRedisUserOutput `json:"outputs,omitempty"`
}

type AddRedisUserOutput struct {
	CallBackParameter
	Result
	UserGuid     string `json:"userGuid,omitempty"`
	UserPassword string `json:"userPassword,omitempty"`
}

func (action *AddRedisUserAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *AddRedisUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddRedisUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

// ValidateChinesePunctuation 校验字符串是否包含常见的中文标点符号

func ValidateChinesePunctuation(str string) bool {
	// 匹配范围包含：
	// 1. CJK 符号和标点 (U+3000 ~ U+303F), 表示Unicode编码范围内的中文标点符号
	// 2. 全角 ASCII 字符 (U+FF00 ~ U+FFEF) 内也包含了部分中文标点，如 "，", "。", "！" 等, 表示Unicode编码范围内的全角字符
	// 也可以补充或移除其它常见符号，比如书名号《》等
	var reZhPunc = regexp.MustCompile("[\u3000-\u303F\uFF00-\uFFEF]")

	if reZhPunc.MatchString(str) {
		return true
	}
	return false
}

func (action *AddRedisUserAction) checkAddRedisUser(input *AddRedisUserInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}

	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Port == "" {
		return getParamEmptyError(action.Language, "port")
	}
	/*
		if input.AdminUserName == "" {
			return getParamEmptyError(action.Language, "adminUserName")
		}
	*/
	if input.AdminPassword == "" {
		return getParamEmptyError(action.Language, "adminPassword")
	}

	if input.UserGuid == "" {
		return getParamEmptyError(action.Language, "userGuid")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	/*
		if input.UserPassword == "" {
			return getParamEmptyError(action.Language, "userPassword")
		}
	*/
	if ValidateChinesePunctuation(input.UserReadKeyPattern) {
		return getParamValidateError(action.Language, "userReadKeyPattern", "Contains illegal character")
	}
	if ValidateChinesePunctuation(input.UserWriteKeyPattern) {
		return getParamValidateError(action.Language, "userWriteKeyPattern", "Contains illegal character")
	}
	return nil
}

func (action *AddRedisUserAction) addRedisUser(input *AddRedisUserInput) (output AddRedisUserOutput, err error) {
	defer func() {
		output.UserGuid = input.UserGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkAddRedisUser(input); err != nil {
		return
	}

	//get admin password
	input.Seed = getEncryptSeed(input.Seed)
	adminPassword, err := AesDePassword(input.Guid, input.Seed, input.AdminPassword)
	if err != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode admin password fail,%s ", err.Error()))
		return
	}

	userPassword, decodeErr := AesDePassword(input.UserGuid, input.Seed, input.UserPassword)
	if decodeErr != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode user password fail,%s ", decodeErr.Error()))
		return
	}

	//create user
	//userPassword := input.DatabaseUserPassword
	if userPassword == "" {
		userPassword = createRandomPassword()
	}

	// create encrypt user password
	encryptUserPassword, encodeErr := AesEnPassword(input.UserGuid, input.Seed, userPassword, DEFALT_CIPHER)
	if encodeErr != nil {
		err = getPasswordEncodeError(action.Language, encodeErr)
		return
	}

	userReadKeyPatterns := splitWithCustomFlag(input.UserReadKeyPattern)
	userWriteKeyPatterns := splitWithCustomFlag(input.UserWriteKeyPattern)
	err = redisCreateUser(input.Host, input.Port, input.AdminUserName, adminPassword, input.UserName, userPassword, userReadKeyPatterns, userWriteKeyPatterns)
	if err != nil {
		err = getRedisAddUserError(action.Language, input.UserName, err.Error())
		return
	}

	output.UserPassword = encryptUserPassword
	return
}

func (action *AddRedisUserAction) Do(input interface{}) (interface{}, error) {
	outputs := AddRedisUserOutputs{}
	var finalErr error

	inputs, isOk := input.(AddRedisUserInputs)
	if !isOk {
		finalErr = fmt.Errorf("input:%v is not the type: AddRedisUserInputs", input)
		return outputs, finalErr
	}

	for _, inputData := range inputs.Inputs {
		output, err := action.addRedisUser(&inputData)
		if err != nil {
			log.Logger.Error("add redis user action failed", log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}

// DeleteRedisUserAction delete redis user -----------------
type DeleteRedisUserAction struct{ Language string }

type DeleteRedisUserInputs struct {
	Inputs []DeleteRedisUserInput `json:"inputs,omitempty"`
}

type DeleteRedisUserInput struct {
	CallBackParameter
	Guid          string `json:"guid,omitempty"`
	Seed          string `json:"seed,omitempty"`
	Host          string `json:"host,omitempty"`
	Port          string `json:"port,omitempty"`
	AdminUserName string `json:"adminUserName,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`

	//user info
	UserGuid string `json:"userGuid,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type DeleteRedisUserOutputs struct {
	Outputs []DeleteRedisUserOutput `json:"outputs,omitempty"`
}

type DeleteRedisUserOutput struct {
	CallBackParameter
	Result
	UserGuid string `json:"userGuid,omitempty"`
}

func (action *DeleteRedisUserAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *DeleteRedisUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteRedisUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *DeleteRedisUserAction) checkDeleteRedisUser(input *DeleteRedisUserInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}

	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Port == "" {
		return getParamEmptyError(action.Language, "port")
	}
	/*
		if input.AdminUserName == "" {
			return getParamEmptyError(action.Language, "adminUserName")
		}
	*/
	if input.AdminPassword == "" {
		return getParamEmptyError(action.Language, "adminPassword")
	}

	if input.UserGuid == "" {
		return getParamEmptyError(action.Language, "userGuid")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}

	return nil
}

func (action *DeleteRedisUserAction) DeleteRedisUser(input *DeleteRedisUserInput) (output DeleteRedisUserOutput, err error) {
	defer func() {
		output.UserGuid = input.UserGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkDeleteRedisUser(input); err != nil {
		return
	}

	//get admin password
	input.Seed = getEncryptSeed(input.Seed)
	adminPassword, err := AesDePassword(input.Guid, input.Seed, input.AdminPassword)
	if err != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode admin password fail,%s ", err.Error()))
		return
	}

	err = redisDeleteUser(input.Host, input.Port, input.AdminUserName, adminPassword, input.UserName)
	if err != nil {
		err = getRedisDeleteUserError(action.Language, input.UserName, err.Error())
		return
	}
	return
}

func (action *DeleteRedisUserAction) Do(input interface{}) (interface{}, error) {
	outputs := DeleteRedisUserOutputs{}
	var finalErr error

	inputs, isOk := input.(DeleteRedisUserInputs)
	if !isOk {
		finalErr = fmt.Errorf("input:%v is not the type: DeleteRedisUserInputs", input)
		return outputs, finalErr
	}

	for _, inputData := range inputs.Inputs {
		output, err := action.DeleteRedisUser(&inputData)
		if err != nil {
			log.Logger.Error("delete redis user action failed", log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}
