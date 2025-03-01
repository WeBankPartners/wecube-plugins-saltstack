package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var FileActions = make(map[string]Action)

func init() {
	FileActions["copy"] = new(FileCopyAction)
	FileActions["find"] = new(FileFindAction)
	FileActions["create"] = new(FileCreateAction)
	FileActions["delete"] = new(FileDeleteAction)
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

type FileCopyThreadObj struct {
	Data  FileCopyOutput
	Err   error
	Index int
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
	directory := ""
	if lastIndex := strings.LastIndex(input.DestinationPath, "/"); lastIndex >= 0 {
		directory = input.DestinationPath[0:lastIndex]
	} else {
		return fmt.Errorf("destinationPath:%s illegal with absolute path check ", input.DestinationPath)
	}
	directory = strings.TrimSpace(directory)
	if directory == "" {
		return fmt.Errorf("changeDirectoryOwner destinationPath can not empty")
	}
	//cmdRun := "chown -R " + input.FileOwner + "  " + directory
	cmdRun := "chown -R " + input.FileOwner + "  " + directory + "/ && chown -R " + input.FileOwner + "  " + directory
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

	fileName, tmpErr := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, true, action.Language)
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
		unpackOutput, unpackErr := CallSaltApi("https://127.0.0.1:8080", *unpackRequest, action.Language)
		if unpackErr != nil {
			err = unpackErr
			return output, err
		}
		if strings.Contains(unpackOutput, "'archive.cmd_unzip' not found") || strings.Contains(unpackOutput, "'archive.tar' not found") {
			err = fmt.Errorf("can not find unzip or tar command in target host")
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
	currentDirectory := ""
	if lastIndex := strings.LastIndex(input.DestinationPath, "/"); lastIndex >= 0 {
		currentDirectory = input.DestinationPath[0:lastIndex]
	} else {
		return &request, fmt.Errorf("destinationPath:%s illegal with absolute path check ", input.DestinationPath)
	}
	//currentDirectory := input.DestinationPath[0:strings.LastIndex(input.DestinationPath, "/")]

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
	outputChan := make(chan FileCopyThreadObj, len(files.Inputs))
	concurrentChan := make(chan int, ApiConcurrentNum)
	wg := sync.WaitGroup{}
	for i, file := range files.Inputs {
		concurrentChan <- 1
		wg.Add(1)
		go func(tmpInput FileCopyInput, index int) {
			output, err := action.copyFile(&tmpInput)
			outputChan <- FileCopyThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
			<-concurrentChan
		}(file, i)
		outputs.Outputs = append(outputs.Outputs, FileCopyOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("File copy action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}

	return &outputs, finalErr
}

// Create file plugin

type FileCreateInputs struct {
	Inputs []FileCreateInput `json:"inputs,omitempty"`
}

type FileCreateInput struct {
	CallBackParameter
	Guid            string `json:"guid,omitempty"`
	Target          string `json:"target,omitempty"`
	FileContent     string `json:"fileContent,omitempty"`
	DestinationPath string `json:"destinationPath,omitempty"`
	FileOwner       string `json:"fileOwner,omitempty"`
}

type FileCreateOutputs struct {
	Outputs []FileCreateOutput `json:"outputs,omitempty"`
}

type FileCreateOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type FileCreateThreadObj struct {
	Data  FileCreateOutput
	Err   error
	Index int
}

type FileCreateAction struct{ Language string }

func (action *FileCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FileCreateInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *FileCreateAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *FileCreateAction) CheckParam(input FileCreateInput) error {
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.FileContent == "" {
		return getParamEmptyError(action.Language, "fileContent")
	}
	if input.DestinationPath == "" {
		return getParamEmptyError(action.Language, "destinationPath")
	}

	return nil
}

func (action *FileCreateAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(FileCreateInputs)
	outputs := FileCreateOutputs{}
	var finalErr error
	outputChan := make(chan FileCreateThreadObj, len(files.Inputs))
	concurrentChan := make(chan int, ApiConcurrentNum)
	wg := sync.WaitGroup{}
	for i, file := range files.Inputs {
		concurrentChan <- 1
		wg.Add(1)
		go func(tmpInput FileCreateInput, index int) {
			output, err := action.createFile(&tmpInput)
			outputChan <- FileCreateThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
			<-concurrentChan
		}(file, i)
		outputs.Outputs = append(outputs.Outputs, FileCreateOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("File create action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}

	return &outputs, finalErr
}

func (action *FileCreateAction) deriveMd5SumRequest(input *FileCreateInput) (*SaltApiRequest, error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "file.get_hash"
	request.Args = append(request.Args, input.DestinationPath)
	request.Args = append(request.Args, "md5")

	return &request, nil
}

func (action *FileCreateAction) deriveCopyFileRequest(basePath string, input *FileCreateInput) (*SaltApiRequest, error) {

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

func (action *FileCreateAction) changeDirectoryOwner(input *FileCreateInput) error {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = input.Target
	request.Function = "cmd.run"
	if !strings.Contains(input.FileOwner, ":") {
		input.FileOwner = fmt.Sprintf("%s:%s", input.FileOwner, input.FileOwner)
	}
	directory := ""
	if lastIndex := strings.LastIndex(input.DestinationPath, "/"); lastIndex >= 0 {
		directory = input.DestinationPath[0:lastIndex]
	} else {
		return fmt.Errorf("destinationPath:%s illegal with absolute path check ", input.DestinationPath)
	}
	//directory := input.DestinationPath[0:strings.LastIndex(input.DestinationPath, "/")]
	cmdRun := "chown -R " + input.FileOwner + "  " + directory
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

func (action *FileCreateAction) createFile(input *FileCreateInput) (output FileCreateOutput, err error) {
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

	// Check user if not exist
	if input.FileOwner != "" {
		if userExist, errOut := checkRunUserIsExists(input.Target, input.FileOwner, action.Language); !userExist {
			return output, fmt.Errorf(errOut)
		}
	}

	// Write content to tmp file in salt base dir
	//tmpBaseDir := path.Join(SCRIPT_SAVE_PATH, input.Target)
	//tmpFile, err := os.CreateTemp(tmpBaseDir, "fast-tmp-")
	tmpFile, err := ioutil.TempFile(SCRIPT_SAVE_PATH, "fast-tmp-")
	if err != nil {
		return output, fmt.Errorf("new tmp file error, %s", err.Error())
	}
	log.Logger.Debug("Create tmp file success", log.String("tmpBaseDir", SCRIPT_SAVE_PATH),
		log.String("file", tmpFile.Name()))

	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err = tmpFile.Write([]byte(input.FileContent)); err != nil {
		return output, fmt.Errorf("write fileContent to tmp file error, %s", err.Error())
	}
	log.Logger.Debug("Write fileContent to tmp file success", log.String("file", tmpFile.Name()))

	// Convert path such as /srv/salt/base/target/t1 -> salt://base/target/t1
	paths := strings.Split(tmpFile.Name(), "base")
	saltPath := fmt.Sprintf("salt://base%s", paths[1])

	//md5sum, err := SendFile(saltPath, input.DestinationPath, input.FileOwner, input.Target)
	//if err != nil {
	//	return output, fmt.Errorf("send tmp file error,%s ", err.Error())
	//}

	//copy file
	copyRequest, err := action.deriveCopyFileRequest(saltPath, input)
	_, err = CallSaltApi("https://127.0.0.1:8080", *copyRequest, action.Language)
	if err != nil {
		return output, err
	}

	md5SumRequest, _ := action.deriveMd5SumRequest(input)
	md5sum, err := CallSaltApi("https://127.0.0.1:8080", *md5SumRequest, action.Language)
	if err != nil {
		return output, err
	}

	if input.FileOwner != "" {
		if err = action.changeDirectoryOwner(input); err != nil {
			return output, err
		}
	}

	output.Detail = md5sum

	return output, err
}

// Find file plugin

type FileFindInputs struct {
	Inputs []FileFindInput `json:"inputs,omitempty"`
}

type FileFindInput struct {
	CallBackParameter
	Guid        string `json:"guid,omitempty"`
	Target      string `json:"target,omitempty"`
	FilePath    string `json:"filePath,omitempty"`
	FilePattern string `json:"filePattern,omitempty"`
}

type FileFindOutputs struct {
	Outputs []FileFindOutput `json:"outputs,omitempty"`
}

type FileFindOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Files  string `json:"files"`
	Detail string `json:"detail,omitempty"`
}

type FileFindThreadObj struct {
	Data  FileFindOutput
	Err   error
	Index int
}

type FileFindAction struct{ Language string }

func (action *FileFindAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FileFindInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *FileFindAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *FileFindAction) CheckParam(input FileFindInput) error {
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if input.FilePath == "" {
		return getParamEmptyError(action.Language, "filePath")
	}
	if input.FilePattern == "" {
		return getParamEmptyError(action.Language, "FilePattern")
	}

	return nil
}

func (action *FileFindAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(FileFindInputs)
	outputs := FileFindOutputs{}
	var finalErr error
	outputChan := make(chan FileFindThreadObj, len(files.Inputs))
	concurrentChan := make(chan int, ApiConcurrentNum)
	wg := sync.WaitGroup{}
	for i, file := range files.Inputs {
		concurrentChan <- 1
		wg.Add(1)
		go func(tmpInput FileFindInput, index int) {
			output, err := action.findFile(&tmpInput)
			outputChan <- FileFindThreadObj{Data: output, Err: err, Index: index}
			wg.Done()
			<-concurrentChan
		}(file, i)
		outputs.Outputs = append(outputs.Outputs, FileFindOutput{})
	}
	wg.Wait()
	for {
		if len(outputChan) == 0 {
			break
		}
		tmpOutput := <-outputChan
		if tmpOutput.Err != nil {
			log.Logger.Error("File find action", log.Error(tmpOutput.Err))
			finalErr = tmpOutput.Err
		}
		outputs.Outputs[tmpOutput.Index] = tmpOutput.Data
	}

	return &outputs, finalErr
}

func (action *FileFindAction) findFile(input *FileFindInput) (output FileFindOutput, err error) {
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

	log.Logger.Debug("findFile in "+input.FilePath, log.String("pattern", input.FilePattern))
	ret, err := FindGlobFiles(input.FilePath, input.FilePattern, input.Target)
	if err != nil {
		return output, err
	}

	log.Logger.Debug("findFile in "+input.FilePath, log.String("files", ret))
	output.Files = ret
	return output, err
}

// File delete

type FileDeleteAction struct{ Language string }

type FileDeleteInputs struct {
	Inputs []FileDeleteInput `json:"inputs,omitempty"`
}

type FileDeleteInput struct {
	CallBackParameter
	Guid       string `json:"guid,omitempty"`
	Host       string `json:"host"`
	BaseDir    string `json:"baseDir"`
	TargetPath string `json:"targetPath"`
}

type FileDeleteOutputs struct {
	Outputs []FileDeleteOutput `json:"outputs,omitempty"`
}

type FileDeleteOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (action *FileDeleteAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs FileDeleteInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *FileDeleteAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *FileDeleteAction) CheckParam(input FileDeleteInput) error {
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.TargetPath == "" {
		return getParamEmptyError(action.Language, "target")
	}
	return nil
}

func (action *FileDeleteAction) Do(actionInput interface{}) (interface{}, error) {
	param, _ := actionInput.(FileDeleteInputs)
	outputs := FileDeleteOutputs{}
	var finalErr error
	for _, input := range param.Inputs {
		fileDeleteOutput, err := action.deleteFile(&input)
		if err != nil {
			log.Logger.Error("File delete action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, fileDeleteOutput)
	}

	return &outputs, finalErr
}

func (action *FileDeleteAction) deleteTargetPath(targetHost, targetPath string) (output string, err error) {
	request := SaltApiRequest{}
	request.Client = "local"
	request.TargetType = "ipcidr"
	request.Target = targetHost
	request.Function = "cmd.run"
	cmdRun := "rm -rf " + targetPath
	request.Args = append(request.Args, cmdRun)
	output, err = CallSaltApi("https://127.0.0.1:8080", request, action.Language)
	return
}

func (action *FileDeleteAction) deleteFile(input *FileDeleteInput) (output FileDeleteOutput, err error) {
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

	targetPathList := splitWithCustomFlag(input.TargetPath)
	if input.BaseDir != "" {
		if !strings.HasSuffix(input.BaseDir, "/") {
			input.BaseDir = input.BaseDir + "/"
		}
		var newPathList []string
		for _, v := range targetPathList {
			if strings.HasPrefix(v, "/") {
				v = v[1:]
			}
			newPathList = append(newPathList, input.BaseDir+v)
		}
		targetPathList = newPathList
	}
	var newPathList []string
	for _, v := range targetPathList {
		tmpPath := strings.TrimSpace(v)
		if tmpPath == "/" || strings.Contains(tmpPath, "..") {
			err = fmt.Errorf("path illegal")
			return output, err
		}
		if strings.HasSuffix(tmpPath, "/") {
			tmpPath = tmpPath + "*"
		}
		newPathList = append(newPathList, tmpPath)
	}

	for _, v := range newPathList {
		tmpOutput, tmpErr := action.deleteTargetPath(input.Host, v)
		if tmpErr != nil {
			err = tmpErr
			output.Detail = fmt.Sprintf("try to delete %s from host:%s fail,%s ", v, input.Host, tmpOutput)
			return output, err
		}
	}

	return output, err
}
