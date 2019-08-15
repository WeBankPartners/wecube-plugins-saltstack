package plugins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

//VariableActions .
var VariableActions = make(map[string]Action)

func init() {
	VariableActions["replace"] = new(VariableReplaceAction)
}

//VariablePlugin .
type VariablePlugin struct {
}

//GetActionByName .
func (plugin *VariablePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := VariableActions[actionName]

	if !found {
		return nil, fmt.Errorf("File plugin,action = %s not found", actionName)
	}

	return action, nil
}

//VariableReplaceInputs .
type VariableReplaceInputs struct {
	Inputs []VariableReplaceInput `json:"inputs,omitempty"`
}

//VariableReplaceInput .
type VariableReplaceInput struct {
	EndPoint string `json:"endpoint,omitempty"`
	// AccessKey    string `json:"accessKey,omitempty"`
	// SecretKey    string `json:"secretKey,omitempty"`
	Guid         string `json:"guid,omitempty"`
	PkgName      string `json:"pkg_name,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
	VariableList string `json:"variable_list,omitempty"`
}

//VariableReplaceOutputs .
type VariableReplaceOutputs struct {
	Outputs []VariableReplaceOutput `json:"outputs,omitempty"`
}

//VariableReplaceOutput .
type VariableReplaceOutput struct {
	Guid       string `json:"guid,omitempty"`
	Detail     string `json:"detail,omitempty"`
	MD5        string `json:"md5,omitempty"`
	FilePath   string `json:"file_path,omitempty"`
	CosPkgPath string `json:"cos_pkg_path,omitempty"`
}

type VariableReplaceAction struct {
}

func (action *VariableReplaceAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VariableReplaceInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VariableReplaceAction) CheckParam(input interface{}) error {
	data, ok := input.(VariableReplaceInputs)
	if !ok {
		return fmt.Errorf("VariableReplaceAction:input type=%T not right", input)
	}
	if !strings.Contains(data.Inputs[0].PkgName, ".zip") && !strings.Contains(data.Inputs[0].PkgName, ".tar.gz") && !strings.Contains(data.Inputs[0].PkgName, ".tgz") {
		return fmt.Errorf("VariableReplaceAction only support zip and tar.gz type package")
	}
	if len(data.Inputs) > 0 {
		for _, d := range data.Inputs {
			if d.PkgName == "" || d.FilePath == "" || d.VariableList == "" {
				return fmt.Errorf("VariableReplaceAction pkg_name pkg_path file_path file_name variable_list could not be empty")
			}
			if !strings.Contains(d.VariableList, "=") {
				return fmt.Errorf("VariableReplaceAction input variable don't have '=' could't get variable key value pair")
			}
		}

	}

	return nil
}

func (action *VariableReplaceAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(VariableReplaceInputs)
	outputs := VariableReplaceOutputs{}

	//replace variable
	fileList := []string{}
	confPkg := ""
	dirPath := ""
	for _, input := range files.Inputs {
		packageName, err := getPackageNameFromEndpoint(input.EndPoint)
		if err != nil {
			return &outputs, err
		}
		logrus.Info("package name = >", packageName)
		fullPath := getDecompressDirName(packageName)
		dirPath = fullPath
		if err = isDirExist(fullPath); err != nil {
			// comporessedFileFullPath, err := downloadS3File(input.EndPoint, input.AccessKey, input.SecretKey)
			comporessedFileFullPath, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
			if err != nil {
				logrus.Errorf("VariableReplaceAction downloadS3File fullPath=%v,err=%v", comporessedFileFullPath, err)
				return &outputs, err
			}

			if err = decompressFile(comporessedFileFullPath, fullPath); err != nil {
				logrus.Errorf("VariableReplaceAction decompressFile fullPath=%v,err=%v", comporessedFileFullPath, err)
				os.RemoveAll(comporessedFileFullPath)
				return &outputs, err
			}
			os.RemoveAll(comporessedFileFullPath)
		}
		filePath := fullPath + "/" + input.FilePath
		fileList = append(fileList, input.FilePath)
		confPkg = fullPath + "/" + input.PkgName
		resp, err := ReplaceFileVar(filePath, input.VariableList)
		if err != nil {
			outputs.Outputs = append(outputs.Outputs, resp)
			return &outputs, err
		}
		resp.FilePath = input.FilePath
		resp.Guid = input.Guid
		md5, err := GetFileMD5Value(fullPath, input.FilePath)
		if err != nil {
			return outputs, err
		}
		resp.MD5 = md5
		outputs.Outputs = append(outputs.Outputs, resp)
	}
	if strings.Contains(files.Inputs[0].PkgName, ".zip") {
		err := CompressFile(dirPath, fileList, confPkg, "ZIP")
		if err != nil {
			logrus.Errorf("compress zip package error: %s", err)
			return &outputs, err
		}
	}
	if strings.Contains(files.Inputs[0].PkgName, ".tar.gz") || strings.Contains(files.Inputs[0].PkgName, ".tgz") {
		err := CompressFile(dirPath, fileList, confPkg, "TGZ")
		if err != nil {
			logrus.Errorf("compress tar.gz package error: %s", err)
			return &outputs, err
		}
	}
	s3Path, err := UploadConfPackage(files.Inputs[0])
	if err != nil {
		return &outputs, err
	}
	for i := 0; i < len(outputs.Outputs); i++ {
		outputs.Outputs[i].CosPkgPath = s3Path
	}

	logrus.Infof("all files = %v are finished", files)
	return &outputs, nil
}

func ReplaceFileVar(filepath, variablelist string) (VariableReplaceOutput, error) {

	var resp VariableReplaceOutput

	index := strings.LastIndexAny(filepath, "/")
	if index == -1 {
		return resp, fmt.Errorf("Invalid endpoint %s", filepath)
	}
	fileName := filepath[index+1:]
	fileVarList, err := GetFileVariableString(filepath, fileName)
	if err != nil {
		resp.Detail = "get " + fileName + " variable error"
		return resp, err
	}

	if len(fileVarList) == 0 {
		resp.Detail = "file " + fileName + " no variable need to replace"
		logrus.Errorf("file %s no variable need to replace", fileName)
		return resp, fmt.Errorf("file %s no variable need to replace", fileName)
	}

	keyMap, err := GetInputVariableMap(variablelist)
	if err != nil {
		logrus.Errorf("GetInputVariableMap error: %s", err)
		resp.Detail = "GetInputVariableMap error"
		return resp, err
	}

	err = CheckVariableIsAllReady(keyMap, fileVarList)
	if err != nil {
		logrus.Errorf("CheckVariableIsAllReady error: %s", err)
		resp.Detail = "CheckVariableIsAllReady error"
		return resp, err
	}

	err = replaceFileVar(keyMap, filepath)
	if err != nil {
		logrus.Errorf("replaceFileVar error: %s", err)
		resp.Detail = "replaceFileVar error"
		return resp, err
	}

	resp.Detail = "file " + fileName + " variable replace finished"

	return resp, nil
}

func GetInputVariableMap(variable string) (map[string]string, error) {

	if !strings.Contains(variable, "=") {
		return nil, fmt.Errorf("input variable don't have '=' could't get variable key value pair")
	}

	inputVariableMap := make(map[string]string)

	str1 := strings.Replace(variable, " ", "", -1)
	if strings.Contains(str1, ",") {
		str2 := strings.Split(str1, ",")
		for _, v := range str2 {
			str3 := strings.Split(v, "=")
			inputVariableMap[str3[0]] = str3[1]
		}

	} else {
		str2 := strings.Split(str1, "=")
		inputVariableMap[str2[0]] = str2[1]
	}

	return inputVariableMap, nil
}

func CheckVariableIsAllReady(input map[string]string, variablelist []string) (err error) {

	for _, va := range variablelist {
		if _, ok := input[va]; !ok {
			return fmt.Errorf("variable %s not input", va)
		}
	}

	return nil
}

func GetFileVariableString(filepath string, filename string) ([]string, error) {
	_, err := PathExists(filepath)
	if err != nil {
		logrus.Errorf("file %s not exits", filepath)
		return []string{}, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		logrus.Errorf("open file %s error", filepath)
		return []string{}, err
	}
	defer f.Close()

	br := bufio.NewReader(f)

	variablemap := make(map[string]string)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(a) == 0 {
			continue
		}

		flysnowRegexp := regexp.MustCompile(`[@*]\w+|[!*]\w+|[&*]\w+`)
		params := flysnowRegexp.FindAllString(string(a), -1)
		if len(params) > 0 {
			for _, param := range params {
				if strings.Contains(param, "@") {
					s := strings.Split(param, "@")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					variablemap[s[1]] = s[1]
				}
				if strings.Contains(param, "!") {
					s := strings.Split(param, "!")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					variablemap[s[1]] = s[1]
				}
				if strings.Contains(param, "&") {
					s := strings.Split(param, "&")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					variablemap[s[1]] = s[1]
				}
			}
		}
	}
	variableList := []string{}
	if len(variablemap) == 0 {
		logrus.Errorf("file %s don't hava variable need to replace", filepath)

		return []string{}, fmt.Errorf("file %s don't hava variable need to replace", filepath)
	}
	for _, v := range variablemap {
		variableList = append(variableList, v)
	}

	return variableList, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func replaceFileVar(keyMap map[string]string, filepath string) error {
	bf, err := os.Open(filepath)
	if err != nil {
		logrus.Errorf("open file fail: %s", err)
		return err
	}
	defer bf.Close()
	newfilePath := filepath + ".bak"
	f, err := os.OpenFile(newfilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Errorf("open file error: %s", err)
		return err
	}
	defer f.Close()
	br := bufio.NewReader(bf)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("read file line info err: %s", err)
			return err
		}
		newLine := string(line)
		flysnowRegexp := regexp.MustCompile(`[@*]\w+|[!*]\w+|[&*]\w+`)
		keys := flysnowRegexp.FindAllString(string(line), -1)
		if len(keys) > 0 {
			for _, key := range keys {
				if strings.Contains(key, "@") {
					s := strings.Split(key, "@")
					if s[1] == "" {
						return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
					}
					oldStr := "[" + key + "]"
					newLine = strings.Replace(newLine, oldStr, keyMap[s[1]], -1)
				}
				if strings.Contains(key, "!") {
					s := strings.Split(key, "!")
					if s[1] == "" {
						return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
					}
					oldStr := "[" + key + "]"
					newLine = strings.Replace(newLine, oldStr, keyMap[s[1]], -1)
				}
				if strings.Contains(key, "&") {
					s := strings.Split(key, "&")
					if s[1] == "" {
						return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
					}
					oldStr := "[" + key + "]"
					newLine = strings.Replace(newLine, oldStr, keyMap[s[1]], -1)
				}
			}
		}
		_, err = f.WriteString(newLine + "\n")
		if err != nil {
			logrus.Errorf("write to file fail: %s", err)
			return err
		}
	}
	err = os.Rename(newfilePath, filepath)
	if err != nil {
		logrus.Errorf("file rename Error: %s", err)
		return err
	}

	return nil
}

func CompressFile(dir string, filePath []string, pkgName string, pkgType string) error {
	sh := ""
	if pkgType == "ZIP" {
		sh = "cd " + dir + " && zip -r " + pkgName
		for _, file := range filePath {
			sh += " " + file
		}
		sh += " && cp " + pkgName + " " + UPLOADS3FILE_DIR
	}
	if pkgType == "TGZ" {
		sh = "cd " + dir + " && tar -zcvf " + pkgName
		for _, file := range filePath {
			sh += " " + file
		}
		sh += " && cp " + pkgName + " " + UPLOADS3FILE_DIR
	}
	cmd := exec.Command("/bin/sh", "-c", sh)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("can not obtain stdout pipe for command: %s \n", err)
		return err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("conmand start is error: %s \n", err)
		return err
	}
	_, err = LogReadLine(cmd, stdout)
	if err != nil {
		return err
	}

	return nil
}

func LogReadLine(cmd *exec.Cmd, stdout io.ReadCloser) ([]string, error) {

	linelist := []string{}
	outputBuf := bufio.NewReader(stdout)
	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			if err.Error() != "EOF" {
				logrus.Info("readline is error")
				return []string{}, nil
			}
		}

		linelist = append(linelist, string(output))
	}
	if err := cmd.Wait(); err != nil {
		return []string{}, nil
	}

	return linelist, nil
}

func GetFileMD5Value(dir, filePath string) (string, error) {
	sh := "cd " + dir + " && md5sum " + filePath + " |awk '{print $1}'"
	cmd := exec.Command("/bin/sh", "-c", sh)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("get file md5 can not obtain stdout pipe for command: %s \n", err)
		return "", err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("get file md5 conmand start is error: %s \n", err)
		return "", err
	}
	line, err := LogReadLine(cmd, stdout)
	if err != nil {
		return "", err
	}
	if len(line) == 0 {
		return "", fmt.Errorf("get file %s md5 failed", filePath)
	}

	return line[0], nil
}

func UploadConfPackage(input VariableReplaceInput) (string, error) {
	index := strings.LastIndex(input.EndPoint, "/")
	if index == -1 {
		return "", fmt.Errorf("endpoint %s is unvaliable", input.EndPoint)
	}
	point := input.EndPoint[:index]
	newEndPoint := point + "/" + input.PkgName
	logrus.Info("new package s3 path is ========================>>>>>>", newEndPoint)
	// _, err := uploadS3File(newEndPoint, input.AccessKey, input.SecretKey)
	_, err := uploadS3File(newEndPoint, "access_key", "secret_key")
	if err != nil {
		return "", err
	}

	return newEndPoint, nil
}
