package plugins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

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
	CallBackParameter
	EndPoint     string `json:"endpoint,omitempty"`
	Guid         string `json:"guid,omitempty"`
	FilePath     string `json:"confFiles,omitempty"`
	VariableList string `json:"variableList,omitempty"`
	// AccessKey    string `json:"accessKey,omitempty"`
	// SecretKey    string `json:"secretKey,omitempty"`
}

//VariableReplaceOutputs .
type VariableReplaceOutputs struct {
	Outputs []VariableReplaceOutput `json:"outputs,omitempty"`
}

//VariableReplaceOutput .
type VariableReplaceOutput struct {
	CallBackParameter
	Guid         string `json:"guid,omitempty"`
	NewS3PkgPath string `json:"s3PkgPath,omitempty"`
	//Detail     string `json:"detail,omitempty"`
	//MD5        string `json:"md5,omitempty"`
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
	if len(data.Inputs) > 0 {
		for _, d := range data.Inputs {
			if d.EndPoint == "" || d.FilePath == "" || d.VariableList == "" {
				return fmt.Errorf("VariableReplaceAction endpoint, file_path, variable_list could not be empty")
			}
			if !strings.Contains(d.VariableList, "=") {
				return fmt.Errorf("VariableReplaceAction input variable don't have '=' could't get variable key value pair")
			}
		}
	}

	return nil
}

func getNewS3EndpointName(endpoint string, newPackageName string) string {
	index := strings.LastIndexAny(endpoint, "/")
	return endpoint[0:index+1] + newPackageName
}

func getPackageNameWithoutSuffix(packageName string) string {
	index := strings.LastIndexAny(packageName, ".")
	return packageName[0:index]
}

func (action *VariableReplaceAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(VariableReplaceInputs)
	outputs := VariableReplaceOutputs{}

	for _, input := range files.Inputs {
		suffix, err := getCompressFileSuffix(input.EndPoint)
		if err != nil {
			return &outputs, err
		}

		packageName, err := getPackageNameFromEndpoint(input.EndPoint)
		if err != nil {
			return &outputs, err
		}
		logrus.Info("package name = >", packageName)

		decompressDirName := getDecompressDirName(packageName)
		if err = isDirExist(decompressDirName); err == nil {
			os.RemoveAll(decompressDirName)
		}

		if err = os.MkdirAll(decompressDirName, os.ModePerm); err != nil {
			return &outputs, err
		}

		compressedFileFullPath, err := downloadS3File(input.EndPoint, "access_key", "secret_key")
		if err != nil {
			logrus.Errorf("VariableReplaceAction downloadS3File fullPath=%v,err=%v", compressedFileFullPath, err)
			return &outputs, err
		}

		if err = decompressFile(compressedFileFullPath, decompressDirName); err != nil {
			logrus.Errorf("VariableReplaceAction decompressFile fullPath=%v,err=%v", compressedFileFullPath, err)
			os.RemoveAll(compressedFileFullPath)
			return &outputs, err
		}
		os.RemoveAll(compressedFileFullPath)

		for _, filePath := range strings.Split(input.FilePath, "|") {
			confFilePath := decompressDirName + "/" + filePath
			if err := ReplaceFileVar(confFilePath, input.VariableList); err != nil {
				os.RemoveAll(decompressDirName)
				return &outputs, err
			}
		}

		//compress file
		nowTime := time.Now().Format("200601021504")
		newPackageName := fmt.Sprintf("%s-%v%s", getPackageNameWithoutSuffix(packageName), nowTime, suffix)
		fmt.Printf("newPackageName=%s\n", newPackageName)
		if err = compressDir(decompressDirName, suffix, newPackageName); err != nil {
			logrus.Errorf("compressDir meet error=%v", err)
			os.RemoveAll(decompressDirName)
			return &outputs, err
		}
		os.RemoveAll(decompressDirName)

		//upload to s3
		newS3Endpoint := getNewS3EndpointName(input.EndPoint, newPackageName)
		fmt.Printf("NewS3EndpointName=%s\n", newS3Endpoint)

		if _, err = uploadS3File(newS3Endpoint, "access_key", "secret_key"); err != nil {
			logrus.Errorf("uploadS3File meet error=%v", err)
			return &outputs, err
		}
		output := VariableReplaceOutput{
			Guid:         input.Guid,
			NewS3PkgPath: newS3Endpoint,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

func ReplaceFileVar(filepath, variablelist string) error {
	index := strings.LastIndexAny(filepath, "/")
	if index == -1 {
		return fmt.Errorf("Invalid endpoint %s", filepath)
	}
	fileName := filepath[index+1:]
	fileVarList, err := GetFileVariableString(filepath, fileName)
	if err != nil {
		return err
	}

	if len(fileVarList) == 0 {
		return fmt.Errorf("file %s no variable need to replace", fileName)
	}

	keyMap, err := GetInputVariableMap(variablelist)
	if err != nil {
		logrus.Errorf("GetInputVariableMap error: %s", err)
		return err
	}

	err = CheckVariableIsAllReady(keyMap, fileVarList)
	if err != nil {
		logrus.Errorf("CheckVariableIsAllReady error: %s", err)
		return err
	}

	err = replaceFileVar(keyMap, filepath)
	if err != nil {
		logrus.Errorf("replaceFileVar error: %s", err)
		return err
	}

	return nil
}

//var1,var2=value1,value2  to  var1=value1, var2=value2
func changeVariableListFormat(variableList string) (string, error) {
	index := strings.Index(variableList, "=")
	if index == -1 {
			return "", fmt.Errorf("variableList(%s) do not have =", variableList)
	}

	varsStr := variableList[0:index]
	valuesStr := variableList[index+1 : len(variableList)]
	vars := strings.Split(varsStr, ",")
	values := strings.Split(valuesStr, ",")
	if len(vars) == 0 || len(values) == 0 || len(vars) != len(values) {
			return "", fmt.Errorf("len(vars)=%v,len(values)=%v", len(vars), len(values))
	}

	result := ""
	for i, _ := range vars {
			text := vars[i] + "=" + values[i]
			result += text
			if i != len(vars)-1 {
					result += ","
			}
	}
	return result, nil
}


func GetInputVariableMap(variable string) (map[string]string, error) {
	if !strings.Contains(variable, "=") {
		return nil, fmt.Errorf("input variable don't have '=' could't get variable key value pair")
	}

	newVariableList,err := changeVariableListFormat(strings.Replace(variable, " ", "", -1))
	if err != nil {
		return nil, err
	}

	inputVariableMap := make(map[string]string)
	if strings.Contains(newVariableList, ",") {
		str2 := strings.Split(newVariableList, ",")
		for _, v := range str2 {
			str3 := strings.Split(v, "=")
			inputVariableMap[str3[0]] = str3[1]
		}
	} else {
		str2 := strings.Split(newVariableList, "=")
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

func compressDir(decompressDirName string, suffix string, newPackageName string) error {
	sh := ""
	if suffix != ".zip" && suffix != ".tgz" && suffix != ".tar.gz" {
		return fmt.Errorf("%s is invalid suffix", suffix)
	}

	if suffix == ".zip" {
		sh = "cd " + decompressDirName + " && " + "zip -r " + UPLOADS3FILE_DIR + newPackageName + " *"
	}
	if suffix == ".tgz" || suffix == ".tar.gz" {
		sh = "cd " + decompressDirName + " && " + "tar czf  " + UPLOADS3FILE_DIR + newPackageName + " *"
	}
	fmt.Printf("compressDir sh=%s\n", sh)

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
