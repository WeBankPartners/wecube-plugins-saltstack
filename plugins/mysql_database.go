package plugins

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var MysqlDatabasePluginActions = make(map[string]Action)

func init() {
	MysqlDatabasePluginActions["add"] = new(AddMysqlDatabaseAction)
	MysqlDatabasePluginActions["delete"] = new(DeleteMysqlDatabaseAction)
}

type MysqlDatabasePlugin struct {
}

func (plugin *MysqlDatabasePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MysqlDatabasePluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("mysql database plugin,action = %s not found", actionName)
	}

	return action, nil
}

type AddMysqlDatabaseAction struct { Language string }

type AddMysqlDatabaseInputs struct {
	Inputs []AddMysqlDatabaseInput `json:"inputs,omitempty"`
}

type AddMysqlDatabaseInput struct {
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
	DatabaseName          string `json:"databaseName,omitempty"`
	DatabaseOwnerGuid     string `json:"databaseOwnerGuid,omitempty"`
	DatabaseOwnerName     string `json:"databaseOwnerName,omitempty"`
	DatabaseOwnerPassword string `json:"databaseOwnerPassword,omitempty"`
	//DatabaseOwnerPermissions string `json:"databaseOwnerPermissions,omitempty"`
}

type AddMysqlDatabaseOutputs struct {
	Outputs []AddMysqlDatabaseOutput `json:"outputs,omitempty"`
}

type AddMysqlDatabaseOutput struct {
	CallBackParameter
	Result
	DatabaseOwnerGuid     string `json:"databaseOwnerGuid,omitempty"`
	DatabaseOwnerPassword string `json:"databaseOwnerPassword,omitempty"`
}

func (action *AddMysqlDatabaseAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *AddMysqlDatabaseAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddMysqlDatabaseInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *AddMysqlDatabaseAction) addMysqlDatabaseCheckParam(input *AddMysqlDatabaseInput) error {
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	if input.DatabaseName == "" {
		return getParamEmptyError(action.Language, "databaseName")
	}
	if input.DatabaseOwnerName == "" {
		return getParamEmptyError(action.Language, "databaseOwnerName")
	}
	if input.DatabaseOwnerGuid == "" {
		return getParamEmptyError(action.Language, "databaseOwnerGuid")
	}
	return nil
}

func (action *AddMysqlDatabaseAction) addMysqlDatabaseAndUser(input *AddMysqlDatabaseInput) (output AddMysqlDatabaseOutput, err error) {
	defer func() {
		output.DatabaseOwnerGuid = input.DatabaseOwnerGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err := action.addMysqlDatabaseCheckParam(input); err != nil {
		return output, err
	}

	//get root password
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	// check database database whether is existed.
	dbIsExist, err := checkDBExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseName, action.Language)
	if err != nil {
		return output, err
	}
	if dbIsExist == true {
		err = getAddMysqlDatabaseError(action.Language, fmt.Sprintf("database %s already exists", input.DatabaseName))
		return output, err
	}

	// check database user whether is existed.
	isExist, err := checkUserExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseOwnerName, action.Language)
	if err != nil {
		return output, err
	}
	if isExist == true {
		err = getAddMysqlDatabaseError(action.Language, fmt.Sprintf("user %s already exists", input.DatabaseOwnerName))
		return output, err
	}

	// create database
	cmd := fmt.Sprintf("create database %s ", input.DatabaseName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	dbOwnerPassword := input.DatabaseOwnerPassword
	if dbOwnerPassword == "" {
		dbOwnerPassword = createRandomPassword()
	}

	// create user
	cmd = fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseOwnerName, dbOwnerPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	// encrypt password
	encryptPassword, err := AesEnPassword(input.DatabaseOwnerGuid, input.Seed, dbOwnerPassword, DEFALT_CIPHER)
	if err != nil {
		err = getPasswordEncodeError(action.Language, err)
		return output, err
	}
	output.DatabaseOwnerPassword = encryptPassword

	// grant permission
	permission := "ALL PRIVILEGES"
	cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseOwnerName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	return output, err
}

func (action *AddMysqlDatabaseAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddMysqlDatabaseInputs)
	outputs := AddMysqlDatabaseOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.addMysqlDatabaseAndUser(&input)
		if err != nil {
			log.Logger.Error("Add mysql database action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}

type DeleteMysqlDatabaseAction struct { Language string }

type DeleteMysqlDatabaseInputs struct {
	Inputs []DeleteMysqlDatabaseInput `json:"inputs,omitempty"`
}

type DeleteMysqlDatabaseInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	// database info
	DatabaseName      string `json:"databaseName,omitempty"`
	DatabaseOwnerGuid string `json:"databaseOwnerGuid,omitempty"`
}

type DeleteMysqlDatabaseOutputs struct {
	Outputs []DeleteMysqlDatabaseOutput `json:"outputs,omitempty"`
}

type DeleteMysqlDatabaseOutput struct {
	CallBackParameter
	Result
	DatabaseOwnerGuid string `json:"databaseOwnerGuid,omitempty"`
}

func (action *DeleteMysqlDatabaseAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *DeleteMysqlDatabaseAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteMysqlDatabaseInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *DeleteMysqlDatabaseAction) deleteMysqlDatabaseCheckParam(input DeleteMysqlDatabaseInput) error {
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	if input.DatabaseName == "" {
		return getParamEmptyError(action.Language, "databaseName")
	}
	if input.DatabaseOwnerGuid == "" {
		return getParamEmptyError(action.Language, "databaseOwnerGuid")
	}

	return nil
}

func (action *DeleteMysqlDatabaseAction) deleteMysqlDatabase(input *DeleteMysqlDatabaseInput) (output DeleteMysqlDatabaseOutput, err error) {
	defer func() {
		output.DatabaseOwnerGuid = input.DatabaseOwnerGuid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.deleteMysqlDatabaseCheckParam(*input)
	if err != nil {
		return output, err
	}

	password, tmpErr := AesDePassword(input.Guid, input.Seed, input.Password)
	if tmpErr != nil {
		err = getPasswordDecodeError(action.Language, tmpErr)
		return output, err
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	// check database database whether is existed.
	dbIsExist, err := checkDBExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseName, action.Language)
	if err != nil {
		return output, err
	}
	if dbIsExist == true {
		var users []string
		users, err = getAllUserByDB(input.Host, input.Port, input.UserName, password, input.DatabaseName, action.Language)
		if err != nil {
			err = getDeleteMysqlDatabaseError(action.Language, err.Error())
			return output, err
		}

		for _, user := range users {
			// revoke permission
			permission := "ALL PRIVILEGES"
			cmd := fmt.Sprintf("REVOKE %s ON %s.* FROM %s ", permission, input.DatabaseName, user)
			if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
				err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
				return output, err
			}
		}
	}

	// delete database
	cmd := fmt.Sprintf("DROP DATABASE %s", input.DatabaseName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		err = getRunMysqlCommnandError(action.Language, cmd, err.Error())
		return output, err
	}

	return output, err
}

func (action *DeleteMysqlDatabaseAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(DeleteMysqlDatabaseInputs)
	outputs := DeleteMysqlDatabaseOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.deleteMysqlDatabase(&input)
		if err != nil {
			log.Logger.Error("Delete mysql database action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
