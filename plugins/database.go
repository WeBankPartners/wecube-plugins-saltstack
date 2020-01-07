package plugins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

var DatabasePluginActions = make(map[string]Action)

func init() {
	DatabasePluginActions["runScript"] = new(RunDatabaseScriptAction)
	DatabasePluginActions["addDatabase"] = new(AddDatabaseAction)
	DatabasePluginActions["addUser"] = new(AddDatabaseUserAction)
}

type DatabasePlugin struct {
}

func (plugin *DatabasePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := DatabasePluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("database plugin,action = %s not found", actionName)
	}

	return action, nil
}

type RunDatabaseScriptInputs struct {
	Inputs []RunDatabaseScriptInput `json:"inputs,omitempty"`
}

type RunDatabaseScriptInput struct {
	CallBackParameter
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
	Guid         string `json:"guid,omitempty"`
	Seed         string `json:"seed,omitempty"`
	Host         string `json:"host,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
	Port         string `json:"port,omitempty"`
}

type RunDatabaseScriptOutputs struct {
	Outputs []RunDatabaseScriptOutput `json:"outputs,omitempty"`
}

type RunDatabaseScriptOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type RunDatabaseScriptAction struct {
}

func (action *RunDatabaseScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunDatabaseScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func runDatabaseScriptCheckParam(input RunDatabaseScriptInput) error {
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
	if input.EndPoint == "" {
		return errors.New("EndPoint is empty")
	}

	if input.Port == "" {
		input.Port = "3306"
	}

	return nil
}

func execSqlScript(hostName string, port string, userName string, password string, databaseName string, fileName string) (string, error) {
	argv := []string{
		"-h" + hostName,
		"-u" + userName,
		"-p" + password,
		"-P" + port,
	}

	if databaseName != "" {
		argv = append(argv, "-D"+databaseName)
	}

	cmd := exec.Command("/usr/bin/mysql", argv...)
	f, err := os.Open(fileName)
	if err != nil {
		logrus.Errorf("open file failed err=%v", err)
		return "", err
	}

	defer f.Close()
	cmd.Stdin = f

	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("failed to execSqlScript err=%v,output=%v", err, string(out))
		return "", fmt.Errorf("failed to execSqlScript, err = %v,output=%v", err, string(out))
	}

	return string(out), nil
}

func (action *RunDatabaseScriptAction) runDatabaseScript(input *RunDatabaseScriptInput) (output RunDatabaseScriptOutput, err error) {
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
	err = runDatabaseScriptCheckParam(*input)
	if err != nil {
		return output, err
	}

	// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
	fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
	if err != nil {
		return output, err
	}

	md5sum := Md5Encode(input.Guid + input.Seed)
	password, err := AesDecode(md5sum[0:16], input.Password)
	if err != nil {
		return output, err
	}

	// new dir to place all *.sql
	Info := strings.Split(fileName, "/")
	newDir := strings.Join(Info[0:len(Info)-2], "/") + "/sql"
	err = ensureDirExist(newDir)
	if err != nil {
		return output, err
	}

	// whether the fileName is *.sql or other
	if !strings.HasSuffix(fileName, ".sql") {
		// if the fileName is unpack package, unpack
		input := FileCopyInput{
			Target:          fileName,
			DestinationPath: newDir,
		}
		actionFileCopy := &FileCopyAction{}
		unpackRequest, er := actionFileCopy.deriveUnpackRequest(&input)
		if er != nil {
			err = er
			return output, err
		}
		if _, err = CallSaltApi("https://127.0.0.1:8080", *unpackRequest); err != nil {
			return output, err
		}
	} else {
		// move the *.sql to newDir directly
		command := exec.Command("mv", fileName, newDir)
		out, er := command.CombinedOutput()
		logrus.Infof("runDatabaseCommand(%v) output=%v,err=%v\n", command, string(out), err)
		if er != nil {
			err = er
			return output, err
		}
	}

	// list all *.sql in newDir and run sql scripts foreach
	files, err := listFile(newDir)
	if err != nil {
		return output, err
	}

	for _, file := range files {
		_, err = execSqlScript(input.Host, input.Port, input.UserName, password, input.DatabaseName, file)
		if err != nil {
			return output, err
		}

	}

	return output, err
}

func (action *RunDatabaseScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunDatabaseScriptInputs)
	outputs := RunDatabaseScriptOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.runDatabaseScript(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	logrus.Infof("all scripts = %v have been run", inputs)

	return &outputs, finalErr
}

//----------------------add db user----------------------//
type AddDatabaseAction struct {
}

type AddDatabaseInputs struct {
	Inputs []AddDatabaseInput `json:"inputs,omitempty"`
}

type AddDatabaseInput struct {
	CallBackParameter
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	//new database info
	DatabaseName          string `json:"databaseName,omitempty"`
	DatabaseOwnerGuid     string `json:"databaseOwnerGuid,omitempty"`
	DatabaseOwnerName     string `json:"databaseOwnerName,omitempty"`
	DatabaseOwnerPassword string `json:"databaseOwnerPassword,omitempty"`
	//DatabaseOwnerPermissions string `json:"databaseOwnerPermissions,omitempty"`
}

type AddDatabaseOutputs struct {
	Outputs []AddDatabaseOutput `json:"outputs,omitempty"`
}

type AddDatabaseOutput struct {
	CallBackParameter
	Result
	DatabaseOwnerGuid     string `json:"databaseOwnerGuid,omitempty"`
	DatabaseOwnerPassword string `json:"databaseOwnerPassword,omitempty"`
}

func (action *AddDatabaseAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddDatabaseInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func addDatabaseCheckParam(input *AddDatabaseInput) error {
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
	if input.Port == "" {
		return errors.New("Port is empty")
	}
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

func AddDatabaseAndUser(input *AddDatabaseInput) (string, error) {
	if err := addDatabaseCheckParam(input); err != nil {
		return "", err
	}

	//get root password
	md5sum := Md5Encode(input.Guid + input.Seed)
	password, err := AesDecode(md5sum[0:16], input.Password)
	if err != nil {
		return "", err
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
	md5sum = Md5Encode(input.DatabaseOwnerGuid + input.Seed)
	encryptPassword, err := AesEncode(md5sum[0:16], dbOwnerPassword)
	return encryptPassword, err
}

func (action *AddDatabaseAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddDatabaseInputs)
	outputs := AddDatabaseOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := AddDatabaseOutput{
			DatabaseOwnerGuid: input.DatabaseOwnerGuid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		password, err := AddDatabaseAndUser(&input)
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

//------------AddDatabaseUserAction--------------//
type AddDatabaseUserAction struct {
}

type AddDatabaseUserInputs struct {
	Inputs []AddDatabaseUserInput `json:"inputs,omitempty"`
}

type AddDatabaseUserInput struct {
	CallBackParameter
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     string `json:"port,omitempty"`

	//new database info
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseName         string `json:"databaseName,omitempty"`
	DatabaseUserName     string `json:"databaseUserName,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

type AddDatabaseUserOutputs struct {
	Outputs []AddDatabaseUserOutput `json:"outputs,omitempty"`
}

type AddDatabaseUserOutput struct {
	CallBackParameter
	Result
	DatabaseUserGuid     string `json:"databaseUserGuid,omitempty"`
	DatabaseUserPassword string `json:"databaseUserPassword,omitempty"`
}

func (action *AddDatabaseUserAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AddDatabaseUserInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func checkAddDatabaseUser(input *AddDatabaseUserInput) error {
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

func createUserForExistedDatabase(input *AddDatabaseUserInput) (string, error) {
	if err := checkAddDatabaseUser(input); err != nil {
		return "", err
	}

	//get root password
	md5sum := Md5Encode(input.Guid + input.Seed)
	password, err := AesDecode(md5sum[0:16], input.Password)
	if err != nil {
		return "", err
	}

	//create user
	userPassword := input.DatabaseUserPassword
	if userPassword == "" {
		userPassword = createRandomPassword()
	}

	cmd := fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s' ", input.DatabaseUserName, userPassword)
	if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
		return "", err
	}

	//grant permission
	if input.DatabaseName != "" {
		permission := "ALL PRIVILEGES"
		cmd = fmt.Sprintf("GRANT %s ON %s.* TO %s ", permission, input.DatabaseName, input.DatabaseUserName)
		if err = runDatabaseCommand(input.Host, input.Port, input.UserName, password, cmd); err != nil {
			return "", err
		}
	}

	//create encrypt password
	md5sum = Md5Encode(input.DatabaseUserGuid + input.Seed)
	encryptPassword, err := AesEncode(md5sum[0:16], userPassword)
	return encryptPassword, err
}

func (action *AddDatabaseUserAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(AddDatabaseUserInputs)
	outputs := AddDatabaseUserOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := AddDatabaseUserOutput{
			DatabaseUserGuid: input.DatabaseUserGuid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		password, err := createUserForExistedDatabase(&input)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
		output.DatabaseUserPassword = password
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}
