package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"os/exec"
	"time"
	"strings"
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
	ScriptPluginActions["ssh-run"] = new(SSHRunScriptAction)
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
	Password      string `json:"password,omitempty"`
	Seed          string `json:"seed,omitempty"`
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
	request.FullReturn = true

	logrus.Infof("executeLocalScript fileName=%s,target=%s,runAs=%s,execArgs=%s", fileName, target, runAs, execArg)

	fileAbsPath := fileName[:strings.LastIndex(fileName, "/")]
	if fileAbsPath == "" {
		fileAbsPath = "/"
	}
	fileShellName := fileName[strings.LastIndex(fileName, "/")+1:]
	//cmdRun := "/bin/bash " + fileName
	cmdRun := fmt.Sprintf("/bin/bash -c 'cd %s && ./%s", fileAbsPath, fileShellName)
	if len(execArg) > 0 {
		cmdRun = cmdRun + " " + execArg
	}
	cmdRun += "'"
	logrus.Infof("exec script to %s : %s ", target, cmdRun)
	request.Args = append(request.Args, cmdRun)

	if len(runAs) > 0 {
		request.Args = append(request.Args, "runas="+runAs)
	}
	logrus.Infof("executeLocalScript request=%v", request)

	result, err := CallSaltApi("https://127.0.0.1:8080", request)
	if err != nil {
		return "", err
	}
	logrus.Infof("executeLocalScript result: %++v", result)

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

func downLoadScript(input RunScriptInput) ([]string, error) {
	var result []string
	for _,v := range splitWithCustomFlag(input.EndPoint) {
		fileName, err := downloadS3File(v, DefaultS3Key, DefaultS3Password, false)
		if err != nil {
			logrus.Errorf("RunScriptAction downloads3 file:%s error=%v", v, err)
			return result, fmt.Errorf("RunScriptAction downloads3 file:%s error=%v", v, err)
		}

		scriptPath, err := saveFileToSaltMasterBaseDir(fileName)
		os.Remove(fileName)
		if err != nil {
			logrus.Errorf("saveFileToSaltMasterBaseDir file:%s meet error=%v", fileName, err)
			return result, fmt.Errorf("saveFileToSaltMasterBaseDir file:%s meet error=%v", fileName, err)
		}

		result = append(result, scriptPath)
	}

	return result, nil
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
			if v.RetCode != 0 {
				err = fmt.Errorf("script run ip[%s] meet error = %v", k, v.RetDetail)
				return k + ": " + v.RetDetail, err
			}
			output = k + ": " + v.RetDetail
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

	//scriptPath := input.EndPoint
	scriptPathList := splitWithCustomFlag(input.EndPoint)
	if input.EndPointType == END_POINT_TYPE_S3 {
		scriptPathList, err = downLoadScript(*input)
		if err != nil {
			return output, err
		}
	}

	if input.EndPointType == END_POINT_TYPE_USER_PARAM {
		var scriptPath string
		scriptPath, err = writeScriptContentToTempFile(input.ScriptContent)
		if err != nil {
			return output, err
		}
		scriptPathList = append(scriptPathList, scriptPath)
	}

	var stdOut string
	for i,v := range scriptPathList {
		stdOut, err = runScript(v, *input)
		stdOut = fmt.Sprintf("script %d result: %s ", i+1, stdOut)
		output.Detail += stdOut
		if err != nil {
			logrus.Errorf(stdOut)
			err = fmt.Errorf(stdOut)
			return output, err
		}else{
			logrus.Infof("%s success ", stdOut)
		}
	}

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

type SSHRunScriptAction struct {}

func (action *SSHRunScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *SSHRunScriptAction) CheckParam(input RunScriptInput) error {
	if input.EndPointType != END_POINT_TYPE_LOCAL && input.EndPointType != END_POINT_TYPE_S3 && input.EndPointType != END_POINT_TYPE_USER_PARAM {
		return errors.New("Wrong EndPointType")
	}
	if input.EndPoint == "" && input.EndPointType == END_POINT_TYPE_S3 {
		return errors.New("Endpoint is empty")
	}
	if input.Target == "" {
		return errors.New("Target is empty")
	}
	if input.Password == "" {
		return errors.New("Password is empty")
	}

	return nil
}

func sshRunScript(scriptPath string, input RunScriptInput) (string, error) {
	var output string
	//var cmdOut []byte
	var err error
	if input.RunAs == "" {
		input.RunAs = "root"
	}
	remoteParam := ExecRemoteParam{User:input.RunAs,Password:input.Password,Host:input.Target,Timeout:300}
	switch input.EndPointType {
	case END_POINT_TYPE_LOCAL:
		localAbsPath := scriptPath[:strings.LastIndex(scriptPath, "/")]
		if localAbsPath == "" {
			localAbsPath = "/"
		}
		localShellPath := scriptPath[strings.LastIndex(scriptPath, "/")+1:]
		remoteParam.Command = fmt.Sprintf("/bin/bash cd %s && ./%s", localAbsPath, localShellPath)
		if len(input.ExecArg) > 0 {
			remoteParam.Command = remoteParam.Command + " " + input.ExecArg
		}
		logrus.Infof("ssh run script local: %s", remoteParam.Command)
		execRemoteWithTimeout(&remoteParam)
		//cmdOut,err = execRemote(input.RunAs, input.Password, input.Target, fmt.Sprintf("bash %s", scriptPath))
		err = remoteParam.Err
		output = remoteParam.Output
		logrus.Infof("exec ssh script:%s in target:%s output:%s \n", scriptPath, input.Target, output)
		if err != nil {
			return fmt.Sprintf("exec ssh to run the script:%s in %s,output:%s ,meet error=%v", scriptPath, input.Target, output, err), err
		}
	case END_POINT_TYPE_S3, END_POINT_TYPE_USER_PARAM:
		newScriptName := fmt.Sprintf("ssh-script-%s-%d", strings.Replace(input.Target, ".", "-", -1), time.Now().Unix())
		err = exec.Command("/bin/cp", "-f", scriptPath, fmt.Sprintf("/var/www/html/tmp/%s", newScriptName)).Run()
		if err != nil {
			return fmt.Sprintf("exec ssh script,cp %s %s meet error=%v", scriptPath, newScriptName, err), err
		}
		err = exec.Command("bash", "-c", fmt.Sprintf("chmod 666 /var/www/html/tmp/%s", newScriptName)).Run()
		if err != nil {
			return fmt.Sprintf("exec ssh script,chmod to %s meet error=%v", newScriptName, err), err
		}
		remoteParam.Command = fmt.Sprintf("curl http://%s:9099/tmp/%s | bash ", MasterHostIp, newScriptName)
		execRemoteWithTimeout(&remoteParam)
		err = remoteParam.Err
		//cmdOut,err = execRemote(input.RunAs, input.Password, input.Target, fmt.Sprintf("curl http://%s:9099/tmp/%s | bash ", MasterHostIp, newScriptName))
		os.Remove(scriptPath)
		os.Remove(fmt.Sprintf("/var/www/html/tmp/%s", newScriptName))
		output = remoteParam.Output
		logrus.Infof("exec ssh script:%s ,target:%s output:%s \n",fmt.Sprintf("curl http://%s:9099/tmp/%s | bash ", MasterHostIp, newScriptName), input.Target, output)
		if err != nil {
			return fmt.Sprintf("exec ssh script error,target:%s output:%s error:%v", input.Target, output, err),err
		}
	default:
		err = fmt.Errorf("no such EndPointType")
		logrus.Errorf("runScript meet error=%v", err)
		return fmt.Sprintf("runScript meet error=%v", err), err
	}

	return output, nil
}

func (action *SSHRunScriptAction) runScript(input *RunScriptInput) (output RunScriptOutput, err error) {
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

	//scriptPath := input.EndPoint
	scriptPathList := splitWithCustomFlag(input.EndPoint)
	if input.EndPointType == END_POINT_TYPE_S3 {
		scriptPathList, err = downLoadScript(*input)
		if err != nil {
			return output, err
		}
	}

	if input.EndPointType == END_POINT_TYPE_USER_PARAM {
		var scriptPath string
		scriptPath, err = writeScriptContentToTempFile(input.ScriptContent)
		if err != nil {
			return output, err
		}
		scriptPathList = append(scriptPathList, scriptPath)
	}

	input.Password,_ = AesDePassword(input.Guid, input.Seed, input.Password)

	var stdOut string
	for i,v := range scriptPathList {
		stdOut, err = sshRunScript(v, *input)
		stdOut = fmt.Sprintf("script %d result: %s ", i+1, stdOut)
		output.Detail += stdOut
		if err != nil {
			logrus.Errorf(stdOut)
			err = fmt.Errorf(stdOut)
			return output, err
		}else{
			logrus.Infof("%s success ", stdOut)
		}
	}

	return output, err
}

func (action *SSHRunScriptAction) Do(input interface{}) (interface{}, error) {
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