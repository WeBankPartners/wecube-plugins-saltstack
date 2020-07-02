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
	Guid       string `json:"guid,omitempty"`
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
	Result
	Guid  string     `json:"guid,omitempty"`
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

func (action *ListCurrentDirAction) CheckParam(input ListFilesInput) error {
	if input.EndPoint == "" {
		return errors.New("Endpoint is empty")
	}

	// if input.AccessKey == "" {
	// 	return errors.New("AccessKey is empty")
	// }

	// if input.SecretKey == "" {
	// 	return errors.New("SecretKey is empty")
	// }
	return nil
}

func getPackageNameFromEndpoint(endpoint string) (string, error) {
	index := strings.LastIndexAny(endpoint, "/")
	if index == -1 {
		return "", fmt.Errorf("Invalid endpoint %s", endpoint)
	}

	return endpoint[index+1:], nil
}

func (action *ListCurrentDirAction) listCurrentDir(input *ListFilesInput) (output ListFilesOutput, err error) {
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

	packageName, err := getPackageNameFromEndpoint(input.EndPoint)
	if err != nil {
		return output, err
	}

	if err := validateCompressedFile(packageName); err != nil {
		return output, err
	}

	fullPath := getDecompressDirName(packageName)
	if err = isDirExist(fullPath); err != nil {
		// comporessedFileFullPath, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
		comporessedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, true)
		if err != nil {
			logrus.Errorf("ListCurrentDirAction downloadS3File fullPath=%v,err=%v", comporessedFileFullPath, err)
			return output, err
		}

		if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
			logrus.Errorf("ListCurrentDirAction decompressFile fullPath=%v,err=%v", comporessedFileFullPath, err)
			os.RemoveAll(comporessedFileFullPath)
			return output, err
		}
		os.RemoveAll(comporessedFileFullPath)
	}

	nodes, err := listCurrentDirectory(fullPath + "/" + input.CurrentDir)
	if err != nil {
		return output, err
	}
	output.Files = nodes

	return output, err
}

func (action *ListCurrentDirAction) Do(input interface{}) (interface{}, error) {
	outputs := ListFilesOutputs{}
	var finalErr error
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
		output, err := action.listCurrentDir(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

type GetConfigFileKeyInputs struct {
	Inputs []GetConfigFileKeyInput `json:"inputs,omitempty"`
}

type GetConfigFileKeyInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
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
	Result
	Guid           string          `json:"guid,omitempty"`
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

func (action *GetConfigFileKeyAction) CheckParam(input GetConfigFileKeyInput) error {
	if input.FilePath == "" {
		return errors.New("FilePath is empty")
	}
	if input.EndPoint == "" {
		return errors.New("Endpoint is empty")
	}

	return nil
}

func (action *GetConfigFileKeyAction) getConfigFileKey(input *GetConfigFileKeyInput) (output GetConfigFileKeyOutput, err error) {
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

	packageName, err := getPackageNameFromEndpoint(input.EndPoint)
	if err != nil {
		return output, err
	}
	logrus.Info("package name = >", packageName)

	fullPath := getDecompressDirName(packageName)
	if err = isDirExist(fullPath); err != nil {
		// comporessedFileFullPath, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
		comporessedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, true)
		if err != nil {
			logrus.Errorf("GetConfigFileKeyAction downloadS3File fullPath=%v,err=%v", comporessedFileFullPath, err)
			return output, err
		}

		if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
			logrus.Errorf("GetConfigFileKeyAction decompressFile fullPath=%v,err=%v", comporessedFileFullPath, err)
			os.RemoveAll(comporessedFileFullPath)
			return output, err
		}
		os.RemoveAll(comporessedFileFullPath)
	}
	logrus.Info("full path = >", fullPath)

	if fullPath[len(fullPath)-1] == '/' {
		fullPath = fullPath[:len(fullPath)-1]
	}
	if input.FilePath[0] == '/' {
		input.FilePath = input.FilePath[1:]
	}

	logrus.Infof("ConfigFile=%v", fullPath+"/"+input.FilePath)
	keys, err := GetVariable(fullPath + "/" + input.FilePath, DefaultSpecialReplaceList)
	if err != nil {
		return output, err
	}

	output.FilePath = input.FilePath
	output.ConfigKeyInfos = keys

	return output, err
}

func (action *GetConfigFileKeyAction) Do(input interface{}) (interface{}, error) {
	outputs := GetConfigFileKeyOutputs{}
	inputs, ok := input.(GetConfigFileKeyInputs)
	if !ok {
		return &outputs, fmt.Errorf("GetConfigFileKeyAction:input type=%T not right", input)
	}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.getConfigFileKey(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
