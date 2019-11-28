package plugins

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var ReleasedPackagePluginActions = make(map[string]Action)

func init() {
	ReleasedPackagePluginActions["listCurrentDir"] = new(ListCurrentDirAction)
	ReleasedPackagePluginActions["getConfigFileKey"] = new(GetConfigFileKeyAction)
}

type ReleasedPackagePlugin struct {
}

func (ReleasedPackagePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ReleasedPackagePluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("User plugin,action = %s not found", actionName)
	}

	return action, nil
}

type ListFilesInputs struct {
	Inputs []ListFilesInput `json:"inputs,omitempty"`
}

type ListFilesInput struct {
	CallBackParameter
	EndPoint   string `json:"endpoint,omitempty"`
	CurrentDir string `json:"currentDir,omitempty"`
	// AccessKey  string `json:"accessKey,omitempty"`
	// SecretKey  string `json:"secretKey,omitempty"`
}

type ListFilesOutputs struct {
	Outputs []ListFilesOutput `json:"outputs,omitempty"`
}

type ListFilesOutput struct {
	CallBackParameter
	Files []FileNode `json:"files,omitempty"`
}

type ListCurrentDirAction struct {
}

func (action *ListCurrentDirAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ListFilesInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *ListCurrentDirAction) CheckParam(input interface{}) error {
	inputs, ok := input.(ListFilesInputs)
	if !ok {
		return fmt.Errorf("ListCurrentDirAction:input type=%T not right", input)
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
	}

	return nil
}

func getPackageNameFromEndpoint(endpoint string) (string, error) {
	index := strings.LastIndexAny(endpoint, "/")
	if index == -1 {
		return "", fmt.Errorf("Invalid endpoint %s", endpoint)
	}

	return endpoint[index+1:], nil
}

func (action *ListCurrentDirAction) Do(input interface{}) (interface{}, error) {
	outputs := ListFilesOutputs{}
	inputs, ok := input.(ListFilesInputs)
	if !ok {
		return &outputs, fmt.Errorf("ListCurrentDirAction:input type=%T not right", input)
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic err=%v", err)
		}
	}()

	for _, input := range inputs.Inputs {
		packageName, err := getPackageNameFromEndpoint(input.EndPoint)
		if err != nil {
			return &outputs, err
		}

		if err := validateCompressedFile(packageName); err != nil {
			return &outputs, err
		}

		fullPath := getDecompressDirName(packageName)
		if err = isDirExist(fullPath); err != nil {
			// comporessedFileFullPath, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
			comporessedFileFullPath, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
			if err != nil {
				logrus.Errorf("ListCurrentDirAction downloadS3File fullPath=%v,err=%v", comporessedFileFullPath, err)
				return &outputs, err
			}

			if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
				logrus.Errorf("ListCurrentDirAction decompressFile fullPath=%v,err=%v", comporessedFileFullPath, err)
				os.RemoveAll(comporessedFileFullPath)
				return &outputs, err
			}
			os.RemoveAll(comporessedFileFullPath)
		}

		nodes, err := listCurrentDirectory(fullPath + "/" + input.CurrentDir)
		if err != nil {
			return &outputs, err
		}

		output := ListFilesOutput{
			Files: nodes,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

type GetConfigFileKeyInputs struct {
	Inputs []GetConfigFileKeyInput `json:"inputs,omitempty"`
}

type GetConfigFileKeyInput struct {
	CallBackParameter
	EndPoint string `json:"endpoint,omitempty"`
	FilePath string `json:"filePath,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
}

type GetConfigFileKeyOutputs struct {
	Outputs []GetConfigFileKeyOutput `json:"outputs,omitempty"`
}

type GetConfigFileKeyOutput struct {
	CallBackParameter
	FilePath       string          `json:"filePath,omitempty"`
	ConfigKeyInfos []ConfigKeyInfo `json:"configKeyInfos"`
}

type ConfigKeyInfo struct {
	Line string `json:"line,omitempty"`
	Key  string `json:"key,omitempty"`
}

type GetConfigFileKeyAction struct {
}

func (action *GetConfigFileKeyAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs GetConfigFileKeyInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetConfigFileKeyAction) CheckParam(input interface{}) error {
	inputs, ok := input.(GetConfigFileKeyInputs)
	if !ok {
		return fmt.Errorf("ListCurrentDirAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.FilePath == "" {
			return errors.New("FilePath is empty")
		}
		if input.EndPoint == "" {
			return errors.New("Endpoint is empty")
		}
	}

	return nil
}

func (action *GetConfigFileKeyAction) Do(input interface{}) (interface{}, error) {
	outputs := GetConfigFileKeyOutputs{}
	inputs, ok := input.(GetConfigFileKeyInputs)
	if !ok {
		return &outputs, fmt.Errorf("GetConfigFileKeyAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		output := GetConfigFileKeyOutput{}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		packageName, err := getPackageNameFromEndpoint(input.EndPoint)
		if err != nil {
			return &outputs, err
		}
		logrus.Info("package name = >", packageName)

		fullPath := getDecompressDirName(packageName)
		if err = isDirExist(fullPath); err != nil {
			// comporessedFileFullPath, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
			comporessedFileFullPath, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
			if err != nil {
				logrus.Errorf("GetConfigFileKeyAction downloadS3File fullPath=%v,err=%v", comporessedFileFullPath, err)
				return &outputs, err
			}

			if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
				logrus.Errorf("GetConfigFileKeyAction decompressFile fullPath=%v,err=%v", comporessedFileFullPath, err)
				os.RemoveAll(comporessedFileFullPath)
				return &outputs, err
			}
			os.RemoveAll(comporessedFileFullPath)
		}
		logrus.Info("full path = >", fullPath)
		keys, err := GetVariable(fullPath, input.FilePath)
		if err != nil {
			return nil, err
		}

		output.FilePath = input.FilePath
		output.ConfigKeyInfos = keys
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, nil
}
