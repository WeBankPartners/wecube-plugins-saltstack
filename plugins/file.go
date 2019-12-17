package plugins

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
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

type FileCopyAction struct {
}

func (action *FileCopyAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FileCopyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *FileCopyAction) CheckParam(input FileCopyInput) error {
	if input.EndPoint == "" {
		return errors.New("EndPoint is empty")
	}
	if input.Target == "" {
		return errors.New("Target is empty")
	}
	if input.DestinationPath == "" {
		return errors.New("DestinationPath is empty")
	}

	// if input.SecretKey == "" {
	// 	return errors.New("SecretKey is empty")
	// }
	// if input.AccessKey == "" {
	// 	return errors.New("AccessKey is empty")
	// }

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

func changeDirecoryOwner(input *FileCopyInput) error {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "cmd.run"

	directory := input.DestinationPath[0:strings.LastIndex(input.DestinationPath, "/")]
	cmdRun := "chown -R " + input.FileOwner + "  " + directory
	request.Args = append(request.Args, cmdRun)

	_, err := CallSaltApi("https://127.0.0.1:8080", request)
	if err != nil {
		return err
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

	fileName, er := downloadS3File(input.EndPoint, "access_key", "secret_key")
	if er != nil {
		logrus.Errorf("CopyFile downloads3 file error=%v", er)
		err = fmt.Errorf("CopyFile downloads3 file error=%v", er)
		return output, err
	}

	input.DestinationPath = buildFileDestinationPath(input.EndPoint, input.DestinationPath)

	savePath, er := saveFileToSaltMasterBaseDir(fileName)
	os.Remove(fileName)
	if er != nil {
		err = fmt.Errorf("saveFileToSaltMasterBaseDir meet error=%v", err)
		return output, err
	}

	//copy file
	copyRequest, err := action.deriveCopyFileRequest("salt://base/"+filepath.Base(savePath), input)
	_, err = CallSaltApi("https://127.0.0.1:8080", *copyRequest)
	os.Remove(savePath)
	if err != nil {
		return output, err
	}

	md5SumRequest, _ := action.deriveMd5SumRequest(input)
	md5sum, err := CallSaltApi("https://127.0.0.1:8080", *md5SumRequest)
	if err != nil {
		return output, err
	}

	var unpackRequest *SaltApiRequest
	if input.Unpack == "true" {
		unpackRequest, err = action.deriveUnpackRequest(input)
		if err != nil {
			return output, err
		}

		if _, err = CallSaltApi("https://127.0.0.1:8080", *unpackRequest); err != nil {
			return output, err
		}
		if input.FileOwner != "" {
			if err = changeDirecoryOwner(input); err != nil {
				return output, err
			}
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
	} else if strings.HasSuffix(lowerFilepath, ".rar") {
		request.Function = "archive.unrar"
		request.Args = append(request.Args, input.DestinationPath)
		request.Args = append(request.Args, currentDirectory)
	} else if strings.HasSuffix(lowerFilepath, ".tar") {
		request.Function = "archive.tar"
		request.Args = append(request.Args, "xf")
		request.Args = append(request.Args, input.DestinationPath)
		request.Args = append(request.Args, "dest="+currentDirectory)
	} else if strings.HasSuffix(lowerFilepath, ".tar.gz") || strings.HasSuffix(lowerFilepath, ".tgz") {
		request.Function = "archive.tar"
		request.Args = append(request.Args, "zxf")
		request.Args = append(request.Args, input.DestinationPath)
		request.Args = append(request.Args, "dest="+currentDirectory)
	} else if strings.HasSuffix(lowerFilepath, ".gz") {
		request.Function = "archive.gunzip"
		request.Args = append(request.Args, input.DestinationPath)
	} else {
		return &request, fmt.Errorf("%s has invalid compressed format", lowerFilepath)
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
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, fileCopyOutput)
	}

	logrus.Infof("all files = %v are copied", files)
	return &outputs, finalErr
}
