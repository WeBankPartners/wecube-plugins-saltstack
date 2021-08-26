package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var FileActions = make(map[string]Action)

func init() {
	FileActions["copy"] = new(FileCopyAction)
}

type FilePlugin struct {
}

func (plugin *FilePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := FileActions[actionName]
	if !found {
		return nil, fmt.Errorf("File plugin,action = %s not found", actionName)
	}

	return action, nil
}

type FileCopyInputs struct {
	Inputs []FileCopyInput `json:"inputs,omitempty"`
}

type FileCopyInput struct {
	CallBackParameter
	EndPoint        string `json:"endpoint,omitempty"`
	Guid            string `json:"guid,omitempty"`
	Target          string `json:"target,omitempty"`
	DestinationPath string `json:"destinationPath,omitempty"`
	Unpack          string `json:"unpack,omitempty"`
	FileOwner       string `json:"fileOwner,omitempty"`
}

type FileCopyOutputs struct {
	Outputs []FileCopyOutput `json:"outputs,omitempty"`
}

type FileCopyOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type FileCopyAction struct{ Language string }

func (action *FileCopyAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FileCopyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *FileCopyAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *FileCopyAction) CheckParam(input FileCopyInput) error {
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.DestinationPath == "" {
		return getParamEmptyError(action.Language, "destinationPath")
	}
	if input.Unpack == "true" {
		if input.FileOwner == "" {
			return getParamEmptyError(action.Language, "fileOwner")
		}
	}

	return nil
}

func buildFileDestinationPath(endpoint string, destPath string) string {
	index := strings.LastIndexAny(destPath, "/")
	if index != len([]rune(destPath))-1 {
		return destPath
	}

	packageName, _ := getPackageNameFromEndpoint(endpoint)
	return destPath + packageName
}

func (action *FileCopyAction) changeDirectoryOwner(input *FileCopyInput) error {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "cmd.run"
	if !strings.Contains(input.FileOwner, ":") {
		input.FileOwner = fmt.Sprintf("%s:%s", input.FileOwner, input.FileOwner)
	}

	//directory := input.DestinationPath[0:strings.LastIndex(input.DestinationPath, "/")]
	cmdRun := "chown -R " + input.FileOwner + "  " + input.DestinationPath
	request.Args = append(request.Args, cmdRun)

	output, err := CallSaltApi("https://127.0.0.1:8080", request, action.Language)
	if err != nil {
		return err
	}
	log.Logger.Debug("Change dir owner", log.String("command", cmdRun), log.String("output", output))
	if strings.Contains(output, "chown") {
		return fmt.Errorf(output)
	}

	return nil
}

func (action *FileCopyAction) copyFile(input *FileCopyInput) (output FileCopyOutput, err error) {
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

	if input.FileOwner != "" {
		userExist, errOut := checkRunUserIsExists(input.Target, input.FileOwner, action.Language)
		if !userExist {
			err = fmt.Errorf(errOut)
			return output, err
		}
	}

	fileName, tmpErr := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false, action.Language)
	if tmpErr != nil {
		log.Logger.Error("Download s3 file", log.String("path", input.EndPoint), log.Error(tmpErr))
		err = tmpErr
		return output, err
	}

	input.DestinationPath = buildFileDestinationPath(input.EndPoint, input.DestinationPath)

	savePath, tmpErr := saveFileToSaltMasterBaseDir(fileName)
	os.Remove(fileName)
	if tmpErr != nil {
		err = getS3DownloadError(action.Language, input.EndPoint, fmt.Sprintf("move download file to salt-dir error:%s", tmpErr.Error()))
		return output, err
	}

	//copy file
	copyRequest, err := action.deriveCopyFileRequest("salt://base/"+filepath.Base(savePath), input)
	_, err = CallSaltApi("https://127.0.0.1:8080", *copyRequest, action.Language)
	os.Remove(savePath)
	if err != nil {
		return output, err
	}

	md5SumRequest, _ := action.deriveMd5SumRequest(input)
	md5sum, err := CallSaltApi("https://127.0.0.1:8080", *md5SumRequest, action.Language)
	if err != nil {
		return output, err
	}

	var unpackRequest *SaltApiRequest
	if input.Unpack == "true" {
		unpackRequest, err = action.deriveUnpackRequest(input)
		if err != nil {
			return output, err
		}

		if _, err = CallSaltApi("https://127.0.0.1:8080", *unpackRequest, action.Language); err != nil {
			return output, err
		}
	}
	if input.FileOwner != "" {
		if err = action.changeDirectoryOwner(input); err != nil {
			return output, err
		}
	}

	output.Detail = md5sum
	return output, err
}

func (action *FileCopyAction) deriveMd5SumRequest(input *FileCopyInput) (*SaltApiRequest, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "file.get_hash"
	request.Args = append(request.Args, input.DestinationPath)
	request.Args = append(request.Args, "md5")

	return &request, nil
}

func (action *FileCopyAction) deriveCopyFileRequest(basePath string, input *FileCopyInput) (*SaltApiRequest, error) {

	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "cp.get_file"

	request.Args = append(request.Args, basePath)
	request.Args = append(request.Args, input.DestinationPath)
	request.Args = append(request.Args, "makedirs=true")
	request.Args = append(request.Args, "gzip=5")

	return &request, nil
}

func (action *FileCopyAction) deriveUnpackRequest(input *FileCopyInput) (*SaltApiRequest, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target

	lowerFilepath := strings.ToLower(input.DestinationPath)
	currentDirectory := input.DestinationPath[0:strings.LastIndex(input.DestinationPath, "/")]

	if strings.HasSuffix(lowerFilepath, ".zip") {
		request.Function = "archive.cmd_unzip"
		request.Args = append(request.Args, input.DestinationPath)
		request.Args = append(request.Args, currentDirectory)
		request.Args = append(request.Args, "options=-o")
	} else if strings.HasSuffix(lowerFilepath, ".tar.gz") || strings.HasSuffix(lowerFilepath, ".tgz") {
		request.Function = "archive.tar"
		request.Args = append(request.Args, "zxf")
		// request.Args = append(request.Args, "--overwrite")
		request.Args = append(request.Args, input.DestinationPath)
		request.Args = append(request.Args, "dest="+currentDirectory)
	} else {
		return &request, getDecompressSuffixError(action.Language, lowerFilepath)
	}

	return &request, nil
}

func (action *FileCopyAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(FileCopyInputs)
	outputs := FileCopyOutputs{}
	var finalErr error
	for _, file := range files.Inputs {
		fileCopyOutput, err := action.copyFile(&file)
		if err != nil {
			log.Logger.Error("File copy action", log.JsonObj("param", file), log.Error(err))
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, fileCopyOutput)
	}

	return &outputs, finalErr
}
