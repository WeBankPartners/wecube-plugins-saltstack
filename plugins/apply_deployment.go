package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

var ApplyDeploymentActions = make(map[string]Action)

func init() {
	ApplyDeploymentActions["new"] = new(ApplyNewDeploymentAction)
	ApplyDeploymentActions["update"] = new(ApplyUpdateDeploymentAction)
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
	EncryptVariblePrefix string `json:"encryptVariblePrefix,omitempty"`
	Seed                 string `json:"seed,omitempty"`
	AppPublicKey         string `json:"appPublicKey,omitempty"`
	SysPrivateKey        string `json:"sysPrivateKey,omitempty"`
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
}

type ApplyNewDeploymentAction struct {
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
		return errors.New("EndPoint is empty")
	}
	if input.UserName == "" {
		return errors.New("UserName is empty")
	}
	if input.Target == "" {
		return errors.New("Target is empty")
	}
	// if input.VariableFilePath == "" {
	// 	return errors.New("VariableFilePath is empty")
	// }
	if input.StartScriptPath == "" {
		return errors.New("StartScriptPath is empty")
	}
	if input.DestinationPath == "" {
		return errors.New("DestinationPath is empty")
	}
	// if input.VariableList == "" {
	// 	return errors.New("VariableList is empty")
	// }

	return nil
}

func (action *ApplyNewDeploymentAction) applyNewDeployment(input *ApplyNewDeploymentInput) (output ApplyNewDeploymentOutput, err error) {
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

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	//create apply deployment user
	addUserRequest := AddUserInputs{
		Inputs: []AddUserInput{
			AddUserInput{
				Guid:     input.Guid,
				Target:   input.Target,
				UserName: input.UserName,
			},
		},
	}

	logrus.Infof("ApplyNewDeploymentAction createApplyUser: input=%++v", addUserRequest)
	userAddOutputs, err := createApplyUser(addUserRequest)
	if err != nil {
		logrus.Errorf("ApplyNewDeploymentAction createApplyUser meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyNewDeploymentAction: userAddOutputs=%++v", userAddOutputs.(*AddUserOutputs))
	output.UserDetail = userAddOutputs.(*AddUserOutputs).Outputs[0].Detail
	logrus.Infof("ApplyNewDeploymentAction: output=%++v", output)

	// replace apply variable
	var variableReplaceOutputs interface{}
	if input.VariableFilePath != "" {
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:                 input.Guid,
					EndPoint:             input.EndPoint,
					FilePath:             input.VariableFilePath,
					VariableList:         input.VariableList,
					EncryptVariblePrefix: input.EncryptVariblePrefix,
					Seed:                 input.Seed,
					AppPublicKey:         input.AppPublicKey,
					SysPrivateKey:        input.SysPrivateKey,
				},
			},
		}

		logrus.Infof("ApplyNewDeploymentAction replaceApplyVariable: input=%++v", variableReplaceRequest)
		variableReplaceOutputs, err = replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			logrus.Errorf("ApplyNewDeploymentAction replaceApplyVariable meet error=%v", err)
			return output, err
		}
		logrus.Infof("ApplyNewDeploymentAction: variableReplaceOutputs=%++v", variableReplaceOutputs.(*VariableReplaceOutputs))
		output.NewS3PkgPath = variableReplaceOutputs.(*VariableReplaceOutputs).Outputs[0].NewS3PkgPath
		logrus.Infof("ApplyNewDeploymentAction: output=%++v", output)
	} else {
		output.NewS3PkgPath = input.VariableFilePath
		logrus.Infof("ApplyNewDeploymentAction: output=%++v", output)
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

	logrus.Infof("ApplyNewDeploymentAction copyApplyFile: input=%++v", fileCopyRequest)
	fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
	if err != nil {
		logrus.Errorf("ApplyNewDeploymentAction copyApplyFile meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyNewDeploymentAction: fileCopyOutputs=%++v", fileCopyOutputs.(*FileCopyOutputs))
	output.FileDetail = fileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail
	logrus.Infof("ApplyNewDeploymentAction: output=%++v", output)

	// start apply script
	runScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StartScriptPath,
				Target:       input.Target,
				RunAs:        input.UserName,
				Guid:         input.Guid,
			},
		},
	}
	if input.ExecArg != "" {
		runScriptRequest.Inputs[0].ExecArg = input.ExecArg
	}

	logrus.Infof("ApplyNewDeploymentAction runApplyScript: input=%++v", runScriptRequest)
	runScriptOutputs, err := runApplyScript(runScriptRequest)
	if err != nil {
		logrus.Errorf("ApplyNewDeploymentAction runApplyScript meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyNewDeploymentAction: runScriptOutputs=%++v", runScriptOutputs.(*RunScriptOutputs))
	output.RunScriptDetail = runScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	return output, err
}

func (action *ApplyNewDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyNewDeploymentInputs)
	logrus.Infof("ApplyNewDeploymentAction Do: input=%++v", inputs)

	outputs := ApplyNewDeploymentOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.applyNewDeployment(&input)
		if err != nil {
			finalErr = err
		}
		logrus.Infof("ApplyNewDeploymentAction: output=%++v", output)

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("All new applications = %v have been done", inputs)
	return &outputs, finalErr
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

	EncryptVariblePrefix string `json:"encryptVariblePrefix,omitempty"`
	Seed                 string `json:"seed,omitempty"`
	AppPublicKey         string `json:"appPublicKey,omitempty"`
	SysPrivateKey        string `json:"sysPrivateKey,omitempty"`
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
}

type ApplyUpdateDeploymentAction struct {
}

func (action *ApplyUpdateDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyUpdateDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *ApplyUpdateDeploymentAction) CheckParam(input ApplyUpdateDeploymentInput) error {
	return fmt.Errorf("ApplyUpdateInputs:input type=%T not right", input)

	if input.EndPoint == "" {
		return errors.New("EndPoint is empty")
	}
	if input.UserName == "" {
		return errors.New("UserName is empty")
	}
	if input.Target == "" {
		return errors.New("Target is empty")
	}
	// if input.VariableFilePath == "" {
	// 	return errors.New("VariableFilePath is empty")
	// }
	if input.StartScriptPath == "" {
		return errors.New("StartScriptPath is empty")
	}
	if input.DestinationPath == "" {
		return errors.New("DestinationPath is empty")
	}
	// if input.VariableList == "" {
	// 	return errors.New("VariableList is empty")
	// }
	if input.StopScriptPath == "" {
		return errors.New("StopScriptPath is empty")
	}

	return nil
}

func (action *ApplyUpdateDeploymentAction) applyUpdateDeployment(input *ApplyUpdateDeploymentInput) (output ApplyUpdateDeploymentOutput, err error) {
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

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
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

	logrus.Infof("ApplyUpdateAction runApplyScript: input=%++v", runStopScriptRequest)
	runStopScriptOutputs, err := runApplyScript(runStopScriptRequest)
	if err != nil {
		logrus.Errorf("ApplyUpdateAction runApplyScript meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyUpdateAction: runStopScriptOutputs=%++v", runStopScriptOutputs.(*RunScriptOutputs))
	output.RunStopScriptDetail = runStopScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail
	logrus.Infof("ApplyUpdateAction: output=%++v", output)

	// replace apply variable
	var variableReplaceOutputs interface{}
	if input.VariableFilePath != "" {
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:                 input.Guid,
					EndPoint:             input.EndPoint,
					FilePath:             input.VariableFilePath,
					VariableList:         input.VariableList,
					EncryptVariblePrefix: input.EncryptVariblePrefix,
					Seed:                 input.Seed,
					AppPublicKey:         input.AppPublicKey,
					SysPrivateKey:        input.SysPrivateKey,
				},
			},
		}

		logrus.Infof("ApplyUpdateAction replaceApplyVariable: input=%++v", variableReplaceRequest)
		variableReplaceOutputs, err = replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			logrus.Errorf("ApplyUpdateAction replaceApplyVariable meet error=%v", err)
			return output, err
		}
		logrus.Infof("ApplyUpdateAction: variableReplaceOutputs=%++v", variableReplaceOutputs.(*VariableReplaceOutputs))
		output.NewS3PkgPath = variableReplaceOutputs.(*VariableReplaceOutputs).Outputs[0].NewS3PkgPath
		logrus.Infof("ApplyUpdateAction: output=%++v", output)
	} else {
		output.NewS3PkgPath = input.EndPoint
		logrus.Infof("ApplyUpdateAction: output=%++v", output)
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

	logrus.Infof("ApplyUpdateAction copyApplyFile: input=%++v", fileCopyRequest)
	fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
	if err != nil {
		logrus.Errorf("ApplyUpdateAction copyApplyFile meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyUpdateAction: fileCopyOutputs=%++v", fileCopyOutputs.(*FileCopyOutputs))
	output.FileDetail = fileCopyOutputs.(*FileCopyOutputs).Outputs[0].Detail
	logrus.Infof("ApplyUpdateAction: output=%++v", output)

	// start apply script
	runStartScriptRequest := RunScriptInputs{
		Inputs: []RunScriptInput{
			RunScriptInput{
				EndPointType: "LOCAL",
				EndPoint:     input.StartScriptPath,
				Target:       input.Target,
				RunAs:        input.UserName,
				Guid:         input.Guid,
			},
		},
	}
	if input.ExecArg != "" {
		runStartScriptRequest.Inputs[0].ExecArg = input.ExecArg
	}

	logrus.Infof("ApplyUpdateAction runApplyScript: input=%++v", runStartScriptRequest)
	runStartScriptOutputs, err := runApplyScript(runStartScriptRequest)
	if err != nil {
		logrus.Errorf("ApplyUpdateAction runApplyScript meet error=%v", err)
		return output, err
	}
	logrus.Infof("ApplyUpdateAction: runStartScriptOutputs=%++v", runStartScriptOutputs.(*RunScriptOutputs))
	output.RunStartScriptDetail = runStartScriptOutputs.(*RunScriptOutputs).Outputs[0].Detail

	return output, err
}

func (action *ApplyUpdateDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyUpdateDeploymentInputs)
	outputs := ApplyUpdateDeploymentOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.applyUpdateDeployment(&input)
		if err != nil {
			finalErr = err
		}
		logrus.Infof("ApplyUpdateAction: output=%++v", output)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("All applictions = %v have been updated", inputs)
	return &outputs, finalErr
}

func createApplyUser(input interface{}) (interface{}, error) {
	addUserAction := new(AddUserAction)

	userAddOutputs, err := addUserAction.Do(input)
	if err != nil {
		logrus.Errorf("createApplyUser Do meet error=%v", err)
		return nil, err
	}

	return userAddOutputs, nil
}

func replaceApplyVariable(input interface{}) (interface{}, error) {
	variableReplaceAction := new(VariableReplaceAction)

	variableReplaceOutputs, err := variableReplaceAction.Do(input)
	if err != nil {
		logrus.Errorf("replaceApplyVariable Do meet error=%v", err)
		return nil, err
	}

	return variableReplaceOutputs, nil
}

func copyApplyFile(input interface{}) (interface{}, error) {
	fileCopyAction := new(FileCopyAction)

	fileCopyOutputs, err := fileCopyAction.Do(input)
	if err != nil {
		logrus.Errorf("copyApplyFile Do meet error=%v", err)
		return nil, err
	}

	return fileCopyOutputs, nil
}

func runApplyScript(input interface{}) (interface{}, error) {
	runScriptAction := new(RunScriptAction)

	runScriptOutputs, err := runScriptAction.Do(input)
	if err != nil {
		logrus.Errorf("runApplyScript Do meet error=%v", err)
		return nil, err
	}

	return runScriptOutputs, nil
}
