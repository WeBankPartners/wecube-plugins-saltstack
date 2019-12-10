package plugins

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

var DiskPluginActions = make(map[string]Action)

func init() {
	DiskPluginActions["getUnformatedDisk"] = new(GetUnformatedDiskAction)
	DiskPluginActions["formatAndMountDisk"] = new(FormatAndMountDiskAction)
}

type DiskPlugin struct {
}

func (plugin *DiskPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := DiskPluginActions[actionName]

	if !found {
		return nil, fmt.Errorf("Script plugin,action = %s not found", actionName)
	}

	return action, nil
}

type GetUnformatedDiskAction struct {
}

type GetUnformatedDiskInputs struct {
	Inputs []GetUnformatedDiskInput `json:"inputs,omitempty"`
}

type GetUnformatedDiskInput struct {
	CallBackParameter
	Guid   string `json:"guid,omitempty"`
	Target string `json:"target,omitempty"`
}

type GetUnformatedDiskOutputs struct {
	Outputs []GetUnformatedDiskOutput `json:"outputs,omitempty"`
}

type GetUnformatedDiskOutput struct {
	CallBackParameter
	Result
	Guid            string   `json:"guid,omitempty"`
	UnformatedDisks []string `json:"unformatedDisks,omitempty"`
}

func (action *GetUnformatedDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs GetUnformatedDiskInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetUnformatedDiskAction) CheckParam(input GetUnformatedDiskInput) error {
	if input.Target == "" {
		return errors.New("Target is empty")
	}

	return nil
}

func (action *GetUnformatedDiskAction) getUnformatedDisk(input *GetUnformatedDiskInput) (output GetUnformatedDiskOutput, err error) {
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

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	result, err := executeS3Script("getUnformatedDisk.py", input.Target, "", "")
	if err != nil {
		return output, err
	}

	saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
	if err != nil {
		logrus.Errorf("parseSaltApiCmdScriptCallResult meet err=%v,rawStr=%s", err, result)
		return output, err
	}

	for k, v := range saltApiResult.Results[0] {
		if v.RetCode != 0 {
			logrus.Errorf("GetUnformatedDiskAction ip=%v,stderr=%v", k, v.Stderr)
			err = fmt.Errorf("GetUnformatedDiskAction ip=%v,stderr=%v", k, v.Stderr)
			return output, err
		}
		if err = json.Unmarshal([]byte(v.Stdout), &output); err != nil {
			logrus.Errorf("GetUnformatedDiskAction Unmarshal failed err=%v,stdOut=%v", err, v.Stdout)
			err = fmt.Errorf("GetUnformatedDiskAction Unmarshal failed err=%v,stdOut=%v", err, v.Stdout)
			return output, err
		}
	}

	return output, err
}

func (action *GetUnformatedDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(GetUnformatedDiskInputs)
	outputs := GetUnformatedDiskOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.getUnformatedDisk(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

type FormatAndMountDiskAction struct {
}
type FormatAndMountDiskInputs struct {
	Inputs []FormatAndMountDiskInput `json:"inputs,omitempty"`
}

type FormatAndMountDiskInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	Target         string `json:"target,omitempty"`
	DiskName       string `json:"diskName,omitempty"`
	FileSystemType string `json:"fileSystemType,omitempty"`
	MountDir       string `json:"mountDir,omitempty"`
}

type FormatAndMountDiskOutputs struct {
	Outputs []FormatAndMountDiskOutput `json:"outputs,omitempty"`
}

type FormatAndMountDiskOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (action *FormatAndMountDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FormatAndMountDiskInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func isValidFileSystemType(fileSystemType string) error {
	validFileSystemTypes := []string{"ext3", "ext4", "xfs"}
	for _, valid := range validFileSystemTypes {
		if valid == fileSystemType {
			return nil
		}
	}
	return fmt.Errorf("invalid fileSystemType(%s)", fileSystemType)
}

func (action *FormatAndMountDiskAction) CheckParam(input FormatAndMountDiskInput) error {
	if input.Target == "" {
		return errors.New("Target is empty")
	}
	if input.DiskName == "" {
		return errors.New("DiskName is empty")
	}
	if input.FileSystemType == "" {
		return errors.New("FileSystemType is empty")
	}
	if input.MountDir == "" {
		return errors.New("MountDir is empty")
	}
	if err := isValidFileSystemType(input.FileSystemType); err != nil {
		return err
	}

	return nil
}

func (action *FormatAndMountDiskAction) formatAndMountDisk(input *FormatAndMountDiskInput) (output FormatAndMountDiskOutput, err error) {
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

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	execArgs := "-d " + input.DiskName + " -f " + input.FileSystemType + " -m " + input.MountDir
	result, err := executeS3Script("formatAndMountDisk.py", input.Target, "", execArgs)
	if err != nil {
		return output, err
	}

	saltApiResult, err := parseSaltApiCmdScriptCallResult(result)
	if err != nil {
		logrus.Errorf("parseSaltApiCmdScriptCallResult meet err=%v,rawStr=%s", err, result)
		return output, err
	}

	for k, v := range saltApiResult.Results[0] {
		if v.RetCode != 0 {
			logrus.Errorf("FormatAndMountDiskAction ip=%v,stderr=%v", k, v.Stderr)
			err = fmt.Errorf("FormatAndMountDiskAction ip=%v,stderr=%v", k, v.Stderr)
			return output, err
		}
		output.Detail = v.Stdout
	}

	return output, err
}

func (action *FormatAndMountDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(FormatAndMountDiskInputs)
	outputs := FormatAndMountDiskOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.formatAndMountDisk(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}
