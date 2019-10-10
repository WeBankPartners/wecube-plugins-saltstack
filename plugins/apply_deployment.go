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
	return nil, nil
}

type ApplyNewDeploymentInputs struct {
	Inputs []ApplyNewDeploymentInput `json:"inputs,omitempty"`
}
type ApplyNewDeploymentInput struct {
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey    string `json:"accessKey,omitempty"`
	// SecretKey    string `json:"secretKey,omitempty"`
	Guid             string `json:"guid,omitempty"`
	UserName         string `json:"user_name,omitempty"`
	Target           string `json:"target,omitempty"`
	DestinationPath  string `json:"destination_path,omitempty"`
	VariableFilePath string `json:"conf_files,omitempty"`
	VariableList     string `json:"variable_list,omitempty"`
	ExecArg          string `json:"args,omitempty"`
	StartScriptPath  string `json:"start_script,omitempty"`
}

type ApplyNewDeploymentOutputs struct {
	Outputs []ApplyNewDeploymentOutput `json:"outputs,omitempty"`
}

type ApplyNewDeploymentOutput struct {
	Guid            string `json:"guid,omitempty"`
	UserDetail      string `json:"user_detail,omitempty"`
	FileDetail      string `json:"file_detail"`
	NewS3PkgPath    string `json:"s3_pkg_path,omitempty"`
	Target          string `json:"target"`
	RetCode         int    `json:"ret_code"`
	RunScriptDetail string `json:"run_script_detail"`
}

type ApplyNewDeploymentAction struct {
}

func (action *ApplyNewDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyNewDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return &inputs, nil
}

func (action *ApplyNewDeploymentAction) CheckParam(input interface{}) error {
	inputs, ok := input.(ApplyNewDeploymentInputs)
	if !ok {
		return fmt.Errorf("ApplyNewDeploymentAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.EndPoint == "" {
			return errors.New("EndPoint is empty")
		}
		if input.UserName == "" {
			return errors.New("UserName is empty")
		}
		if input.Target == "" {
			return errors.New("Target is empty")
		}
		if input.VariableFilePath == "" {
			return errors.New("VariableFilePath is empty")
		}
		if input.StartScriptPath == "" {
			return errors.New("StartScriptPath is empty")
		}
		if input.DestinationPath == "" {
			return errors.New("DestinationPath is empty")
		}
		if input.VariableList == "" {
			return errors.New("VariableList is empty")
		}
	}

	return nil
}

func (action *ApplyNewDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyNewDeploymentInputs)
	outputs := ApplyNewDeploymentOutputs{}

	for _, input := range inputs.Inputs {
		output := ApplyNewDeploymentOutput{
			Guid:   input.Guid,
			Target: input.Target,
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

		userAddOutputs, err := createApplyUser(addUserRequest)
		if err != nil {
			logrus.Errorf("ApplyNewDeploymentAction createApplyUser meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.UserDetail = userAddOutputs.(AddUserOutputs).Outputs[0].Detail

		// replace apply variable
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:         input.Guid,
					EndPoint:     input.EndPoint,
					FilePath:     input.VariableFilePath,
					VariableList: input.VariableList,
				},
			},
		}
		variableReplaceOutputs, err := replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			logrus.Errorf("ApplyNewDeploymentAction replaceApplyVariable meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.NewS3PkgPath = variableReplaceOutputs.(VariableReplaceOutputs).Outputs[0].NewS3PkgPath

		// copy apply package
		fileCopyRequest := FileCopyInputs{
			Inputs: []FileCopyInput{
				FileCopyInput{
					EndPoint:        output.NewS3PkgPath,
					Guid:            input.Guid,
					Target:          input.Target,
					DestinationPath: input.DestinationPath,
					Unpack:          "true",
				},
			},
		}
		fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
		if err != nil {
			logrus.Errorf("ApplyNewDeploymentAction copyApplyFile meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.FileDetail = fileCopyOutputs.(FileCopyOutputs).Outputs[0].Detail

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
		runScriptOutputs, err := runApplyScript(runScriptRequest)
		if err != nil {
			logrus.Errorf("ApplyNewDeploymentAction runApplyScript meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.RunScriptDetail = runScriptOutputs.(RunScriptOutputs).Outputs[0].Detail
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

type ApplyUpdateDeploymentInputs struct {
	Inputs []ApplyUpdateDeploymentInput `json:"inputs,omitempty"`
}

type ApplyUpdateDeploymentInput struct {
	EndPoint string `json:"end_point,omitempty"`
	// AccessKey    string `json:"accessKey,omitempty"`
	// SecretKey    string `json:"secretKey,omitempty"`
	Guid             string `json:"guid,omitempty"`
	UserName         string `json:"user_name,omitempty"`
	Target           string `json:"target,omitempty"`
	DestinationPath  string `json:"destination_path,omitempty"`
	VariableFilePath string `json:"conf_files,omitempty"`
	VariableList     string `json:"variable_list,omitempty"`
	ExecArg          string `json:"args,omitempty"`
	StopScriptPath   string `json:"stop_script,omitempty"`
	StartScriptPath  string `json:"start_script,omitempty"`
}

type ApplyUpdateDeploymentOutputs struct {
	Outputs []ApplyUpdateDeploymentOutput `json:"outputs,omitempty"`
}

type ApplyUpdateDeploymentOutput struct {
	Guid                 string `json:"guid,omitempty"`
	FileDetail           string `json:"file_detail"`
	NewS3PkgPath         string `json:"s3_pkg_path,omitempty"`
	Target               string `json:"target"`
	RetCode              int    `json:"ret_code"`
	RunStartScriptDetail string `json:"run_start_script_detail"`
	RunStopScriptDetail  string `json:"run_stop_script_detail"`
}

type ApplyUpdateDeploymentAction struct {
}

func (action *ApplyUpdateDeploymentAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ApplyUpdateDeploymentInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return &inputs, nil
}

func (action *ApplyUpdateDeploymentAction) CheckParam(input interface{}) error {
	inputs, ok := input.(ApplyUpdateDeploymentInputs)
	if !ok {
		return fmt.Errorf("ApplyUpdateInputs:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.EndPoint == "" {
			return errors.New("EndPoint is empty")
		}
		if input.UserName == "" {
			return errors.New("UserName is empty")
		}
		if input.Target == "" {
			return errors.New("Target is empty")
		}
		if input.VariableFilePath == "" {
			return errors.New("VariableFilePath is empty")
		}
		if input.StartScriptPath == "" {
			return errors.New("StartScriptPath is empty")
		}
		if input.DestinationPath == "" {
			return errors.New("DestinationPath is empty")
		}
		if input.VariableList == "" {
			return errors.New("VariableList is empty")
		}
		if input.StopScriptPath == "" {
			return errors.New("StopScriptPath is empty")
		}
	}

	return nil
}

func (action *ApplyUpdateDeploymentAction) Do(input interface{}) (interface{}, error) {
	inputs := input.(ApplyUpdateDeploymentInputs)
	outputs := ApplyUpdateDeploymentOutputs{}

	for _, input := range inputs.Inputs {
		output := ApplyUpdateDeploymentOutput{
			Guid:   input.Guid,
			Target: input.Target,
		}

		// stop apply script
		runStopScriptRequest := RunScriptInputs{
			Inputs: []RunScriptInput{
				RunScriptInput{
					EndPointType: "LOCAL",
					EndPoint:     input.StopScriptPath,
					Target:       input.Target,
					RunAs:        "",
					Guid:         input.Guid,
				},
			},
		}
		runStopScriptOutputs, err := runApplyScript(runStopScriptRequest)
		if err != nil {
			logrus.Errorf("ApplyUpdateAction runApplyScript meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.RunStopScriptDetail = runStopScriptOutputs.(RunScriptOutputs).Outputs[0].Detail

		// replace apply variable
		variableReplaceRequest := VariableReplaceInputs{
			Inputs: []VariableReplaceInput{
				VariableReplaceInput{
					Guid:         input.Guid,
					EndPoint:     input.EndPoint,
					FilePath:     input.VariableFilePath,
					VariableList: input.VariableList,
				},
			},
		}
		variableReplaceOutputs, err := replaceApplyVariable(variableReplaceRequest)
		if err != nil {
			logrus.Errorf("ApplyUpdateAction replaceApplyVariable meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.NewS3PkgPath = variableReplaceOutputs.(VariableReplaceOutputs).Outputs[0].NewS3PkgPath

		// copy apply package
		fileCopyRequest := FileCopyInputs{
			Inputs: []FileCopyInput{
				FileCopyInput{
					EndPoint:        output.NewS3PkgPath,
					Guid:            input.Guid,
					Target:          input.Target,
					DestinationPath: input.DestinationPath,
					Unpack:          "true",
				},
			},
		}
		fileCopyOutputs, err := copyApplyFile(fileCopyRequest)
		if err != nil {
			logrus.Errorf("ApplyUpdateAction copyApplyFile meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.FileDetail = fileCopyOutputs.(FileCopyOutputs).Outputs[0].Detail

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
		runStartScriptOutputs, err := runApplyScript(runStartScriptRequest)
		if err != nil {
			logrus.Errorf("ApplyUpdateAction runApplyScript meet error=%v", err)
			output.RetCode = 1
			outputs.Outputs = append(outputs.Outputs, output)
			return &outputs, err
		}
		output.RunStartScriptDetail = runStartScriptOutputs.(RunScriptOutputs).Outputs[0].Detail
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, nil
}

func createApplyUser(input AddUserInputs) (interface{}, error) {
	addUserAction := new(AddUserAction)

	userAddInpurt, err := addUserAction.ReadParam(input)
	if err != nil {
		logrus.Errorf("createApplyUser ReadParam meet error=%v", err)
		return nil, err
	}
	if err = addUserAction.CheckParam(userAddInpurt); err != nil {
		logrus.Errorf("createApplyUser CheckParam meet error=%v", err)
		return nil, err
	}
	userAddOutputs, err := addUserAction.Do(userAddInpurt)
	if err != nil {
		logrus.Errorf("createApplyUser Do meet error=%v", err)
		return nil, err
	}

	return userAddOutputs, nil
}

func replaceApplyVariable(input VariableReplaceInputs) (interface{}, error) {
	variableReplaceAction := new(VariableReplaceAction)

	variableReplaceInput, err := variableReplaceAction.ReadParam(input)
	if err != nil {
		logrus.Errorf("replaceApplyVariable ReadParam meet error=%v", err)
		return nil, err
	}
	if err = variableReplaceAction.CheckParam(variableReplaceInput); err != nil {
		logrus.Errorf("replaceApplyVariable CheckParam meet error=%v", err)
		return nil, err
	}
	variableReplaceOutputs, err := variableReplaceAction.Do(variableReplaceInput)
	if err != nil {
		logrus.Errorf("replaceApplyVariable Do meet error=%v", err)
		return nil, err
	}

	return variableReplaceOutputs, nil
}

func copyApplyFile(input FileCopyInputs) (interface{}, error) {
	fileCopyAction := new(FileCopyAction)

	fileCopyInput, err := fileCopyAction.ReadParam(input)
	if err != nil {
		logrus.Errorf("copyApplyFile ReadParam meet error=%v", err)
		return nil, err
	}
	if err = fileCopyAction.CheckParam(fileCopyInput); err != nil {
		logrus.Errorf("copyApplyFile CheckParam meet error=%v", err)
		return nil, err
	}
	fileCopyOutputs, err := fileCopyAction.Do(fileCopyInput)
	if err != nil {
		logrus.Errorf("copyApplyFile Do meet error=%v", err)
		return nil, err
	}

	return fileCopyOutputs, nil
}

func runApplyScript(input RunScriptInputs) (interface{}, error) {
	runScriptAction := new(RunScriptAction)

	runScriptInput, err := runScriptAction.ReadParam(input)
	if err != nil {
		logrus.Errorf("runApplyScript ReadParam meet error=%v", err)
		return nil, err
	}
	if err = runScriptAction.CheckParam(runScriptInput); err != nil {
		logrus.Errorf("runApplyScript CheckParam meet error=%v", err)
		return nil, err
	}
	runScriptOutputs, err := runScriptAction.Do(runScriptInput)
	if err != nil {
		logrus.Errorf("runApplyScript Do meet error=%v", err)
		return nil, err
	}

	return runScriptOutputs, nil
}
