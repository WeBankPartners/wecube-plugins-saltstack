package plugins

import (
	"bufio"
	"errors"
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

var (
	SEPERATOR                   = string([]byte{0x01})
	VARIABLE_KEY_SEPERATOR      = SEPERATOR + "=" + SEPERATOR
	VARIABLE_VARIABLE_SEPERATOR = "," + SEPERATOR
	KEY_KEY_SEPERATOR           = VARIABLE_VARIABLE_SEPERATOR
	ONE_VARIABLE_SEPERATOR      = "&" + SEPERATOR
)

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

	//support aomp password encrypt
	EncryptVariblePrefix string `json:"encryptVariblePrefix,omitempty"`
	Seed                 string `json:"seed,omitempty"`
	AppPublicKey         string `json:"appPublicKey,omitempty"`
	SysPrivateKey        string `json:"sysPrivateKey,omitempty"`
}

//VariableReplaceOutputs .
type VariableReplaceOutputs struct {
	Outputs []VariableReplaceOutput `json:"outputs,omitempty"`
}

//VariableReplaceOutput .
type VariableReplaceOutput struct {
	CallBackParameter
	Result
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

func (action *VariableReplaceAction) CheckParam(input VariableReplaceInput) error {
	if input.EndPoint == "" {
		return fmt.Errorf("VariableReplaceAction endpoint could not be empty")
	}
	if input.VariableList != "" {
		if !strings.Contains(input.VariableList, "=") {
			return fmt.Errorf("VariableReplaceAction input variable don't have '=' could't get variable key value pair")
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

func (action *VariableReplaceAction) variableReplace(input *VariableReplaceInput) (output VariableReplaceOutput, err error) {
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

	suffix, err := getCompressFileSuffix(input.EndPoint)
	if err != nil {
		return output, err
	}

	packageName, err := getPackageNameFromEndpoint(input.EndPoint)
	if err != nil {
		return output, err
	}
	logrus.Info("package name = >", packageName)

	decompressDirName := getDecompressDirName(packageName)
	if err = isDirExist(decompressDirName); err == nil {
		os.RemoveAll(decompressDirName)
	}

	if err = os.MkdirAll(decompressDirName, os.ModePerm); err != nil {
		return output, err
	}

	compressedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password)
	if err != nil {
		logrus.Errorf("VariableReplaceAction downloadS3File fullPath=%v,err=%v", compressedFileFullPath, err)
		return output, err
	}

	if err = decompressFile(compressedFileFullPath, decompressDirName); err != nil {
		logrus.Errorf("VariableReplaceAction decompressFile fullPath=%v,err=%v", compressedFileFullPath, err)
		os.RemoveAll(compressedFileFullPath)
		return output, err
	}
	os.RemoveAll(compressedFileFullPath)

	if input.FilePath != "" && input.VariableList != "" {
		for _, filePath := range strings.Split(input.FilePath, "|") {
			confFilePath := ""

			if decompressDirName[len(decompressDirName)-1] == '/' {
				decompressDirName = decompressDirName[:len(decompressDirName)-1]
			}

			if filePath[0] == '/' {
				confFilePath = decompressDirName + filePath
			} else {
				confFilePath = decompressDirName + "/" + filePath
			}

			logrus.Infof("confFilePath=%v", confFilePath)

			if err = ReplaceFileVar(confFilePath, input); err != nil {
				os.RemoveAll(decompressDirName)
				return output, err
			}
		}
	}

	//compress file
	nowTime := time.Now().Format("200601021504")
	newPackageName := fmt.Sprintf("%s-%v%s", getPackageNameWithoutSuffix(packageName), nowTime, suffix)
	fmt.Printf("newPackageName=%s\n", newPackageName)
	if err = compressDir(decompressDirName, suffix, newPackageName); err != nil {
		logrus.Errorf("compressDir meet error=%v", err)
		os.RemoveAll(decompressDirName)
		return output, err
	}
	os.RemoveAll(decompressDirName)

	//upload to s3
	newS3Endpoint := getNewS3EndpointName(input.EndPoint, newPackageName)
	logrus.Infof("NewS3EndpointName=%s\n", newS3Endpoint)

	if _, err = uploadS3File(newS3Endpoint, DefaultS3Key, DefaultS3Password); err != nil {
		logrus.Errorf("uploadS3File meet error=%v", err)
		return output, err
	}
	output.NewS3PkgPath = newS3Endpoint

	return output, err
}

func (action *VariableReplaceAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(VariableReplaceInputs)
	outputs := VariableReplaceOutputs{}
	var finalErr error

	for _, input := range files.Inputs {
		output, err := action.variableReplace(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

//variablelist,seed,publicKey,privateKey string
func ReplaceFileVar(filepath string, input *VariableReplaceInput) error {
	variablelist := input.VariableList
	seed := input.Seed
	publicKey := input.AppPublicKey
	privateKey := input.SysPrivateKey
	prefix := input.EncryptVariblePrefix

	index := strings.LastIndexAny(filepath, "/")
	if index == -1 {
		return fmt.Errorf("Invalid endpoint %s", filepath)
	}
	fileName := filepath[index+1:]
	fileVarMap, err := GetVariable(filepath)
	if err != nil {
		return err
	}

	if len(fileVarMap) == 0 {
		return fmt.Errorf("file %s no variable need to replace", fileName)
	}

	fileVarList := []string{}
	for _, v := range fileVarMap {
		fileVarList = append(fileVarList, v.Key)
	}

	keyMap, err := GetInputVariableMap(variablelist, seed)
	if err != nil {
		logrus.Errorf("GetInputVariableMap error: %s", err)
		return err
	}

	err = CheckVariableIsAllReady(keyMap, fileVarList)
	if err != nil {
		logrus.Errorf("CheckVariableIsAllReady error: %s", err)
		return err
	}

	err = replaceFileVar(keyMap, filepath, seed, publicKey, privateKey, prefix)
	if err != nil {
		logrus.Errorf("replaceFileVar error: %s", err)
		return err
	}

	return nil
}

func getRawKeyValue(key, value, seed string) (string, string, error) {
	values := strings.Split(value, ONE_VARIABLE_SEPERATOR)

	if len(values) == 1 {
		return key, values[0], nil
	}
	if len(values) != 2 {
		return key, "", fmt.Errorf("getRawKeyValue key=%v,value=%v is not right formt", key, value)
	}

	//need to decode
	afterDecode, err := AesDePassword(values[1], seed, values[0])
	if err != nil {
		logrus.Errorf("AesDePassword meet error=%v", err)
	}
	return key, afterDecode, err

	// pass code
	//guid := values[1]
	//md5sum := Md5Encode(guid + seed)
	//
	//// judge whether has cipher and remove it
	//var cipher string
	//enCode := values[0]
	//for _, _cipher := range CIPHER_MAP {
	//	if strings.HasPrefix(values[0], _cipher) {
	//		cipher = _cipher
	//		break
	//	}
	//}
	//if cipher != "" {
	//	enCode = enCode[len(cipher):]
	//}
	//
	//data, err := AesDecode(md5sum[0:16], enCode)
	//if err != nil {
	//	logrus.Errorf("AesDecode meet error=%v", err)
	//}
	//return key, data, err
}

func GetInputVariableMap(variable string, seed string) (map[string]string, error) {
	inputMap := make(map[string]string)
	kvs := strings.Split(variable, VARIABLE_KEY_SEPERATOR)
	if len(kvs) != 2 {
		logrus.Errorf("varialbeList(%v) format error", variable)
		return inputMap, fmt.Errorf("varialbeList(%v) format error", variable)
	}

	keys := strings.Split(kvs[0], VARIABLE_VARIABLE_SEPERATOR)
	values := strings.Split(kvs[1], KEY_KEY_SEPERATOR)

	if len(keys) != len(values) {
		logrus.Errorf("varialbeList(%v) format error", variable)
		return inputMap, fmt.Errorf("varialbeList(%v) format error", variable)
	}

	for i, _ := range keys {
		key, value, err := getRawKeyValue(keys[i], values[i], seed)
		if err != nil {
			logrus.Errorf("getRawKeyValue meet error=%v", err)
			return inputMap, err
		}
		inputMap[key] = value
	}
	return inputMap, nil
}

func CheckVariableIsAllReady(input map[string]string, variablelist []string) (err error) {
	for _, va := range variablelist {
		if _, ok := input[va]; !ok {
			return fmt.Errorf("variable %s not input", va)
		}
	}

	return nil
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

func isKeyNeedEncrypt(key string, prefix string) bool {
	return strings.HasPrefix(key, prefix)
}

func encrpytSenstiveData(rawData, publicKey, privateKey string) (string, error) {
	publicKeyFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(publicKeyFile)

	privateKeyFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(privateKeyFile)

	rawDataFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(rawDataFile)

	encrpyDataFile, err := getTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(encrpyDataFile)

	if err = writeStringToFile(publicKey, publicKeyFile); err != nil {
		return "", err
	}

	if err = writeStringToFile(privateKey, privateKeyFile); err != nil {
		return "", err
	}
	if err = writeStringToFile(rawData, rawDataFile); err != nil {
		return "", err
	}

	args := []string{
		"enc",
		publicKeyFile,
		privateKeyFile,
		rawDataFile,
		encrpyDataFile,
	}
	out, err := runBashScript("/home/app/wecube-plugins-saltstack/scripts/rsautil.sh", args)
	fmt.Printf("run script out=%v\n,err=%v\n", out, err)
	if err != nil {
		fmt.Printf("encrpytSenstiveData out=%v,err=%v\n", out, err)
		return "", err
	}

	encryptData, err := readStringFromFile(encrpyDataFile)
	if err != nil {
		return "", err
	}

	return encryptData, nil
}

func getVariableValue(key string, value string, publicKey string, privateKey string, prefix string) (string, error) {
	needEncryt := isKeyNeedEncrypt(key, prefix)
	if !needEncryt {
		return value, nil
	}

	if publicKey == "" {
		return "", errors.New("getVariableValue publicKey is empty")
	}
	if privateKey == "" {
		return "", errors.New("getVariableValue privateKey is empty")
	}

	return encrpytSenstiveData(value, publicKey, privateKey)
}

func replaceFileVar(keyMap map[string]string, filepath, seed, publicKey, privateKey, prefix string) error {
	bf, err := os.Open(filepath)
	if err != nil {
		logrus.Errorf("open file fail: %s", err)
		return err
	}
	defer bf.Close()

	// get file info
	fileInfo, err := bf.Stat()
	if err != nil {
		return err
	}

	// get file mode
	fileMode := fileInfo.Mode()

	newfilePath := filepath + ".bak"
	// f, err := os.OpenFile(newfilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	f, err := os.OpenFile(newfilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
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
		flysnowRegexp := regexp.MustCompile(`[^\[]*]`)
		keys := flysnowRegexp.FindAllString(string(line), -1)
		if len(keys) > 0 {
			for _, key := range keys {
				if false == strings.HasSuffix(key, "]") {
					continue
				}
				key = key[0 : len(key)-1]

				for _, specialFlag := range DefaultSpecialReplaceList {
					if specialFlag == "" {
						continue
					}
					if strings.Contains(key, specialFlag) {
						s := strings.Split(key, specialFlag)
						if s[1] == "" {
							return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
						}
						oldStr := "[" + key + "]"
						variableValue, err := getVariableValue(key, keyMap[s[1]], publicKey, privateKey, prefix)
						if err != nil {
							return err
						}
						newLine = strings.Replace(newLine, oldStr, variableValue, -1)
					}
				}

				//if strings.Contains(key, "@") {
				//	s := strings.Split(key, "@")
				//	if s[1] == "" {
				//		return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
				//	}
				//	oldStr := "[" + key + "]"
				//	variableValue, err := getVariableValue(key, keyMap[s[1]], publicKey, privateKey, prefix)
				//	if err != nil {
				//		return err
				//	}
				//	newLine = strings.Replace(newLine, oldStr, variableValue, -1)
				//}
				//if strings.Contains(key, "!") {
				//	s := strings.Split(key, "!")
				//	if s[1] == "" {
				//		return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
				//	}
				//	oldStr := "[" + key + "]"
				//	variableValue, err := getVariableValue(key, keyMap[s[1]], publicKey, privateKey, prefix)
				//	if err != nil {
				//		return err
				//	}
				//	newLine = strings.Replace(newLine, oldStr, variableValue, -1)
				//}
				//if strings.Contains(key, "&") {
				//	s := strings.Split(key, "&")
				//	if s[1] == "" {
				//		return fmt.Errorf("file %s have unvaliable variable %s", filepath, key)
				//	}
				//	oldStr := "[" + key + "]"
				//	variableValue, err := getVariableValue(key, keyMap[s[1]], publicKey, privateKey, prefix)
				//	if err != nil {
				//		return err
				//	}
				//	newLine = strings.Replace(newLine, oldStr, variableValue, -1)
				//}
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
		sh = "cd " + decompressDirName + " && " + "zip -r " + UPLOADS3FILE_DIR + newPackageName + " * .[^.]*"
	}
	if suffix == ".tgz" || suffix == ".tar.gz" {
		sh = "cd " + decompressDirName + " && " + "tar czf  " + UPLOADS3FILE_DIR + newPackageName + " * .[^.]*"
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
