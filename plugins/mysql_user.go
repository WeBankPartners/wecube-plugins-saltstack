package plugins

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var MysqlDatabaseUserPluginActions = make(map[string]Action)

func init() {
	MysqlDatabaseUserPluginActions["add"] = new(AddMysqlDatabaseUserAction)
	MysqlDatabaseUserPluginActions["delete"] = new(DeleteMysqlDatabaseUserAction)
}

type MysqlUserPlugin struct {
}

func (plugin *MysqlUserPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MysqlDatabaseUserPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("mysql user plugin,action = %s not found", actionName)
	}

	return action, nil
}

//------------AddMysqlDatabaseUserAction--------------
type AddMysqlDatabaseUserAction struct {  Language string  }

type AddMysqlDatabaseUserInputs struct {
	Inputs []AddMysqlDatabaseUserInput `json:"inputs,omitempty"`
}

type AddMysqlDatabaseUserInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`

	//new database info
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseName         string `json:"databaseName,omitempty"`
	DatabaseUserName     string `json:"databaseUserName,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

type AddMysqlDatabaseUserOutputs struct {
	Outputs []AddMysqlDatabaseUserOutput `json:"outputs,omitempty"`
}

type AddMysqlDatabaseUserOutput struct {
	CallBackParameter
	Result
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

func (action *AddMysqlDatabaseUserAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *AddMysqlDatabaseUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddMysqlDatabaseUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *AddMysqlDatabaseUserAction) checkAddMysqlDatabaseUser(input *AddMysqlDatabaseUserInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	if input.DatabaseUserName == "" {
		return getParamEmptyError(action.Language, "databaseUserName")
	}
	if input.DatabaseUserGuid == "" {
		return getParamEmptyError(action.Language, "databaseUserGuid")
	}
	return nil
}

func (action *AddMysqlDatabaseUserAction) createUserForExistedDatabase(input *AddMysqlDatabaseUserInput) (output AddMysqlDatabaseUserOutput, err error) {
	defer func() {
		output.DatabaseUserGuid = input.DatabaseUserGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkAddMysqlDatabaseUser(input); err != nil {
		return output, err
	}

	//get root password
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	// check database user whether is existed.
	isExist, err := checkUserExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseUserName, action.Language)
	if err != nil {
		return output, err
	}
	if isExist == true {
		err = getMysqlCreateUserError(action.Language, input.DatabaseUserName, "user already exist")
		return output, err
	}

	//create user
	userPassword := input.DatabaseUserPassword
	if userPassword == "" {
		userPassword = createRandomPassword()
	}

	cmd := fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseUserName, userPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	// grant permission
	if input.DatabaseName != "" {
		permission := "ALL PRIVILEGES"
		cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseUserName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
			return output, err
		}
	}

	// create encrypt password
	encryptPassword, err := AesEnPassword(input.Guid, input.Seed, userPassword, DEFALT_CIPHER)
	if err != nil {
		err = getPasswordEncodeError(action.Language, err)
		return output, err
	}
	output.DatabaseUserPassword = encryptPassword
	return output, err
}

func (action *AddMysqlDatabaseUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddMysqlDatabaseUserInputs)
	outputs := AddMysqlDatabaseUserOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.createUserForExistedDatabase(&input)
		if err != nil {
			log.Logger.Error("Add mysql user action", log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}

type DeleteMysqlDatabaseUserAction struct { Language string }

type DeleteMysqlDatabaseUserInputs struct {
	Inputs []DeleteMysqlDatabaseUserInput `json:"inputs,omitempty"`
}

type DeleteMysqlDatabaseUserInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	//database info
	DatabaseUserName string `json:"databaseUserName,omitempty"`
	DatabaseUserGuid string `json:"databaseUserGuid,omitempty"`
}

type DeleteMysqlDatabaseUserOutputs struct {
	Outputs []DeleteMysqlDatabaseUserOutput `json:"outputs,omitempty"`
}

type DeleteMysqlDatabaseUserOutput struct {
	CallBackParameter
	Result
	DatabaseUserGuid string `json:"databaseUserGuid,omitempty"`
}

func (action *DeleteMysqlDatabaseUserAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *DeleteMysqlDatabaseUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteMysqlDatabaseUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action DeleteMysqlDatabaseUserAction) deleteMysqlDatabaseUserCheckParam(input *DeleteMysqlDatabaseUserInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	if input.DatabaseUserName == "" {
		return getParamEmptyError(action.Language, "databaseUserName")
	}
	if input.DatabaseUserGuid == "" {
		return getParamEmptyError(action.Language, "databaseUserGuid")
	}
	return nil
}

func (action *DeleteMysqlDatabaseUserAction) deleteMysqlDatabaseUser(input *DeleteMysqlDatabaseUserInput) (output DeleteMysqlDatabaseUserOutput, err error) {
	defer func() {
		output.DatabaseUserGuid = input.DatabaseUserGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.deleteMysqlDatabaseUserCheckParam(input); err != nil {
		return output, err
	}

	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	dbs, err := getAllDBByUser(input.Host, input.Port, input.UserName, password, input.DatabaseUserName, action.Language)
	if err != nil {
		return output, err
	}

	for _, db := range dbs {
		// revoke permissions
		permission := "ALL PRIVILEGES"
		cmd := fmt.Sprintf("REVOKE %s ON %s.* FROM %s ", permission, db, input.DatabaseUserName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
			return output, err
		}
	}

	// delete user
	cmd := fmt.Sprintf("DROP USER %s", input.DatabaseUserName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	return output, err
}

func (action *DeleteMysqlDatabaseUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(DeleteMysqlDatabaseUserInputs)
	outputs := DeleteMysqlDatabaseUserOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.deleteMysqlDatabaseUser(&input)
		if err != nil {
			log.Logger.Error("Delete mysql user action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
