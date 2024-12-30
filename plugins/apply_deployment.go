package plugins

import (
	"fmt"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"sync"
)

var ApplyDeploymentActions = make(map[string]Action)

func init() {
	ApplyDeploymentActions["new"] = new(ApplyNewDeploymentAction)
	ApplyDeploymentActions["update"] = new(ApplyUpdateDeploymentAction)
	ApplyDeploymentActions["delete"] = new(ApplyDeleteDeploymentAction)
	//ApplyDeploymentActions["rollback"] = new(ApplyRollbackDeploymentAction)
}

type ApplyDeploymentPlugin struct {
}

func (plugin *ApplyDeploymentPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ApplyDeploymentActions[actionName]
	if !found {
		return nil, fmt.Errorf("ApplyDeployment plugin,action = %s not found", actionName)
	}

	return action, nil
}

type ApplyNewDeploymentInputs struct {
	Inputs []ApplyNewDeploymentInput `json:"inputs,omitempty"`
}
type ApplyNewDeploymentInput struct {
	CallBackParameter
	EndPoint         string `json:"endpoint,omitempty"`
	Guid             string `json:"guid,omitempty"`
	UserName         string `json:"userName,omitempty"`
	Target           string `json:"target,omitempty"`
	DestinationPath  string `json:"destinationPath,omitempty"`
	VariableFilePath string `json:"confFiles,omitempty"`
	VariableList     string `json:"variableList,omitempty"`
	ExecArg          string `json:"args,omitempty"`
	StartScriptPath  string `json:"startScript,omitempty"`
	// AccessKey    string `json:"accessKey,omitempty"`
	// SecretKey    string `json:"secretKey,omitempty"`
	// EncryptVariblePrefix string `json:"encryptVariblePrefix,omitempty"`
	Seed          string `json:"seed,omitempty"`
	AppPublicKey  string `json:"appPublicKey,omitempty"`
	SysPrivateKey string `json:"sysPrivateKey,omitempty"`
	Password      string `json:"password,omitempty"`
	RwDir         string `json:"rwDir,omitempty"`
	RwFile        string `json:"rwFile,omitempty"`

	LogFileTrade   string `json:"logFileTrade,omitempty"`
	LogFileTrace   string `json:"logFileTrace,omitempty"`
	LogFileMetric  string `json:"logFileMetric,omitempty"`
	LogFileKeyword string `json:"logFileKeyword,omitempty"`

	SignFileSrc string `json:"SignFileSrc,omitempty"`
	SignFileDst string `json:"SignFileDst,omitempty"`
}

type ApplyNewDeploymentOutputs struct {
	Outputs []ApplyNewDeploymentOutput `json:"outputs,omitempty"`
}

type ApplyNewDeploymentOutput struct {
	CallBackParameter
	Result
	Guid            string `json:"guid,omitempty"`
	UserDetail      string `json:"userDetail,omitempty"`
	FileDetail      string `json:"fileDetail,omitempty"`
	NewS3PkgPath    string `json:"s3PkgPath,omitempty"`
	Target          string `json:"target,omitempty"`
	RetCode         int    `json:"retCode,omitempty"`
	RunScriptDetail string `json:"runScriptDetail,omitempty"`
	Password        string `json:"password,omitempty"`

	LogFileTrade   string `json:"logFileTrade,omitempty"`
	LogFileTrace   string `json:"logFileTrace,omitempty"`
	LogFileMetric  string `json:"logFileMetric,omitempty"`
	LogFileKeyword string `json:"logFileKeyword,omitempty"`
}

type ApplyNewDeploymentThreadObj struct {
	Data  ApplyNewDeploymentOutput
	Err   error
	Index int
}

type ApplyNewDeploymentAction struct {
	Language string
}

func (action *ApplyNewDeploymentAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *ApplyNewDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyNewDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ApplyNewDeploymentAction) CheckParam(input ApplyNewDeploymentInput) error {
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.StartScriptPath == "" {
		return getParamEmptyError(action.Language, "startScript")
	}
	if input.DestinationPath == "" {
		return getParamEmptyError(action.Language, "destinationPath")
	}

	return nil
}

func (action *ApplyNewDeploymentAction) applyNewDeployment(input ApplyNewDeploymentInput) (output ApplyNewDeploymentOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Target = input.Target
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.RetCode = 1
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(input)
	if err != nil {
		return output, err
	}
	input.Seed = getEncryptSeed(input.Seed)
	userGroup := input.UserName
	if strings.Contains(userGroup, ":") {
		userGroup = strings.Split(userGroup, ":")[1]
	}

	if !strings.Contains(input.UserName, ":") {
		input.UserName = fmt.Sprintf("%s:%s", input.UserName, input.UserName)
	}
	if !strings.HasPrefix(input.StartScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newScriptPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StartScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript[1:])
			} else {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript)
			}
		}
		input.StartScriptPath = strings.Join(newScriptPathList, ",")
	}
	//create apply deployment user
	addUserRequest := AddUserInputs{
		Inputs: []AddUserInput{
			AddUserInput{
				Guid:      input.Guid,
				Target:    input.Target,
				UserName:  strings.Split(input.UserName, ":")[0],
				Password:  input.Password,
				Seed:      input.Seed,
				RwDir:     input.RwDir,
				RwFile:    input.RwFile,
				UserGroup: userGroup,
			},
		},
	}

	log.Logger.Debug("App deploy", log.String("step", "create user"), log.JsonObj("param", addUserRequest))
	userAddOutputs, err := createApplyUser(addUserRequest)
	if err != nil {
		return output, fmt.Errorf("Create user fail,%s ", err.Error())
	}
	output.UserDetail = userAddOutputs.(*AddUserOutputs).Outputs[0].Detail
	output.Password = userAddOutputs.(*AddUserOutputs).Outputs[0].Password

	// replace apply variable
	var variableReplaceOutputs interface{}
	if input.VariableFilePath != "" {
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:         input.Guid,
					EndPoint:     input.EndPoint,
					FilePath:     input.VariableFilePath,
					VariableList: input.VariableList,
					// EncryptVariblePrefix: input.EncryptVariblePrefix,
					Seed:          input.Seed,
					AppPublicKey:  input.AppPublicKey,
					SysPrivateKey: input.SysPrivateKey,
				},
			},
		}

		log.Logger.Debug("App deploy", log.String("step", "replace variable"), log.JsonObj("param", variableReplaceRequest))
		variableReplaceOutputs, err = replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			return output, fmt.Errorf("Replace variable fail,%s ", err.Error())
		}
		output.NewS3PkgPath = variableReplaceOutputs.(*VariableReplaceOutputs).Outputs[0].NewS3PkgPath
	} else {
		output.NewS3PkgPath = input.EndPoint
	}

	// copy apply package
	fileCopyRequest := FileCopyInputs{
		Inputs: []FileCopyInput{
			FileCopyInput{
				EndPoint:        output.NewS3PkgPath,
				Guid:            input.Guid,
				Target:          input.Target,
				DestinationPath: input.DestinationPath,
				Unpack:          "true",
				FileOwner:       input.UserName,
			},
		},
	}

	log.Logger.Debug("App deploy", log.String("step", "file copy"), log.JsonObj("param", fileCopyRequest))
	fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
	if err != nil {
		return output, fmt.Errorf("Copy app package to target fail,%s ", err.Error())
	}
	output.FileDetail = fileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail

	// copy apply package's sign file from s3 to dst
	if input.SignFileSrc != "" && input.SignFileDst != "" {
		signFileCopyRequest := FileCopyInputs{
			Inputs: []FileCopyInput{
				{
					EndPoint:        input.SignFileSrc,
					Guid:            input.Guid,
					Target:          input.Target,
					DestinationPath: input.SignFileDst,
					Unpack:          "false",
					FileOwner:       input.UserName,
				},
			},
		}

		log.Logger.Debug("App deploy", log.String("step", "sign file copy"), log.JsonObj("param", signFileCopyRequest))
		signFileCopyOutputs, err := copyApplyFile(signFileCopyRequest)
		if err != nil {
			return output, fmt.Errorf("Copy app package's sign file to target fail,%s ", err.Error())
		}
		output.FileDetail = signFileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail
	}

	// start apply script
	runScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StartScriptPath,
				Target:       input.Target,
				RunAs:        strings.Split(input.UserName, ":")[0],
				Guid:         input.Guid,
			},
		},
	}
	if input.ExecArg != "" {
		runScriptRequest.Inputs[0].ExecArg = input.ExecArg
	}

	log.Logger.Debug("App deploy", log.String("step", "run script"), log.JsonObj("param", runScriptRequest))
	runScriptOutputs, err := runApplyScript(runScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run start script fail,%s ", err.Error())
	}
	output.RunScriptDetail = runScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	// parse log files from glob name
	action.parseGlobLogFiles(&input, &output)
	return output, err
}

func (action *ApplyNewDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyNewDeploymentInputs)
	outputs := ApplyNewDeploymentOutputs{}
	var finalErr error
	outputChan := make(chan ApplyNewDeploymentThreadObj, len(inputs.Inputs))
	wg := sync.WaitGroup{}
	for i, input := range inputs.Inputs {
		wg.Add(1)
		go func(tmpInput ApplyNewDeploymentInput, index int) {
			output, err := action.applyNewDeployment(tmpInput)
			outputChan <- ApplyNewDeploymentThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
		}(input, i)
		outputs.Outputs = append(outputs.Outputs, ApplyNewDeploymentOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("App new deploy action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}
	return &outputs, finalErr
}

func (action *ApplyNewDeploymentAction) parseGlobLogFiles(input *ApplyNewDeploymentInput, output *ApplyNewDeploymentOutput) {
	for _, logFile := range []struct {
		logFilePattern string
		outputField    *string
	}{
		{input.LogFileTrade, &output.LogFileTrade},
		{input.LogFileTrace, &output.LogFileTrace},
		{input.LogFileMetric, &output.LogFileMetric},
		{input.LogFileKeyword, &output.LogFileKeyword},
	} {
		if logFile.logFilePattern == "" {
			continue
		}
		log.Logger.Debug("App deploy", log.String("step", "parse log file: "+logFile.logFilePattern))
		if ret, err := FindGlobFiles(input.DestinationPath, logFile.logFilePattern, input.Target); err == nil {
			*logFile.outputField = ret
		}
	}
}

type ApplyUpdateDeploymentInputs struct {
	Inputs []ApplyUpdateDeploymentInput `json:"inputs,omitempty"`
}

type ApplyUpdateDeploymentInput struct {
	CallBackParameter
	EndPoint         string `json:"endpoint,omitempty"`
	Guid             string `json:"guid,omitempty"`
	UserName         string `json:"userName,omitempty"`
	Target           string `json:"target,omitempty"`
	DestinationPath  string `json:"destinationPath,omitempty"`
	VariableFilePath string `json:"confFiles,omitempty"`
	VariableList     string `json:"variableList,omitempty"`
	ExecArg          string `json:"args,omitempty"`
	StopScriptPath   string `json:"stopScript,omitempty"`
	StartScriptPath  string `json:"startScript,omitempty"`

	// EncryptVariblePrefix string `json:"encryptVariblePrefix,omitempty"`
	Seed          string `json:"seed,omitempty"`
	AppPublicKey  string `json:"appPublicKey,omitempty"`
	SysPrivateKey string `json:"sysPrivateKey,omitempty"`

	LogFileTrade   string `json:"logFileTrade,omitempty"`
	LogFileTrace   string `json:"logFileTrace,omitempty"`
	LogFileMetric  string `json:"logFileMetric,omitempty"`
	LogFileKeyword string `json:"logFileKeyword,omitempty"`

	//AppBackUpEnabled string `json:"appBackUpEnabled,omitempty"`
	//AppBackUpPath    string `json:"appBackUpPath,omitempty"`
	//ExcludePath      string `json:"excludePath,omitempty"`

	SignFileSrc string `json:"SignFileSrc,omitempty"`
	SignFileDst string `json:"SignFileDst,omitempty"`
}

type ApplyUpdateDeploymentOutputs struct {
	Outputs []ApplyUpdateDeploymentOutput `json:"outputs,omitempty"`
}

type ApplyUpdateDeploymentOutput struct {
	CallBackParameter
	Result
	Guid                 string `json:"guid,omitempty"`
	FileDetail           string `json:"fileDetail,omitempty"`
	NewS3PkgPath         string `json:"s3PkgPath,omitempty"`
	Target               string `json:"target,omitempty"`
	RetCode              int    `json:"retCode,omitempty"`
	RunStartScriptDetail string `json:"runStartScriptDetail,omitempty"`
	RunStopScriptDetail  string `json:"runStopScriptDetail,omitempty"`

	LogFileTrade   string `json:"logFileTrade,omitempty"`
	LogFileTrace   string `json:"logFileTrace,omitempty"`
	LogFileMetric  string `json:"logFileMetric,omitempty"`
	LogFileKeyword string `json:"logFileKeyword,omitempty"`
}

type ApplyUpdateDeploymentAction struct {
	Language string
}

type ApplyUpdateDeploymentThreadObj struct {
	Data  ApplyUpdateDeploymentOutput
	Err   error
	Index int
}

func (action *ApplyUpdateDeploymentAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *ApplyUpdateDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyUpdateDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ApplyUpdateDeploymentAction) CheckParam(input ApplyUpdateDeploymentInput) error {
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.StartScriptPath == "" {
		return getParamEmptyError(action.Language, "startScriptPath")
	}
	if input.DestinationPath == "" {
		return getParamEmptyError(action.Language, "destinationPath")
	}
	if input.StopScriptPath == "" {
		return getParamEmptyError(action.Language, "stopScriptPath")
	}

	return nil
}

func (action *ApplyUpdateDeploymentAction) applyUpdateDeployment(input ApplyUpdateDeploymentInput) (output ApplyUpdateDeploymentOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Target = input.Target
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.RetCode = 1
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(input)
	if err != nil {
		return output, err
	}
	input.Seed = getEncryptSeed(input.Seed)
	if !strings.Contains(input.UserName, ":") {
		input.UserName = fmt.Sprintf("%s:%s", input.UserName, input.UserName)
	}
	if !strings.HasPrefix(input.StartScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newScriptPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StartScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript[1:])
			} else {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript)
			}
		}
		input.StartScriptPath = strings.Join(newScriptPathList, ",")
	}
	if !strings.HasPrefix(input.StopScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newStopPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StopScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript[1:])
			} else {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript)
			}
		}
		input.StopScriptPath = strings.Join(newStopPathList, ",")
	}
	// stop apply script
	runStopScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StopScriptPath,
				Target:       input.Target,
				RunAs:        strings.Split(input.UserName, ":")[0],
				Guid:         input.Guid,
			},
		},
	}

	log.Logger.Debug("App update", log.String("step", "run stop script"), log.JsonObj("param", runStopScriptRequest))
	runStopScriptOutputs, err := runApplyScript(runStopScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run stop script fail,%s ", err.Error())
	}
	output.RunStopScriptDetail = runStopScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	// replace apply variable
	var variableReplaceOutputs interface{}
	if input.VariableFilePath != "" {
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:         input.Guid,
					EndPoint:     input.EndPoint,
					FilePath:     input.VariableFilePath,
					VariableList: input.VariableList,
					// EncryptVariblePrefix: input.EncryptVariblePrefix,
					Seed:          input.Seed,
					AppPublicKey:  input.AppPublicKey,
					SysPrivateKey: input.SysPrivateKey,
				},
			},
		}

		log.Logger.Debug("App update", log.String("step", "variable replace"), log.JsonObj("param", variableReplaceRequest))
		variableReplaceOutputs, err = replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			return output, fmt.Errorf("Replace variable fail,%s ", err.Error())
		}
		output.NewS3PkgPath = variableReplaceOutputs.(*VariableReplaceOutputs).Outputs[0].NewS3PkgPath
	} else {
		output.NewS3PkgPath = input.EndPoint
	}

	// backup dest dir to tar guid.tar.gz
	// salt '*' archive.tar zxf {{.AppBackUpPath}}/{{.guid}}.tar.gz dest='{{.SourcePath}}'
	//if input.AppBackUpEnabled == "Y" && input.AppBackUpPath != "" {
	//	log.Logger.Debug("App update", log.String("step", "backup"), log.JsonObj("input", input))
	//
	//	// salt '*' file.mkdir {{.AppBackUpPath}}
	//	if _, saltErr := CallSaltApi("https://127.0.0.1:8080", SaltApiRequest{
	//		Client:   "local",
	//		Function: "file.mkdir",
	//		Target:   input.Target,
	//		Args:     []string{input.AppBackUpPath},
	//	}, action.Language); saltErr != nil {
	//		errMsg := fmt.Sprintf("Failed when mkdir [%s], %s\n", input.AppBackUpPath, saltErr.Error())
	//		output.FileDetail = errMsg
	//		return output, fmt.Errorf(errMsg)
	//	}
	//
	//	// salt '*' archive.tar zcf {{.AppBackUpPath}}/{{.guid}}.tar.gz cwd='{{.SourcePath}}' .
	//	// tar -zcvf tmp.tar.gz --exclude-vcs --exclude={.idea,venv,vendor,.github} .
	//	cmdArgs := []string{
	//		"zcf",
	//		path.Join(input.AppBackUpPath, fmt.Sprintf("%s.tar.gz", input.Guid)),
	//		fmt.Sprintf("cwd='%s'", input.DestinationPath),
	//	}
	//
	//	// support exclude log dir when backup
	//	if input.ExcludePath != "" {
	//		cmdArgs = append(cmdArgs, fmt.Sprintf("--exclude={%s}", input.ExcludePath))
	//	}
	//	cmdArgs = append(cmdArgs, ".")
	//
	//	log.Logger.Debug("App update", log.String("step", "backup"), log.JsonObj("cmdArgs", cmdArgs))
	//	if _, saltErr := CallSaltApi("https://127.0.0.1:8080", SaltApiRequest{
	//		Client:   "local",
	//		Function: "archive.tar",
	//		Target:   input.Target,
	//		Args:     cmdArgs,
	//	}, action.Language); saltErr != nil {
	//		errMsg := fmt.Sprintf("Failed when archive %s to %s, %s\n", input.DestinationPath, input.AppBackUpPath, saltErr.Error())
	//		output.FileDetail = errMsg
	//		return output, fmt.Errorf(errMsg)
	//	}
	//}

	// copy apply package
	fileCopyRequest := FileCopyInputs{
		Inputs: []FileCopyInput{
			FileCopyInput{
				EndPoint:        output.NewS3PkgPath,
				Guid:            input.Guid,
				Target:          input.Target,
				DestinationPath: input.DestinationPath,
				Unpack:          "true",
				FileOwner:       input.UserName,
			},
		},
	}

	log.Logger.Debug("App update", log.String("step", "copy file"), log.JsonObj("param", fileCopyRequest))
	fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
	if err != nil {
		return output, fmt.Errorf("Copy file to target fail,%s ", err.Error())
	}
	output.FileDetail += fileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail

	// copy apply package's sign file from s3 to dst
	if input.SignFileSrc != "" && input.SignFileDst != "" {
		signFileCopyRequest := FileCopyInputs{
			Inputs: []FileCopyInput{
				{
					EndPoint:        input.SignFileSrc,
					Guid:            input.Guid,
					Target:          input.Target,
					DestinationPath: input.SignFileDst,
					Unpack:          "false",
					FileOwner:       input.UserName,
				},
			},
		}

		log.Logger.Debug("App deploy", log.String("step", "sign file copy"), log.JsonObj("param", signFileCopyRequest))
		signFileCopyOutputs, err := copyApplyFile(signFileCopyRequest)
		if err != nil {
			return output, fmt.Errorf("Copy app package's sign file to target fail,%s ", err.Error())
		}
		output.FileDetail = signFileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail
	}

	// start apply script
	runStartScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StartScriptPath,
				Target:       input.Target,
				RunAs:        strings.Split(input.UserName, ":")[0],
				Guid:         input.Guid,
			},
		},
	}
	if input.ExecArg != "" {
		runStartScriptRequest.Inputs[0].ExecArg = input.ExecArg
	}

	log.Logger.Debug("App update", log.String("step", "run start script"), log.JsonObj("param", runStartScriptRequest))
	runStartScriptOutputs, err := runApplyScript(runStartScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run start script fail,%s ", err.Error())
	}
	output.RunStartScriptDetail = runStartScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	// parse log files from glob name
	action.parseGlobLogFiles(&input, &output)
	return output, err
}

func (action *ApplyUpdateDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyUpdateDeploymentInputs)
	outputs := ApplyUpdateDeploymentOutputs{}
	var finalErr error
	outputChan := make(chan ApplyUpdateDeploymentThreadObj, len(inputs.Inputs))
	wg := sync.WaitGroup{}
	for i, input := range inputs.Inputs {
		wg.Add(1)
		go func(tmpInput ApplyUpdateDeploymentInput, index int) {
			output, err := action.applyUpdateDeployment(tmpInput)
			outputChan <- ApplyUpdateDeploymentThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
		}(input, i)
		outputs.Outputs = append(outputs.Outputs, ApplyUpdateDeploymentOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("App update deploy action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}
	return &outputs, finalErr
}

// ApplyRollbackDeploymentAction inherit ApplyUpdateDeploymentAction for backup file reuse
type ApplyRollbackDeploymentAction struct {
	ApplyUpdateDeploymentAction
}

func (action *ApplyRollbackDeploymentAction) applyRollbackDeployment(input ApplyUpdateDeploymentInput) (output ApplyUpdateDeploymentOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Target = input.Target
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.RetCode = 1
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(input)
	if err != nil {
		return output, err
	}
	input.Seed = getEncryptSeed(input.Seed)
	if !strings.Contains(input.UserName, ":") {
		input.UserName = fmt.Sprintf("%s:%s", input.UserName, input.UserName)
	}
	if !strings.HasPrefix(input.StartScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newScriptPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StartScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript[1:])
			} else {
				newScriptPathList = append(newScriptPathList, input.DestinationPath+oldScript)
			}
		}
		input.StartScriptPath = strings.Join(newScriptPathList, ",")
	}
	if !strings.HasPrefix(input.StopScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newStopPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StopScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript[1:])
			} else {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript)
			}
		}
		input.StopScriptPath = strings.Join(newStopPathList, ",")
	}
	// stop apply script
	runStopScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StopScriptPath,
				Target:       input.Target,
				RunAs:        strings.Split(input.UserName, ":")[0],
				Guid:         input.Guid,
			},
		},
	}

	log.Logger.Debug("App update", log.String("step", "run stop script"), log.JsonObj("param", runStopScriptRequest))
	runStopScriptOutputs, err := runApplyScript(runStopScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run stop script fail,%s ", err.Error())
	}
	output.RunStopScriptDetail = runStopScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	log.Logger.Debug("App update", log.String("step", "extract backup file"), log.JsonObj("input", input))

	// extract backup file
	//if input.AppBackUpEnabled != "Y" || input.AppBackUpPath == "" {
	//	errMsg := fmt.Sprintf("rollback disabled, AppBackUpEnabled=%v, AppBackUpPath=%s", input.AppBackUpEnabled, input.AppBackUpPath)
	//	output.FileDetail = errMsg
	//	return output, fmt.Errorf(errMsg)
	//}

	// salt '*' archive.tar zxf {{.AppBackUpPath}}/{{.guid}}.tar.gz dest='{{.SourcePath}}'
	//tarFile := path.Join(input.AppBackUpPath, fmt.Sprintf("%s.tar.gz", input.Guid))
	//if _, saltErr := CallSaltApi("https://127.0.0.1:8080", SaltApiRequest{
	//	Client:   "local",
	//	Function: "archive.tar",
	//	Target:   input.Target,
	//	Args:     []string{"zxf", tarFile, fmt.Sprintf("dest='%s'", input.DestinationPath)},
	//}, action.Language); saltErr != nil {
	//	errMsg := fmt.Sprintf("Failed when extract %s to %s, %s\n", input.DestinationPath, input.AppBackUpPath, saltErr.Error())
	//	output.FileDetail = errMsg
	//	return output, fmt.Errorf(errMsg)
	//}
	//output.FileDetail += fmt.Sprintf("extract backup file success: %s -> %s", tarFile, input.DestinationPath)

	// start apply script
	runStartScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StartScriptPath,
				Target:       input.Target,
				RunAs:        strings.Split(input.UserName, ":")[0],
				Guid:         input.Guid,
			},
		},
	}
	if input.ExecArg != "" {
		runStartScriptRequest.Inputs[0].ExecArg = input.ExecArg
	}

	log.Logger.Debug("App update", log.String("step", "run start script"), log.JsonObj("param", runStartScriptRequest))
	runStartScriptOutputs, err := runApplyScript(runStartScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run start script fail,%s ", err.Error())
	}
	output.RunStartScriptDetail = runStartScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	return output, err
}

func (action *ApplyRollbackDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyUpdateDeploymentInputs)
	outputs := ApplyUpdateDeploymentOutputs{}
	var finalErr error
	outputChan := make(chan ApplyUpdateDeploymentThreadObj, len(inputs.Inputs))
	wg := sync.WaitGroup{}
	for i, input := range inputs.Inputs {
		wg.Add(1)
		go func(tmpInput ApplyUpdateDeploymentInput, index int) {
			output, err := action.applyRollbackDeployment(tmpInput)
			outputChan <- ApplyUpdateDeploymentThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
		}(input, i)
		outputs.Outputs = append(outputs.Outputs, ApplyUpdateDeploymentOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("App rollback deploy action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}
	return &outputs, finalErr
}

func (action *ApplyUpdateDeploymentAction) parseGlobLogFiles(input *ApplyUpdateDeploymentInput, output *ApplyUpdateDeploymentOutput) {
	for _, logFile := range []struct {
		logFilePattern string
		outputField    *string
	}{
		{input.LogFileTrade, &output.LogFileTrade},
		{input.LogFileTrace, &output.LogFileTrace},
		{input.LogFileMetric, &output.LogFileMetric},
		{input.LogFileKeyword, &output.LogFileKeyword},
	} {
		if logFile.logFilePattern == "" {
			continue
		}

		log.Logger.Debug("App update", log.String("step", "parse log file: "+logFile.logFilePattern))
		if ret, err := FindGlobFiles(input.DestinationPath, logFile.logFilePattern, input.Target); err == nil {
			*logFile.outputField = ret
		}
	}
}

func createApplyUser(input interface{}) (interface{}, error) {
	addUserAction := new(AddUserAction)

	userAddOutputs, err := addUserAction.Do(input)
	if err != nil {
		log.Logger.Error("Create app user fail", log.Error(err))
		return nil, err
	}

	return userAddOutputs, nil
}

func replaceApplyVariable(input interface{}) (interface{}, error) {
	variableReplaceAction := new(VariableReplaceAction)

	variableReplaceOutputs, err := variableReplaceAction.Do(input)
	if err != nil {
		log.Logger.Error("Replace app variable fail", log.Error(err))
		return nil, err
	}

	return variableReplaceOutputs, nil
}

func copyApplyFile(input interface{}) (interface{}, error) {
	fileCopyAction := new(FileCopyAction)

	fileCopyOutputs, err := fileCopyAction.Do(input)
	if err != nil {
		log.Logger.Error("Copy app file fail", log.Error(err))
		return nil, err
	}

	return fileCopyOutputs, nil
}

func runApplyScript(input interface{}) (interface{}, error) {
	runScriptAction := new(RunScriptAction)

	runScriptOutputs, err := runScriptAction.Do(input)
	if err != nil {
		log.Logger.Error("Run app script fail", log.Error(err))
		return nil, err
	}

	return runScriptOutputs, nil
}

type ApplyDeleteDeploymentAction struct {
	Language string
}

type ApplyDeleteDeploymentInputs struct {
	Inputs []ApplyDeleteDeploymentInput `json:"inputs,omitempty"`
}

type ApplyDeleteDeploymentInput struct {
	CallBackParameter
	Guid            string `json:"guid,omitempty"`
	UserName        string `json:"userName,omitempty"`
	Target          string `json:"target,omitempty"`
	StopScriptPath  string `json:"stopScript,omitempty"`
	DestinationPath string `json:"destinationPath,omitempty"`
}

type ApplyDeleteDeploymentOutputs struct {
	Outputs []ApplyDeleteDeploymentOutput `json:"outputs,omitempty"`
}

type ApplyDeleteDeploymentOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *ApplyDeleteDeploymentAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *ApplyDeleteDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyDeleteDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ApplyDeleteDeploymentAction) deleteDeploymentCheckParam(input ApplyDeleteDeploymentInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.StopScriptPath == "" {
		return getParamEmptyError(action.Language, "stopScriptPath")
	}
	if input.DestinationPath == "" {
		return getParamEmptyError(action.Language, "destinationPath")
	}

	return nil
}

func (action *ApplyDeleteDeploymentAction) applyDeleteDeployment(input *ApplyDeleteDeploymentInput) (output ApplyDeleteDeploymentOutput, err error) {
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

	if err = action.deleteDeploymentCheckParam(*input); err != nil {
		return output, err
	}
	if !strings.HasPrefix(input.StopScriptPath, input.DestinationPath) {
		if !strings.HasSuffix(input.DestinationPath, "/") {
			input.DestinationPath = input.DestinationPath + "/"
		}
		var newStopPathList []string
		for _, oldScript := range splitWithCustomFlag(input.StopScriptPath) {
			if strings.HasPrefix(oldScript, "/") {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript[1:])
			} else {
				newStopPathList = append(newStopPathList, input.DestinationPath+oldScript)
			}
		}
		input.StopScriptPath = strings.Join(newStopPathList, ",")
	}
	// stop apply script
	runStopScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StopScriptPath,
				Target:       input.Target,
				RunAs:        input.UserName,
				Guid:         input.Guid,
			},
		},
	}

	log.Logger.Debug("App delete", log.JsonObj("param", runStopScriptRequest))
	_, err = runApplyScript(runStopScriptRequest)
	if err != nil {
		return output, fmt.Errorf("Run stop script fail,%s ", err.Error())
	}

	// rm package-dir
	err = deleteDirecory(input.Target, input.DestinationPath, action.Language)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("ApplyDeleteAction remove target[%v] dir[%v] meet error=%v", input.Target, input.DestinationPath, err))
	}

	return output, err
}

func deleteDirecory(target, destinationPath, language string) error {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = target
	request.Function = "cmd.run"

	cmdRun := "rm -rf " + destinationPath
	request.Args = append(request.Args, cmdRun)

	_, err := CallSaltApi("https://127.0.0.1:8080", request, language)
	if err != nil {
		return err
	}

	return nil
}

func (action *ApplyDeleteDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyDeleteDeploymentInputs)
	outputs := ApplyDeleteDeploymentOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.applyDeleteDeployment(&input)
		if err != nil {
			log.Logger.Error("App delete deploy action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
