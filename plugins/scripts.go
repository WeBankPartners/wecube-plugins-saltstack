package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"os/exec"
	"time"
	"strings"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
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
	Language string
}

func (action *RunScriptAction) SetAcceptLanguage(language string) {
	action.Language = language
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
		return getParamValidateError(action.Language, "endpointType", fmt.Sprintf("must in (%s,%s,%s)",END_POINT_TYPE_LOCAL,END_POINT_TYPE_S3,END_POINT_TYPE_USER_PARAM))
	}
	if input.EndPoint == "" && input.EndPointType == END_POINT_TYPE_S3 {
		return getParamValidateError(action.Language, "endpoint", "endpoint cat not empty when endpointType="+END_POINT_TYPE_S3)
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.RunAs == "" {
		return getParamEmptyError(action.Language, "runAs")
	}

	return nil
}

// why not move? TODO
func saveFileToSaltMasterBaseDir(fileName string) (string, error) {
	var err error
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("read %s fail,%s", fileName, err.Error())
	}

	tmpFile, err := ioutil.TempFile(SCRIPT_SAVE_PATH, "script-")
	if err != nil {
		return "", fmt.Errorf("create tmp file fail,%s", err.Error())
	}

	defer func() {
		if err != nil {
			defer os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write(content); err != nil {
		return "", fmt.Errorf("write content to tmp file fail,%s", err.Error())
	}

	if err = tmpFile.Close(); err != nil {
		return "", fmt.Errorf("close tmp file fail,%s", err.Error())
	}

	fullPath := tmpFile.Name()
	return fullPath, err
}

func executeS3Script(fileName string, target string, runAs string, execArg string, language string) (string, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.script"

	request.Args = append(request.Args, "salt://base/"+fileName)
	if len(execArg) > 0 {
		request.Args = append(request.Args, "args="+execArg)
	}

	if len(runAs) > 0 {
		request.Args = append(request.Args, "runas="+runAs)
	}

	result, err := CallSaltApi("https://127.0.0.1:8080", request, language)
	if err != nil {
		return "", err
	}

	return result, nil
}

func executeLocalScript(fileName string, target string, runAs string, execArg string, language string) (string, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.run"
	request.FullReturn = true

	log.Logger.Info("Exec local script", log.String("fileName",fileName), log.String("target",target), log.String("runAs",runAs), log.String("args", execArg))

	fileAbsPath := fileName[:strings.LastIndex(fileName, "/")]
	if fileAbsPath == "" {
		fileAbsPath = "/"
	}
	fileShellName := fileName[strings.LastIndex(fileName, "/")+1:]
	//cmdRun := "/bin/bash " + fileName
	fileShellName = strings.TrimSpace(fileShellName)
	tmpFileShellBin := fileShellName
	if strings.Contains(tmpFileShellBin, " ") {
		tmpFileShellBin = strings.Split(tmpFileShellBin, " ")[0]
	}
	cmdRun := fmt.Sprintf("/bin/bash -c 'cd %s && chmod +x %s && ./%s", fileAbsPath, tmpFileShellBin, fileShellName)
	if len(execArg) > 0 {
		cmdRun = cmdRun + " " + execArg
	}
	cmdRun += "'"
	request.Args = append(request.Args, cmdRun)

	if len(runAs) > 0 {
		request.Args = append(request.Args, "runas="+runAs)
	}
	log.Logger.Debug("Exec script", log.String("target", target), log.JsonObj("request", request))

	result, err := CallSaltApi("https://127.0.0.1:8080", request, language)
	if err != nil {
		return "", err
	}

	return result, nil
}

func checkRunUserIsExists(target,userGroup,language string) (exist bool,output string) {
	if userGroup == "" {
		return false,"user is empty"
	}
	exist = false
	var user,group string
	if strings.Contains(userGroup, ":") {
		tmpList := strings.Split(userGroup, ":")
		user = tmpList[0]
		group = tmpList[1]
	}else{
		user = userGroup
	}
	log.Logger.Info("Check user group if exist", log.String("user", user), log.String("group",group))
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.run"
	request.FullReturn = true
	cmdRun := fmt.Sprintf("/bin/bash -c 'cat /etc/passwd|grep %s:'", user)
	request.Args = append(request.Args, cmdRun)
	result, err := CallSaltApi("https://127.0.0.1:8080", request, language)
	if err != nil {
		log.Logger.Error("Check user exists,call salt api error", log.Error(err))
		return false,fmt.Sprintf("check user exists,call salt api error,%s", err.Error())
	}
	saltApiResult, err := parseSaltApiCmdRunCallResult(result)
	if err != nil {
		log.Logger.Error("check user exists,parse salt api result error", log.Error(err))
		return false,fmt.Sprintf("check user exists,parse salt api result error,%s", err.Error())
	}
	for _, v := range saltApiResult.Results[0] {
		if strings.Contains(v.RetDetail, user) {
			exist = true
		}
	}
	if !exist {
		return false,fmt.Sprintf("user %s not exist", user)
	}
	if group != "" {
		exist = false
		groupRequest := SaltApiRequest{}
		groupRequest.Client = "local"
		groupRequest.TargetType = "ipcidr"
		groupRequest.Target = target
		groupRequest.Function = "cmd.run"
		groupRequest.FullReturn = true
		cmdRun := fmt.Sprintf("/bin/bash -c 'cat /etc/group|grep %s:'", group)
		groupRequest.Args = append(groupRequest.Args, cmdRun)
		result, err := CallSaltApi("https://127.0.0.1:8080", groupRequest, language)
		if err != nil {
			log.Logger.Error("check group exists,call salt api error", log.Error(err))
			return false,fmt.Sprintf("check group exists,call salt api error,%s", err.Error())
		}
		saltApiResult, err := parseSaltApiCmdRunCallResult(result)
		if err != nil {
			log.Logger.Error("check group exists,parse salt api result error", log.Error(err))
			return false,fmt.Sprintf("check group exists,parse salt api result error,%s", err.Error())
		}
		for _, v := range saltApiResult.Results[0] {
			if strings.Contains(v.RetDetail, group) {
				exist = true
			}
		}
		if !exist {
			return false,fmt.Sprintf("group %s not exist", group)
		}
	}
	return true,""
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

func downLoadScript(input RunScriptInput, language string) ([]string, error) {
	var result []string
	for _,v := range splitWithCustomFlag(input.EndPoint) {
		fileName, err := downloadS3File(v, DefaultS3Key, DefaultS3Password, false, language)
		if err != nil {
			return result, err
		}

		scriptPath, err := saveFileToSaltMasterBaseDir(fileName)
		os.Remove(fileName)
		if err != nil {
			return result, fmt.Errorf("Move file:%s to salt-dir fail,%v ", fileName, err)
		}

		result = append(result, scriptPath)
	}

	return result, nil
}

func runScript(scriptPath string, input RunScriptInput, language string) (string, error) {
	var result string
	var output string
	var err error
	if strings.Contains(input.RunAs, ":") {
		input.RunAs = strings.Split(input.RunAs, ":")[0]
	}
	switch input.EndPointType {
	case END_POINT_TYPE_LOCAL:
		result, err = executeLocalScript(scriptPath, input.Target, input.RunAs, input.ExecArg, language)
		if err != nil {
			return "", getRunRemoteScriptError(language, input.Target, result, err)
		}
		saltApiResult, err := parseSaltApiCmdRunCallResult(result)
		if err != nil {
			return "", fmt.Errorf("Parse salt call result fail,%s ", err.Error())
		}
		for k, v := range saltApiResult.Results[0] {
			if k != input.Target {
				err = fmt.Errorf("Script run ip[%s] is not a available target[%s] ", k, input.Target)
				return fmt.Sprintf("parseSaltApiCmdRunCallResult meet error=%v", err), err
			}
			if v.RetCode != 0 || strings.Contains(v.RetDetail, "ERROR") {
				err = fmt.Errorf("Script run ip[%s] meet error = %v ", k, v.RetDetail)
				return k + ": " + v.RetDetail, err
			}
			output = k + ": " + v.RetDetail
			break
		}
	case END_POINT_TYPE_S3, END_POINT_TYPE_USER_PARAM:
		result, err = executeS3Script(filepath.Base(scriptPath), input.Target, input.RunAs, input.ExecArg, language)
		os.Remove(scriptPath)
		if err != nil {
			return "", getRunRemoteScriptError(language, input.Target, result, err)
		}
		saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
		if err != nil {
			return "", fmt.Errorf("Parse salt call result fail,%s ", err.Error())
		}

		for _, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 || strings.Contains(v.Stdout, "ERROR") {
				return v.Stderr, fmt.Errorf("Script run retCode=%d output=%s error=%s ", v.RetCode, v.Stdout, v.Stderr)
			}
			output = v.Stdout + v.Stderr
			break
		}

	default:
		err = fmt.Errorf("No such endPointType ")
		return "", err
	}

	return output, nil
}

func writeScriptContentToTempFile(content string) (fileName string, err error) {
	tmpFile, err := ioutil.TempFile(SCRIPT_SAVE_PATH, "script-")
	if err != nil {
		err = fmt.Errorf("New tmp file error,%s ", err.Error())
		return fileName,err
	}

	defer func() {
		if err != nil {
			defer os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write([]byte(content)); err != nil {
		err = fmt.Errorf("Write script content to tmp file error,%s ", err.Error())
		return fileName,err
	}

	if err = tmpFile.Close(); err != nil {
		err = fmt.Errorf("Tmp file close fail,%s ", err.Error())
		return fileName,err
	}

	fileName = tmpFile.Name()
	return fileName,err
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

	if input.RunAs != "" {
		if strings.Contains(input.RunAs, ":") {
			input.RunAs = strings.Split(input.RunAs, ":")[0]
		}
		userExist,errOut := checkRunUserIsExists(input.Target, input.RunAs, action.Language)
		if !userExist {
			err = fmt.Errorf(errOut)
			return output,err
		}
	}

	//scriptPath := input.EndPoint
	scriptPathList := splitWithCustomFlag(input.EndPoint)
	if input.EndPointType == END_POINT_TYPE_S3 {
		scriptPathList, err = downLoadScript(*input, action.Language)
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
		scriptPathList = []string{scriptPath}
	}

	var stdOut string
	for i,v := range scriptPathList {
		stdOut, err = runScript(v, *input, action.Language)
		stdOut = fmt.Sprintf("script %d result: %s ", i+1, stdOut)
		if i < len(scriptPathList)-1 {
			stdOut += " | "
		}
		output.Detail += stdOut
		if err != nil {
			log.Logger.Error("Run script fail", log.String("output", stdOut), log.Error(err))
			return output, err
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
			log.Logger.Error("Run script action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

type SSHRunScriptAction struct { Language string }

func (action *SSHRunScriptAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *SSHRunScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *SSHRunScriptAction) CheckParam(input RunScriptInput) error {
	if input.EndPointType != END_POINT_TYPE_LOCAL && input.EndPointType != END_POINT_TYPE_S3 && input.EndPointType != END_POINT_TYPE_USER_PARAM {
		return getParamValidateError(action.Language, "endpointType", fmt.Sprintf("must in (%s,%s,%s)",END_POINT_TYPE_LOCAL,END_POINT_TYPE_S3,END_POINT_TYPE_USER_PARAM))
	}
	if input.EndPoint == "" && input.EndPointType == END_POINT_TYPE_S3 {
		return getParamValidateError(action.Language, "endpoint", "endpoint cat not empty when endpointType="+END_POINT_TYPE_S3)
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.RunAs == "" {
		return getParamEmptyError(action.Language, "runAs")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	return nil
}

func sshRunScript(scriptPath string, input RunScriptInput, language string) (string, error) {
	var output string
	//var cmdOut []byte
	var err error
	if strings.Contains(input.RunAs, ":") {
		input.RunAs = strings.Split(input.RunAs, ":")[0]
	}
	remoteParam := ExecRemoteParam{User:input.RunAs,Password:input.Password,Host:input.Target,Timeout:1800}
	switch input.EndPointType {
	case END_POINT_TYPE_LOCAL:
		localAbsPath := scriptPath[:strings.LastIndex(scriptPath, "/")]
		if localAbsPath == "" {
			localAbsPath = "/"
		}
		localShellPath := scriptPath[strings.LastIndex(scriptPath, "/")+1:]
		localShellPath = strings.TrimSpace(localShellPath)
		tmpFileShellBin := localShellPath
		if strings.Contains(tmpFileShellBin, " ") {
			tmpFileShellBin = strings.Split(tmpFileShellBin, " ")[0]
		}
		remoteParam.Command = fmt.Sprintf("/bin/bash -c 'cd %s && chmod +x %s && ./%s", localAbsPath, tmpFileShellBin, localShellPath)
		if len(input.ExecArg) > 0 {
			remoteParam.Command = remoteParam.Command + " " + input.ExecArg
		}
		remoteParam.Command += "'"
		execRemoteWithTimeout(&remoteParam)
		err = remoteParam.Err
		output = remoteParam.Output
		log.Logger.Debug("SSH run script", log.String("type",input.EndPointType), log.String("command", remoteParam.Command), log.String("target",input.Target), log.String("output", output), log.Error(err))
		if err != nil {
			return "", getRunRemoteScriptError(language, input.Target, output, err)
		}
	case END_POINT_TYPE_S3, END_POINT_TYPE_USER_PARAM:
		newScriptName := fmt.Sprintf("ssh-script-%s-%d", strings.Replace(input.Target, ".", "-", -1), time.Now().Unix())
		tmpOut,err := exec.Command("/bin/cp", "-f", scriptPath, fmt.Sprintf("/var/www/html/tmp/%s", newScriptName)).Output()
		if err != nil {
			return fmt.Sprintf("exec ssh script,cp %s %s meet error output=%s,err=%s", scriptPath, newScriptName, string(tmpOut), err.Error()), err
		}
		tmpOut,err = exec.Command("bash", "-c", fmt.Sprintf("chmod 666 /var/www/html/tmp/%s", newScriptName)).Output()
		if err != nil {
			return fmt.Sprintf("exec ssh script,chmod to %s meet output=%s,err=%s", newScriptName, string(tmpOut), err.Error()), err
		}
		remoteParam.Command = fmt.Sprintf("curl http://%s:9099/tmp/%s | bash ", MasterHostIp, newScriptName)
		execRemoteWithTimeout(&remoteParam)
		err = remoteParam.Err
		os.Remove(scriptPath)
		os.Remove(fmt.Sprintf("/var/www/html/tmp/%s", newScriptName))
		output = remoteParam.Output
		log.Logger.Debug("SSH run script", log.String("type",input.EndPointType), log.String("script", newScriptName), log.String("target",input.Target), log.String("output", output), log.Error(err))
		if err != nil {
			return "", getRunRemoteScriptError(language, input.Target, output, err)
		}
	default:
		err = fmt.Errorf("No such endPointType ")
		return "", err
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
		scriptPathList, err = downLoadScript(*input, action.Language)
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
		scriptPathList = []string{scriptPath}
	}

	input.Password,err = AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output,err
	}

	var stdOut string
	for i,v := range scriptPathList {
		stdOut, err = sshRunScript(v, *input, action.Language)
		stdOut = fmt.Sprintf("script %d result: %s ", i+1, stdOut)
		output.Detail += stdOut
		if err != nil {
			log.Logger.Error("Run ssh script", log.String("output", stdOut), log.Error(err))
			return output, err
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
