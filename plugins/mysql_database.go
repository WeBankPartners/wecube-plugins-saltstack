package plugins

import (
	"errors"
	"fmt"
	"os/exec"

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

func addMysqlDatabaseCheckParam(input *AddMysqlDatabaseInput) error {
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
	// if input.Port == "" {
	// 	return errors.New("Port is empty")
	// }
	if input.DatabaseName == "" {
		return errors.New("DatabaseName is empty")
	}
	if input.DatabaseOwnerGuid == "" {
		return errors.New("DatabaseOwnerGuid is empty")
	}
	if input.DatabaseOwnerName == "" {
		return errors.New("DatabaseOwnerName is empty")
	}
	return nil
}

func runDatabaseCommand(host string, port string, loginUser string, loginPwd string, cmd string) error {
	argv := []string{
		"-h" + host,
		"-u" + loginUser,
		"-p" + loginPwd,
		"-P" + port,
		"-e",
		cmd,
	}
	command := exec.Command("/usr/bin/mysql", argv...)
	out, err := command.CombinedOutput()
	fmt.Printf("runDatabaseCommand(%v) output=%v,err=%v\n", command, string(out), err)
	return err
}

func AddMysqlDatabaseAndUser(input *AddMysqlDatabaseInput) (string, error) {
	if err := addMysqlDatabaseCheckParam(input); err != nil {
		return "", err
	}

	//get root password
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		logrus.Errorf("AesDePassword meet error(%v)", err)
		return "", err
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	//create database
	cmd := fmt.Sprintf("create database %s ", input.DatabaseName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return "", err
	}

	//create user
	dbOwnerPassword := input.DatabaseOwnerPassword
	if dbOwnerPassword == "" {
		dbOwnerPassword = createRandomPassword()
	}
	cmd = fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseOwnerName, dbOwnerPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return "", err
	}

	//grant permission
	permission := "ALL PRIVILEGES"
	cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseOwnerName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return "", err
	}

	//create new password
	encryptPassword, err := AesEnPassword(input.DatabaseOwnerGuid, input.Seed, dbOwnerPassword, DEFALT_CIPHER)
	if err != nil {
		logrus.Errorf("AesEnPassword meet error(%v)", err)
		return "", err
	}
	return encryptPassword, err
}

func (action *AddMysqlDatabaseAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddMysqlDatabaseInputs)
	outputs := AddMysqlDatabaseOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := AddMysqlDatabaseOutput{
			DatabaseOwnerGuid: input.DatabaseOwnerGuid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		password, err := AddMysqlDatabaseAndUser(&input)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}
		output.DatabaseOwnerPassword = password
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
	DatabaseOwnerName string `json:"databaseOwnerName,omitempty"`
}

type DeleteMysqlDatabaseOutputs struct {
	Outputs []DeleteMysqlDatabaseOutput `json:"outputs,omitempty"`
}

type DeleteMysqlDatabaseOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
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
	// if input.Port == "" {
	// 	return errors.New("Port is empty")
	// }
	if input.DatabaseName == "" {
		return errors.New("DatabaseName is empty")
	}
	// if input.DatabaseOwnerName == "" {
	// 	return errors.New("DatabaseOwnerName is empty")
	// }
	return nil
}

func (action *DeleteMysqlDatabaseAction) deleteMysqlDatabase(input *DeleteMysqlDatabaseInput) (output DeleteMysqlDatabaseOutput, err error) {
	defer func() {
		output.Guid = input.Guid
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

	if input.DatabaseOwnerName != "" {
		// revoke permission
		permission := "ALL PRIVILEGES"
		cmd := fmt.Sprintf("REVOKE %s ON %s.* FROM %s ", permission, input.DatabaseName, input.DatabaseOwnerName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			return output, err
		}

		// delete user
		cmd = fmt.Sprintf("DROP USER %s", input.DatabaseOwnerName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			return output, err
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
