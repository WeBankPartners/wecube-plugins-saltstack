package plugins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var DatabasePluginActions = make(map[string]Action)

func init() {
	DatabasePluginActions["runScript"] = new(RunDatabaseScriptAction)
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

func (action *RunDatabaseScriptAction) CheckParam(input interface{}) error {
	inputs, ok := input.(RunDatabaseScriptInputs)
	if !ok {
		return fmt.Errorf("RunDatabaseScriptAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
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

func (action *RunDatabaseScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunDatabaseScriptInputs)
	outputs := RunDatabaseScriptOutputs{}

	for _, input := range inputs.Inputs {
		// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
		fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
		if err != nil {
			logrus.Infof("RunScriptAction downloads3 file error=%v", err)
			return nil, err
		}

		md5sum := Md5Encode(input.Guid+input.Seed)
		password,err := AesDecode(md5sum[0:16], input.Password)
		if err != nil {
			logrus.Errorf("AesDecode meet error(%v)", err)
			return nil , err
		}
		
		result, err := execSqlScript(input.Host, input.Port, input.UserName, password, input.DatabaseName, fileName)
		os.Remove(fileName)
		if err != nil {
			return nil, err
		}

		output := RunDatabaseScriptOutput{
			Detail: result,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
