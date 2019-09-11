package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	SCRIPT_SAVE_PATH = "/srv/salt/base/"
)

var ScriptPluginActions = make(map[string]Action)

func init() {
	ScriptPluginActions["run"] = new(RunScriptAction)
}

type ScriptPlugin struct {
}

func (plugin *ScriptPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ScriptPluginActions[actionName]

	if !found {
		return nil, fmt.Errorf("Script plugin,action = %s not found", actionName)
	}

	return action, nil
}

type RunScriptInputs struct {
	Inputs []RunScriptInput `json:"inputs,omitempty"`
}

type RunScriptInput struct {
	Guid         string `json:"guid,omitempty"`
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`

	Target  string `json:"target,omitempty"`
	RunAs   string `json:"runas,omitempty"`
	ExecArg string `json:"args,omitempty"`
}

type RunScriptOutputs struct {
	Outputs []RunScriptOutput `json:"outputs"`
}

type RunScriptOutput struct {
	Target  string `json:"target"`
	RetCode int    `json:"retCode"`
	Detail  string `json:"detail"`
	Guid         string `json:"guid,omitempty"`
}

type RunScriptAction struct {
}

func (action *RunScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *RunScriptAction) CheckParam(input interface{}) error {
	inputs, ok := input.(RunScriptInputs)
	if !ok {
		return fmt.Errorf("RunScriptAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.EndPoint == "" {
			return errors.New("Endpoint is empty")
		}
		// if input.AccessKey == "" {
		// 	return errors.New("AccessKey is empty")
		// }
		// if input.SecretKey == "" {
		// 	return errors.New("SecretKey is empty")
		// }
		if input.Target == "" {
			return errors.New("Target is empty")
		}
	}

	return nil
}

func saveFileToSaltMasterBaseDir(fileName string) (string, error) {
	var err error
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Errorf("saveFileToSaltMasterBaseDir,readfile meet err=%v", err)
		return "", err
	}

	tmpFile, err := ioutil.TempFile(SCRIPT_SAVE_PATH, "script-")
	if err != nil {
		logrus.Errorf("saveScript,create tempfile meet err=%v", err)
		return "", err
	}

	defer func() {
		if err != nil {
			defer os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write(content); err != nil {
		logrus.Errorf("saveScript,write tempfile meet err=%v", err)
		return "", err
	}

	if err = tmpFile.Close(); err != nil {
		logrus.Errorf("saveScript,close tempfile meet err=%v", err)
		return "", err
	}

	fullPath := tmpFile.Name()
	return fullPath, err
}

func executeScript(fileName string, target string, runAs string, execArg string) (string, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.script"

	logrus.Infof("executeScript fileName=%s,target=%s,runAs=%s,execArgs=%s", fileName, target, runAs, execArg)

	request.Args = append(request.Args, "salt://base/"+fileName)
	if len(execArg) > 0 {
		request.Args = append(request.Args, "args="+execArg)
	}

	if len(runAs) > 0 {
		request.Args = append(request.Args, "runas="+runAs)
	}

	result, err := CallSaltApi("https://127.0.0.1:8080", request)
	if err != nil {
		return "", err
	}

	return result, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("downloadFile status code=%v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func downLoadAndRunScript(input RunScriptInput) (error, string) {
	// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
	fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
	if err != nil {
		logrus.Errorf("RunScriptAction downloads3 file error=%v", err)
		return err, fmt.Sprintf("RunScriptAction downloads3 file error=%v", err)
	}

	scriptPath, err := saveFileToSaltMasterBaseDir(fileName)
	os.Remove(fileName)
	if err != nil {
		logrus.Errorf("saveFileToSaltMasterBaseDir meet error=%v", err)
		return err, fmt.Sprintf("saveFileToSaltMasterBaseDir meet error=%v", err)
	}

	result, err := executeScript(filepath.Base(scriptPath), input.Target, input.RunAs, input.ExecArg)
	os.Remove(scriptPath)
	if err != nil {
		return err, fmt.Sprintf("executeScript meet error=%v", err)
	}

	saltApiResult, err := parseSaltApiCallResult(result)
	if err != nil {
		logrus.Errorf("parseSaltApiCallResult meet err=%v,rawStr=%s", err, result)
		return err, fmt.Sprintf("parseSaltApiCallResult meet err=%v", err)
	}

	var output string
	for _, v := range saltApiResult.Results[0] {
		if v.RetCode != 0 {
			return fmt.Errorf("script run retCode =%v", v.RetCode), v.Stderr
		}
		output = v.Stdout + v.Stderr
		break
	}
	return nil, output
}

func (action *RunScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunScriptInputs)
	outputs := RunScriptOutputs{}

	for _, input := range inputs.Inputs {
		output := RunScriptOutput{
			Target: input.Target,
		}

		err, stdOut := downLoadAndRunScript(input)
		if err != nil {
			output.RetCode = 1
			return outputs,err
		}
		output.Detail = stdOut
		output.Guid = input.Guid

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
