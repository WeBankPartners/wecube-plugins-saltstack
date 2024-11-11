package plugins

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var MysqlDatabaseUserPluginActions = make(map[string]Action)

func init() {
	MysqlDatabaseUserPluginActions["add"] = new(AddMysqlDatabaseUserAction)
	MysqlDatabaseUserPluginActions["delete"] = new(DeleteMysqlDatabaseUserAction)
	MysqlDatabaseUserPluginActions["change-password"] = new(ChangeMysqlDatabaseUserPwdAction)
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

// ------------AddMysqlDatabaseUserAction--------------
type AddMysqlDatabaseUserAction struct{ Language string }

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
	input.Seed = getEncryptSeed(input.Seed)
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode manage password fail,%s ", err.Error()))
		return output, err
	}

	userPassword, decodeErr := AesDePassword(input.DatabaseUserGuid, input.Seed, input.DatabaseUserPassword)
	if decodeErr != nil {
		err = getPasswordDecodeError(action.Language, fmt.Errorf("aes decode user password fail,%s ", decodeErr.Error()))
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
	//userPassword := input.DatabaseUserPassword
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
	encryptPassword, encodeErr := AesEnPassword(input.DatabaseUserGuid, input.Seed, userPassword, DEFALT_CIPHER)
	if encodeErr != nil {
		err = getPasswordEncodeError(action.Language, encodeErr)
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

type DeleteMysqlDatabaseUserAction struct{ Language string }

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
	input.Seed = getEncryptSeed(input.Seed)
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

type ChangeMysqlDatabaseUserPwdAction struct{ Language string }

type ChangeMysqlDatabaseUserPwdInputs struct {
	Inputs []ChangeMysqlDatabaseUserPwdInput `json:"inputs,omitempty"`
}

type ChangeMysqlDatabaseUserPwdInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	//new database info
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseUserName     string `json:"databaseUserName,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

type ChangeMysqlDatabaseUserPwdOutputs struct {
	Outputs []ChangeMysqlDatabaseUserPwdOutput `json:"outputs,omitempty"`
}

type ChangeMysqlDatabaseUserPwdOutput struct {
	CallBackParameter
	Result
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

func (action *ChangeMysqlDatabaseUserPwdAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *ChangeMysqlDatabaseUserPwdAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ChangeMysqlDatabaseUserPwdInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ChangeMysqlDatabaseUserPwdAction) checkChangeMysqlDatabaseUserPwd(input *ChangeMysqlDatabaseUserPwdInput) error {
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
	if input.DatabaseUserPassword == "" {
		return getParamEmptyError(action.Language, "databaseUserPassword")
	}
	return nil
}

func (action *ChangeMysqlDatabaseUserPwdAction) changeUserPassword(input *ChangeMysqlDatabaseUserPwdInput) (output ChangeMysqlDatabaseUserPwdOutput, err error) {
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

	if err = action.checkChangeMysqlDatabaseUserPwd(input); err != nil {
		return output, err
	}

	//get root password
	input.Seed = getEncryptSeed(input.Seed)
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	// get new password
	newPassword, decodeErr := AesDePassword(input.DatabaseUserGuid, input.Seed, input.DatabaseUserPassword)
	if decodeErr != nil {
		err = getPasswordDecodeError(action.Language, decodeErr)
		return output, err
	}

	// check database user whether is existed.
	isExist, err := checkUserExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseUserName, action.Language)
	if err != nil {
		return output, err
	}
	if isExist == false {
		err = fmt.Errorf("user %s not exist", input.DatabaseUserName)
		return output, err
	}

	// modify password
	cmd := fmt.Sprintf("set password for '%s'@'%%'=password('%s')", input.DatabaseUserName, newPassword)
	if resetErr := runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); resetErr != nil {
		log.Logger.Warn("Change mysql user password action fail with set password command", log.Error(resetErr))
		cmd = fmt.Sprintf("alter user '%s'@'%%' identified by '%s'", input.DatabaseUserName, newPassword)
		err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd)
		if err != nil {
			log.Logger.Warn("Change mysql user password action fail with alter user command", log.Error(err))
			err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
			return output, err
		}
	}

	// create encrypt password
	encryptPassword, encodeErr := AesEnPassword(input.DatabaseUserGuid, input.Seed, newPassword, DEFALT_CIPHER)
	if encodeErr != nil {
		err = getPasswordEncodeError(action.Language, encodeErr)
		return output, err
	}
	output.DatabaseUserPassword = encryptPassword
	return output, err
}

func (action *ChangeMysqlDatabaseUserPwdAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(ChangeMysqlDatabaseUserPwdInputs)
	outputs := ChangeMysqlDatabaseUserPwdOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.changeUserPassword(&input)
		if err != nil {
			log.Logger.Error("Change mysql user password action", log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}
