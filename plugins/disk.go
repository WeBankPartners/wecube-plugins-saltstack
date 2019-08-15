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
	Target string `json:"target,omitempty"`
}

type GetUnformatedDiskOutputs struct {
	Outputs []GetUnformatedDiskOutput `json:"outputs,omitempty"`
}

type GetUnformatedDiskOutput struct {
	UnformatedDisks []string `json:"unformatedDisks,omitempty"`
}

func (action *GetUnformatedDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs GetUnformatedDiskInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetUnformatedDiskAction) CheckParam(input interface{}) error {
	inputs, ok := input.(GetUnformatedDiskInputs)
	if !ok {
		return fmt.Errorf("GetUnformatedDiskAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.Target == "" {
			return errors.New("Target is empty")
		}
	}

	return nil
}

func (action *GetUnformatedDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(GetUnformatedDiskInputs)
	outputs := GetUnformatedDiskOutputs{}

	for _, input := range inputs.Inputs {
		result, err := executeScript("getUnformatedDisk.py", input.Target, "", "")
		if err != nil {
			return nil, err
		}

		saltApiResult, err := parseSaltApiCallResult(result)
		if err != nil {
			logrus.Errorf("parseSaltApiCallResult meet err=%v,rawStr=%s", err, result)
			return nil, err
		}

		output := GetUnformatedDiskOutput{}
		for k, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				logrus.Errorf("GetUnformatedDiskAction ip=%v,stderr=%v", k, v.Stderr)
				return nil, fmt.Errorf("GetUnformatedDiskAction ip=%v,stderr=%v", k, v.Stderr)
			}
			if err = json.Unmarshal([]byte(v.Stdout), &output); err != nil {
				logrus.Errorf("GetUnformatedDiskAction Unmarshal failed err=%v,stdOut=%v", err, v.Stdout)
				return nil, fmt.Errorf("GetUnformatedDiskAction Unmarshal failed err=%v,stdOut=%v", err, v.Stdout)
			}
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, nil
}

type FormatAndMountDiskAction struct {
}
type FormatAndMountDiskInputs struct {
	Inputs []FormatAndMountDiskInput `json:"inputs,omitempty"`
}

type FormatAndMountDiskInput struct {
	Target         string `json:"target,omitempty"`
	DiskName       string `json:"diskName,omitempty"`
	FileSystemType string `json:"fileSystemType,omitempty"`
	MountDir       string `json:"mountDir,omitempty"`
}

type FormatAndMountDiskOutputs struct {
	Outputs []FormatAndMountDiskOutput `json:"outputs,omitempty"`
}

type FormatAndMountDiskOutput struct {
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

func (action *FormatAndMountDiskAction) CheckParam(input interface{}) error {
	inputs, ok := input.(FormatAndMountDiskInputs)
	if !ok {
		return fmt.Errorf("FormatAndMountDiskInputs:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
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
	}

	return nil
}

func (action *FormatAndMountDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(FormatAndMountDiskInputs)
	outputs := FormatAndMountDiskOutputs{}

	for _, input := range inputs.Inputs {
		execArgs := "-d " + input.DiskName +" -f "+input.FileSystemType + " -m " + input.MountDir
		result, err := executeScript("formatAndMountDisk.py", input.Target, "", execArgs)
		if err != nil {
			return nil, err
		}

		saltApiResult, err := parseSaltApiCallResult(result)
		if err != nil {
			logrus.Errorf("parseSaltApiCallResult meet err=%v,rawStr=%s", err, result)
			return nil, err
		}

		output := FormatAndMountDiskOutput{}
		for k, v := range saltApiResult.Results[0] {
			if v.RetCode != 0 {
				logrus.Errorf("FormatAndMountDiskAction ip=%v,stderr=%v", k, v.Stderr)
				return nil, fmt.Errorf("FormatAndMountDiskAction ip=%v,stderr=%v", k, v.Stderr)
			}
			output.Detail = v.Stdout
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, nil
}
