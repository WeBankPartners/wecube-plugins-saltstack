package plugins

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var RedisUserPluginActions = make(map[string]Action)

func init() {
	RedisUserPluginActions["add"] = new(AddRedisUserAction)
	/*
		RedisUserPluginActions["delete"] = new(DeleteRedisUserAction)
		RedisUserPluginActions["grant"] = new(GrantRedisUserAction)
		RedisUserPluginActions["revoke"] = new(RevokeRedisUserAction)
	*/
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

// AddRedisUserAction add redis user
type AddRedisUserAction struct{ Language string }

type AddRedisUserInputs struct {
	Inputs []AddRedisUserInput `json:"inputs,omitempty"`
}

type AddRedisUserInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`

	//new user info
	NewUserGuid           string `json:"newUserGuid,omitempty"`
	NewUserName           string `json:"newUserName,omitempty"`
	NewUserPassword       string `json:"newUserPassword,omitempty"`
	NewUserReadKeyPrefix  string `json:"newUserReadKeyPrefix,omitempty"`
	NewUserWriteKeyPrefix string `json:"newUserWriteKeyPrefix,omitempty"`
}

type AddRedisUserOutputs struct {
	Outputs []AddRedisUserOutput `json:"outputs,omitempty"`
}

type AddRedisUserOutput struct {
	CallBackParameter
	Result
	NewUserGuid     string `json:"newUserGuid,omitempty"`
	NewUserPassword string `json:"newUserPassword,omitempty"`
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

	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}

	if input.NewUserGuid == "" {
		return getParamEmptyError(action.Language, "newUserGuid")
	}
	if input.NewUserName == "" {
		return getParamEmptyError(action.Language, "newUserName")
	}
	/*
		if input.NewUserPassword == "" {
			return getParamEmptyError(action.Language, "newUserPassword")
		}
	*/
	return nil
}

func (action *AddRedisUserAction) addRedisUser(input *AddRedisUserInput) (output AddRedisUserOutput, err error) {
	defer func() {
		output.NewUserGuid = input.NewUserGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkAddRedisUser(input); err != nil {
		return output, err
	}

	//get login password
	input.Seed = getEncryptSeed(input.Seed)
	loginPassword, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode manage password fail,%s ", err.Error()))
		return output, err
	}

	newUserPassword, decodeErr := AesDePassword(input.NewUserGuid, input.Seed, input.NewUserPassword)
	if decodeErr != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode user password fail,%s ", decodeErr.Error()))
		return output, err
	}

	// check redis user whether is existed.
	isExist, err := redisCheckUserExistOrNot(input.Host, input.Port, input.UserName, loginPassword, input.NewUserName)
	if err != nil {
		return output, err
	}
	if isExist {
		err = getRedisAddUserError(action.Language, input.NewUserName, "user already exist")
		return output, err
	}

	//create user
	//userPassword := input.DatabaseUserPassword
	if newUserPassword == "" {
		newUserPassword = createRandomPassword()
	}

	// create encrypt password
	encryptPassword, encodeErr := AesEnPassword(input.NewUserGuid, input.Seed, newUserPassword, DEFALT_CIPHER)
	if encodeErr != nil {
		err = getPasswordEncodeError(action.Language, encodeErr)
		return output, err
	}

	err = redisCreateUser(input.Host, input.Port, input.UserName, loginPassword, input.NewUserName, newUserPassword, input.NewUserReadKeyPrefix, input.NewUserWriteKeyPrefix)
	if err != nil {
		err = getRedisAddUserError(action.Language, input.NewUserName, err.Error())
		return output, err
	}

	output.NewUserPassword = encryptPassword
	return output, err
}

func (action *AddRedisUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddRedisUserInputs)
	outputs := AddRedisUserOutputs{}
	var finalErr error

	for _, inputData := range inputs.Inputs {
		output, err := action.addRedisUser(&inputData)
		if err != nil {
			log.Logger.Error("Add redis user action", log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}

// DeleteRedisUserAction delete redis user
type DeleteRedisUserAction struct{ Language string }

// todo implement function
