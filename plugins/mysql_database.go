package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
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

type AddMysqlDatabaseAction struct {
}

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

func (action *AddMysqlDatabaseAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddMysqlDatabaseInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *AddMysqlDatabaseAction) addMysqlDatabaseCheckParam(input *AddMysqlDatabaseInput) error {
	if input.Host == "" {
		return errors.New("Host is empty")
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty")
	}
	if input.UserName == "" {
		return errors.New("UserName is empty")
	}
	if input.Password == "" {
		return errors.New("Password is empty")
	}
	if input.DatabaseName == "" {
		return errors.New("DatabaseName is empty")
	}
	if input.DatabaseOwnerName == "" {
		return errors.New("DatabaseOwnerName is empty")
	}
	if input.DatabaseOwnerGuid == "" {
		return errors.New("DatabaseOwnerGuid is empty")
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
		logrus.Errorf("AesDePassword meet error(%v)", err)
		return output, err
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	// check database database whether is existed.
	dbIsExist, err := checkDBExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseName)
	if err != nil {
		logrus.Errorf("check db[%v] exist or not meet error=%v", input.DatabaseName, err)
		return output, err
	}
	if dbIsExist == true {
		logrus.Errorf("db[%v] is existed", input.DatabaseName)
		err = fmt.Errorf("db[%v] is existed", input.DatabaseName)
		return output, err
	}

	// check database user whether is existed.
	isExist, err := checkUserExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseOwnerName)
	if err != nil {
		logrus.Errorf("checking user exist or not meet error=%v", err)
		return output, err
	}
	if isExist == true {
		logrus.Errorf("user[%v] is existed", input.DatabaseOwnerName)
		err = fmt.Errorf("user[%v] is existed", input.DatabaseOwnerName)
		return output, err
	}

	// create database
	cmd := fmt.Sprintf("create database %s ", input.DatabaseName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return output, err
	}

	dbOwnerPassword := input.DatabaseOwnerPassword
	if dbOwnerPassword == "" {
		dbOwnerPassword = createRandomPassword()
	}

	// create user
	cmd = fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseOwnerName, dbOwnerPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return output, err
	}

	// encrypt password
	encryptPassword, err := AesEnPassword(input.DatabaseOwnerGuid, input.Seed, dbOwnerPassword, DEFALT_CIPHER)
	if err != nil {
		logrus.Errorf("AesEnPassword meet error(%v)", err)
		return output, err
	}
	output.DatabaseOwnerPassword = encryptPassword

	// grant permission
	permission := "ALL PRIVILEGES"
	cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseOwnerName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
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
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}

type DeleteMysqlDatabaseAction struct {
}

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

func (action *DeleteMysqlDatabaseAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteMysqlDatabaseInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *DeleteMysqlDatabaseAction) deleteMysqlDatabaseCheckParam(input DeleteMysqlDatabaseInput) error {
	if input.Host == "" {
		return errors.New("Host is empty")
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty")
	}
	if input.UserName == "" {
		return errors.New("UserName is empty")
	}
	if input.Password == "" {
		return errors.New("Password is empty")
	}
	if input.DatabaseName == "" {
		return errors.New("DatabaseName is empty")
	}
	if input.DatabaseOwnerGuid == "" {
		return errors.New("DatabaseOwnerGuid is empty")
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

	password, er := AesDePassword(input.Guid, input.Seed, input.Password)
	if er != nil {
		err = er
		return output, er
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	// check database database whether is existed.
	dbIsExist, err := checkDBExistOrNot(input.Host, input.Port, input.UserName, password, input.DatabaseName)
	if err != nil {
		logrus.Errorf("check db[%v] exist or not meet error=%v", input.DatabaseName, err)
		return output, err
	}
	if dbIsExist == true {
		var users []string
		users, err = getAllUserByDB(input.Host, input.Port, input.UserName, password, input.DatabaseName)
		if err != nil {
			logrus.Errorf("get user by db[%v] meet err=%v", input.DatabaseName, err)
			return output, err
		}

		for _, user := range users {
			// revoke permission
			permission := "ALL PRIVILEGES"
			cmd := fmt.Sprintf("REVOKE %s ON %s.* FROM %s ", permission, input.DatabaseName, user)
			if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
				return output, err
			}
		}
	}

	// delete database
	cmd := fmt.Sprintf("DROP DATABASE %s", input.DatabaseName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
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
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
