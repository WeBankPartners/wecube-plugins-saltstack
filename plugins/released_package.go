package plugins

import (
	"fmt"
	"os"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
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
	Language string
}

func (action *ListCurrentDirAction) SetAcceptLanguage(language string) {
	action.Language = language
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
		return getParamEmptyError(action.Language, "endpoint")
	}
	return nil
}

func getPackageNameFromEndpoint(endpoint string) (string, error) {
	index := strings.LastIndexAny(endpoint, "/")
	if index == -1 {
		return "", fmt.Errorf("Invalid endpoint %s ", endpoint)
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
			if strings.HasPrefix(err.Error(), "exist") {
				output.Result.Code = "2"
			}else {
				output.Result.Code = RESULT_CODE_ERROR
			}
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
		err = getDecompressSuffixError(action.Language, packageName)
		return output, err
	}

	fullPath := getDecompressDirName(packageName)
	if err = isDirExist(fullPath); err != nil {
		comporessedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, true, action.Language)
		if err != nil {
			return output, err
		}

		if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
			err = getUnpackFileError(action.Language, comporessedFileFullPath, err)
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
			log.Logger.Error("List current dir,recover error", log.String("error", fmt.Sprintf("err=%v", err)))
		}
	}()

	for _, input := range inputs.Inputs {
		output, err := action.listCurrentDir(&input)
		if err != nil {
			log.Logger.Error("List current dir action", log.Error(err))
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
	Type string `json:"type,omitempty"`
}

type GetConfigFileKeyAction struct {
	Language string
}

func (action *GetConfigFileKeyAction) SetAcceptLanguage(language string) {
	action.Language = language
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
		return getParamEmptyError(action.Language, "filePath")
	}
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
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
	log.Logger.Debug("Package name", log.String("name", packageName))

	fullPath := getDecompressDirName(packageName)
	if err = isDirExist(fullPath); err != nil {
		comporessedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, true, action.Language)
		if err != nil {
			return output, err
		}

		if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
			err = getUnpackFileError(action.Language, comporessedFileFullPath, err)
			os.RemoveAll(comporessedFileFullPath)
			return output, err
		}
		os.RemoveAll(comporessedFileFullPath)
	}

	if fullPath[len(fullPath)-1] == '/' {
		fullPath = fullPath[:len(fullPath)-1]
	}
	if input.FilePath[0] == '/' {
		input.FilePath = input.FilePath[1:]
	}
	log.Logger.Debug("ConfigFile", log.String("file", fullPath+"/"+input.FilePath))
	tmpSpecialReplaceList := DefaultSpecialReplaceList
	tmpSpecialReplaceList = append(tmpSpecialReplaceList, DefaultEncryptReplaceList...)
	tmpSpecialReplaceList = append(tmpSpecialReplaceList, DefaultFileReplaceList...)
	keys, err := GetVariable(fullPath + "/" + input.FilePath, tmpSpecialReplaceList, true)
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
			log.Logger.Error("Get config file key action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
