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

	END_POINT_TYPE_S3         = "S3"
	END_POINT_TYPE_LOCAL      = "LOCAL"
	END_POINT_TYPE_USER_PARAM = "USER_PARAM"
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
	CallBackParameter
	EndPointType  string `json:"endpointType,omitempty"` // "S3" or "LOCAL", Defalt: "LOCAL"
	EndPoint      string `json:"endpoint,omitempty"`
	ScriptContent string `json:"scriptContent,omitempty"`
	Target        string `json:"target,omitempty"`
	RunAs         string `json:"runAs,omitempty"`
	ExecArg       string `json:"args,omitempty"`
	Guid          string `json:"guid,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
}

type RunScriptOutputs struct {
	Outputs []RunScriptOutput `json:"outputs"`
}

type RunScriptOutput struct {
	CallBackParameter
	Result
	Target  string `json:"target"`
	RetCode int    `json:"retCode"`
	Detail  string `json:"detail"`
	Guid    string `json:"guid,omitempty"`
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

func (action *RunScriptAction) CheckParam(input RunScriptInput) error {
	if input.EndPointType != END_POINT_TYPE_LOCAL && input.EndPointType != END_POINT_TYPE_S3 && input.EndPointType != END_POINT_TYPE_USER_PARAM {
		return errors.New("Wrong EndPointType")
	}
	if input.EndPoint == "" && input.EndPointType == END_POINT_TYPE_S3 {
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

func executeS3Script(fileName string, target string, runAs string, execArg string) (string, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.script"

	logrus.Infof("executeS3Script fileName=%s,target=%s,runAs=%s,execArgs=%s", fileName, target, runAs, execArg)

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

func executeLocalScript(fileName string, target string, runAs string, execArg string) (string, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.run"

	logrus.Infof("executeLocalScript fileName=%s,target=%s,runAs=%s,execArgs=%s", fileName, target, runAs, execArg)

	cmdRun := "chmod +x " + fileName + ";" + fileName
	if len(execArg) > 0 {
		cmdRun = cmdRun + " " + execArg
	}
	request.Args = append(request.Args, cmdRun)

	if len(runAs) > 0 {
		request.Args = append(request.Args, "runas="+runAs)
	}
	logrus.Infof("executeLocalScript request=%v", request)

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

func downLoadScript(input RunScriptInput) (string, error) {
	// fileName, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
	fileName, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
	if err != nil {
		logrus.Errorf("RunScriptAction downloads3 file error=%v", err)
		return fmt.Sprintf("RunScriptAction downloads3 file error=%v", err), err
	}

	scriptPath, err := saveFileToSaltMasterBaseDir(fileName)
	os.Remove(fileName)
	if err != nil {
		logrus.Errorf("saveFileToSaltMasterBaseDir meet error=%v", err)
		return fmt.Sprintf("saveFileToSaltMasterBaseDir meet error=%v", err), err
	}

	return scriptPath, nil
}

func runScript(scriptPath string, input RunScriptInput) (string, error) {
	var result string
	var output string
	var err error

	switch input.EndPointType {
	case END_POINT_TYPE_LOCAL:
		result, err = executeLocalScript(scriptPath, input.Target, input.RunAs, input.ExecArg)
		if err != nil {
			return fmt.Sprintf("executeLocalScript meet error=%v", err), err
		}
		saltApiResult, err := parseSaltApiCmdRunCallResult(result)
		if err != nil {
			logrus.Errorf("parseSaltApiCmdRunCallResult meet err=%v,rawStr=%s", err, result)
			return fmt.Sprintf("parseSaltApiCmdRunCallResult meet err=%v", err), err
		}
		for k, v := range saltApiResult.Results[0] {
			if k != input.Target {
				err = fmt.Errorf("script run ip[%s] is not target[%s]", k, input.Target)
				return fmt.Sprintf("parseSaltApiCmdRunCallResult meet error=%v", err), err
			}
			output = k + ":" + v
			break
		}
	case END_POINT_TYPE_S3, END_POINT_TYPE_USER_PARAM:
		result, err = executeS3Script(filepath.Base(scriptPath), input.Target, input.RunAs, input.ExecArg)
		os.Remove(scriptPath)
		if err != nil {
			return fmt.Sprintf("executeS3Script meet error=%v", err), err
		}
		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			logrus.Errorf("parseSaltApiCmdScriptCallResult meet err=%v,rawStr=%s", err, result)
			return fmt.Sprintf("parseSaltApiCmdScriptCallResult meet err=%v", err), err
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				return v.Stderr, fmt.Errorf("script run retCode=%v", v.RetCode)
			}
			output = v.Stdout + v.Stderr
			break
		}

	default:
		err = fmt.Errorf("no such EndPointType")
		logrus.Errorf("runScript meet error=%v", err)
		return fmt.Sprintf("runScript meet error=%v", err), err
	}

	return output, nil
}

func writeScriptContentToTempFile(content string) (fileName string, err error) {
	tmpFile, err := ioutil.TempFile(SCRIPT_SAVE_PATH, "script-")
	if err != nil {
		logrus.Errorf("writeScriptContentToTempFile,create tempfile meet err=%v", err)
		return
	}

	defer func() {
		if err != nil {
			defer os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write([]byte(content)); err != nil {
		logrus.Errorf("writeScriptContentToTempFile,write tempfile meet err=%v", err)
		return
	}

	if err = tmpFile.Close(); err != nil {
		logrus.Errorf("writeScriptContentToTempFile,close tempfile meet err=%v", err)
		return
	}

	fileName = tmpFile.Name()
	return
}

func (action *RunScriptAction) runScript(input *RunScriptInput) (output RunScriptOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Target = input.Target
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.RetCode = 0
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.RetCode = 1
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	scriptPath := input.EndPoint
	if input.EndPointType == END_POINT_TYPE_S3 {
		scriptPath, err = downLoadScript(*input)
		if err != nil {
			return output, err
		}
	}

	if input.EndPointType == END_POINT_TYPE_USER_PARAM {
		scriptPath, err = writeScriptContentToTempFile(input.ScriptContent)
		if err != nil {
			return output, err
		}
	}

	stdOut, err := runScript(scriptPath, *input)
	if err != nil {
		return output, err
	}
	output.Detail = stdOut

	return output, err
}

func (action *RunScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunScriptInputs)
	outputs := RunScriptOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output, err := action.runScript(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
