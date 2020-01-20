package plugins

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
type AddMysqlDatabaseUserAction struct {
}

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

func (action *AddMysqlDatabaseUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddMysqlDatabaseUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func checkAddMysqlDatabaseUser(input *AddMysqlDatabaseUserInput) error {
	if input.Guid == "" {
		return fmt.Errorf("empty guid")
	}
	if input.Seed == "" {
		return fmt.Errorf("empty seed")
	}
	if input.Password == "" {
		return fmt.Errorf("empty password")
	}
	if input.DatabaseUserName == "" {
		return fmt.Errorf("empty databaseUserName")
	}
	if input.DatabaseUserGuid == "" {
		return fmt.Errorf("empty databaseUserGuid")
	}
	return nil
}

func createUserForExistedDatabase(input *AddMysqlDatabaseUserInput) (output AddMysqlDatabaseUserOutput, err error) {
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

	if err = checkAddMysqlDatabaseUser(input); err != nil {
		return output, err
	}

	// check database user whether is existed.
	isExist, err := checkUserExistOrNot(input.Host, input.Port, input.UserName, input.Password, input.DatabaseUserName)
	if err != nil {
		logrus.Errorf("checking user exist or not meet error=%v", err)
		return output, err
	}
	if isExist == true {
		logrus.Errorf("database user[%v] exsit", input.DatabaseUserName)
		err = fmt.Errorf("database user[%v] exsit", input.DatabaseUserName)
		return output, err
	}

	//get root password
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		logrus.Errorf("AesDePassword meet error(%v)", err)
		return output, err
	}

	//create user
	userPassword := input.DatabaseUserPassword
	if userPassword == "" {
		userPassword = createRandomPassword()
	}

	cmd := fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseUserName, userPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return output, err
	}

	//grant permission
	if input.DatabaseName != "" {
		permission := "ALL PRIVILEGES"
		cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseUserName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			return output, err
		}
	}

	//create encrypt password
	encryptPassword, err := AesEnPassword(input.Guid, input.Seed, userPassword, DEFALT_CIPHER)
	if err != nil {
		logrus.Errorf("AesEnPassword meet error(%v)", err)
		return output, err
	}
	output.DatabaseUserPassword = encryptPassword
	return output, err
}

func checkUserExistOrNot(host, port, loginUser, loginPwd, userName string) (bool, error) {
	// initDB param dbName = "mysql".
	DB, err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return false, err
	}

	querySql := fmt.Sprintf("SELECT 1 FROM mysql.user WHERE user = '%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Errorf("db.query meet err=%v", err)
		return false, err
	}

	return rows.Next(), nil
}

func (action *AddMysqlDatabaseUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddMysqlDatabaseUserInputs)
	outputs := AddMysqlDatabaseUserOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := createUserForExistedDatabase(&input)
		if err != nil {
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}

type DeleteMysqlDatabaseUserAction struct {
}

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
	DatabaseName     string `json:"databaseName,omitempty"`
	DatabaseUserName string `json:"databaseUserName,omitempty"`
}

type DeleteMysqlDatabaseUserOutputs struct {
	Outputs []DeleteMysqlDatabaseUserOutput `json:"outputs,omitempty"`
}

type DeleteMysqlDatabaseUserOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
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
		return fmt.Errorf("empty guid")
	}
	if input.Seed == "" {
		return fmt.Errorf("empty seed")
	}
	if input.Password == "" {
		return fmt.Errorf("empty password")
	}
	if input.DatabaseUserName == "" {
		return fmt.Errorf("empty databaseUserName")
	}
	return nil
}

func (action *DeleteMysqlDatabaseUserAction) deleteMysqlDatabaseUser(input *DeleteMysqlDatabaseUserInput) (output DeleteMysqlDatabaseUserOutput, err error) {
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

	if err = action.deleteMysqlDatabaseUserCheckParam(input); err != nil {
		return output, err
	}

	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		return output, err
	}

	dbs, err := getAllDBByUser(input.Host, input.Port, input.UserName, input.Password, input.DatabaseUserName)
	if err != nil {
		logrus.Errorf("getting dbs by user[%v] meet error=%v", input.DatabaseUserName, err)
		return output, err
	}

	for _, db := range dbs {
		// revoke permissions
		permission := "ALL PRIVILEGES"
		cmd := fmt.Sprintf("REVOKE %s ON %s.* FROM %s ", permission, db, input.DatabaseUserName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			return output, err
		}
	}

	// delete user
	cmd := fmt.Sprintf("DROP USER %s", input.DatabaseUserName)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
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
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}

func getAllDBByUser(host, port, loginUser, loginPwd, userName string) ([]string, error) {
	dbs := []string{}
	// initDB param dbName = "mysql".
	DB, err := initDB(host, port, loginUser, loginPwd, "mysql")
	if err != nil {
		logrus.Errorf("init myhsql db failed, err=%v ", err)
		return dbs, err
	}

	querySql := fmt.Sprintf("select Db from db where db.User='%s'", userName)
	rows, err := DB.Query(querySql)
	if err != nil {
		logrus.Infof("db.query meet err=%v", err)
		return dbs, err
	}
	for rows.Next() {
		var db string
		err := rows.Scan(&db)
		if err != nil {
			logrus.Infof("rows.Scan meet err=%v", err)
			return dbs, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}
