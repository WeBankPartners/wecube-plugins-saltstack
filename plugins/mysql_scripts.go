package plugins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

var MysqlScriptPluginActions = make(map[string]Action)

func init() {
	MysqlScriptPluginActions["run"] = new(RunMysqlScriptAction)
}

type MysqlScriptPlugin struct {
}

func (plugin *MysqlScriptPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MysqlScriptPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("mysql script plugin,action = %s not found", actionName)
	}

	return action, nil
}

type RunMysqlScriptInputs struct {
	Inputs []RunMysqlScriptInput `json:"inputs,omitempty"`
}

type RunMysqlScriptInput struct {
	CallBackParameter
	EndPoint     string `json:"endpoint,omitempty"`
	SqlFiles     string `json:"sql_files,omitempty"`
	Guid         string `json:"guid,omitempty"`
	Seed         string `json:"seed,omitempty"`
	Host         string `json:"host,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
	Port         string `json:"port,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
}

type RunMysqlScriptOutputs struct {
	Outputs []RunMysqlScriptOutput `json:"outputs,omitempty"`
}

type RunMysqlScriptOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type RunMysqlScriptAction struct {
}

func (action *RunMysqlScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunMysqlScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func runMysqlScriptCheckParam(input RunMysqlScriptInput) error {
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

func (action *RunMysqlScriptAction) runMysqlScript(input *RunMysqlScriptInput) (output RunMysqlScriptOutput, err error) {
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
	err = runMysqlScriptCheckParam(*input)
	if err != nil {
		return output, err
	}

	// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
	fileName, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false)
	if err != nil {
		return output, err
	}

	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		logrus.Errorf("AesDecode meet error(%v)", err)
		return output, err
	}

	// new dir to place all *.sql
	Info := strings.Split(fileName, "/")
	newDir := strings.Join(Info[0:len(Info)-2], "/") + "/sql"
	err = ensureDirExist(newDir)
	if err != nil {
		return output, err
	}

	files := []string{}
	// whether the fileName is *.sql or other
	if !strings.HasSuffix(fileName, ".sql") {
		if input.SqlFiles == "" {
			err = errors.New("SqlFiles is empty")
			return output, err
		}

		// unpack file
		er := deriveUnpackfile(fileName, newDir, true)
		if er != nil {
			err = er
			return output, err
		}

		// split SqlFiles to *.sql
		sqlFiles := strings.Split(input.SqlFiles, ",")
		for _, file := range sqlFiles {
			sqlFile := newDir + "/" + strings.TrimSpace(file)
			if !fileExist(sqlFile) {
				err = fmt.Errorf("file [%v] does not exist", sqlFile)
				return output, err
			}
			files = append(files, sqlFile)
		}
	} else {
		// move the *.sql to newDir directly
		command := exec.Command("mv", fileName, newDir)
		out, er := command.CombinedOutput()
		logrus.Infof("runDatabaseCommand(%v) output=%v,err=%v\n", command, string(out), er)
		if er != nil {
			err = er
			return output, err
		}

		sqlFile := newDir + "/" + Info[len(Info)-1]
		if fileExist(sqlFile) {
			err = fmt.Errorf("file [%v] does not exist", sqlFile)
			return output, err
		}
		files = append(files, sqlFile)
	}

	// run sql scripts foreach
	for _, file := range files {
		_, err = execSqlScript(input.Host, input.Port, input.UserName, password, input.DatabaseName, file)
		if err != nil {
			return output, err
		}
	}

	err = os.RemoveAll(fileName)
	if err != nil {
		return output, err
	}
	err = os.RemoveAll(newDir)
	if err != nil {
		return output, err
	}

	return output, err
}

func (action *RunMysqlScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunMysqlScriptInputs)
	outputs := RunMysqlScriptOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.runMysqlScript(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	logrus.Infof("all mysql scripts = %v have been run", inputs)

	return &outputs, finalErr
}
