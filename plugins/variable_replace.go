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

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
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
	FileReplacePrefix    string `json:"fileReplacePrefix,omitempty"`
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
	Language string
}

func (action *VariableReplaceAction) SetAcceptLanguage(language string) {
	action.Language = language
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
		return getParamEmptyError(action.Language, "endpoint")
	}
	if input.VariableList != "" {
		if !strings.Contains(input.VariableList, "=") {
			return getParamValidateError(action.Language, "variableList", "can not find '=' in the content,variable should be k=v")
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
		err = getDecompressSuffixError(action.Language, input.EndPoint)
		return output, err
	}

	packageName, err := getPackageNameFromEndpoint(input.EndPoint)
	if err != nil {
		return output, err
	}

	decompressDirName := getDecompressDirName(packageName)
	if err = isDirExist(decompressDirName); err == nil {
		os.RemoveAll(decompressDirName)
	}

	if err = os.MkdirAll(decompressDirName, os.ModePerm); err != nil {
		return output, err
	}

	compressedFileFullPath, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false, action.Language)
	if err != nil {
		return output, err
	}

	if err = decompressFile(compressedFileFullPath, decompressDirName); err != nil {
		err = getUnpackFileError(action.Language, compressedFileFullPath, err)
		os.RemoveAll(compressedFileFullPath)
		return output, err
	}
	os.RemoveAll(compressedFileFullPath)

	if input.FilePath != "" && input.VariableList != "" {
		for _, filePath := range splitWithCustomFlag(input.FilePath) {
			confFilePath := ""

			if decompressDirName[len(decompressDirName)-1] == '/' {
				decompressDirName = decompressDirName[:len(decompressDirName)-1]
			}

			if filePath[0] == '/' {
				confFilePath = decompressDirName + filePath
			} else {
				confFilePath = decompressDirName + "/" + filePath
			}


			if err = ReplaceFileVar(confFilePath, input, decompressDirName); err != nil {
				os.RemoveAll(decompressDirName)
				return output, err
			}
		}
	}

	//compress file
	//nowTime := time.Now().Format("20060102150405.999999999")
	newPackageName := fmt.Sprintf("%s-%s%s", getPackageNameWithoutSuffix(packageName), time.Now().Format("20060102150405.999999999"), suffix)
	err,newPackageName = compressDir(decompressDirName, suffix, newPackageName)
	if err != nil {
		os.RemoveAll(decompressDirName)
		return output, fmt.Errorf("After replace variable,try to compress %s fail,%s ", newPackageName, err.Error())
	}
	os.RemoveAll(decompressDirName)

	//upload to s3
	newS3Endpoint := getNewS3EndpointName(input.EndPoint, newPackageName)
	log.Logger.Info("Upload new file to s3", log.String("file", newS3Endpoint))

	if _, err = uploadS3File(newS3Endpoint, DefaultS3Key, DefaultS3Password, action.Language); err != nil {
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
func ReplaceFileVar(filepath string, input *VariableReplaceInput, decompressDirName string) error {
	variablelist := input.VariableList
	seed := input.Seed
	publicKey := input.AppPublicKey
	privateKey := input.SysPrivateKey
	prefix := DefaultEncryptReplaceList
	fileReplacePrefix := DefaultFileReplaceList

	index := strings.LastIndexAny(filepath, "/")
	if index == -1 {
		return fmt.Errorf("Invalid endpoint %s ", filepath)
	}
	tmpSpecialReplaceList := DefaultSpecialReplaceList
	tmpSpecialReplaceList = append(tmpSpecialReplaceList, prefix...)
	tmpSpecialReplaceList = append(tmpSpecialReplaceList, fileReplacePrefix...)
	if !checkIsUniqueList(tmpSpecialReplaceList) {
		return fmt.Errorf("Prefix duplicate ,defaultPrefix:%v encryptPrefix:%v fileReplacePrefix:%v ", DefaultSpecialReplaceList, prefix, fileReplacePrefix)
	}
	fileName := filepath[index+1:]
	fileVarMap, err := GetVariable(filepath, tmpSpecialReplaceList, false)
	if err != nil {
		return err
	}

	if len(fileVarMap) == 0 {
		log.Logger.Warn("Replace variable key,no variable need to replace", log.String("file", fileName))
		return nil
	}

	//fileVarList := []string{}
	//for _, v := range fileVarMap {
	//	fileVarList = append(fileVarList, v.Key)
	//}

	keyMap, err := GetInputVariableMap(variablelist, seed, tmpSpecialReplaceList)
	if err != nil {
		return err
	}

	//err = CheckVariableIsAllReady(keyMap, fileVarList)
	//if err != nil {
	//	return err
	//}

	err = replaceFileVar(keyMap, filepath, seed, publicKey, privateKey, decompressDirName, tmpSpecialReplaceList, prefix, fileReplacePrefix)
	if err != nil {
		return err
	}

	return nil
}

func getRawKeyValue(key, value, seed string) (string, string, error) {
	values := strings.Split(value, ONE_VARIABLE_SEPERATOR)

	if len(values) == 1 {
		return key, values[0], nil
	}
	if len(values) != 2 || values[1] == "" {
		return key, "", fmt.Errorf("GetRawKeyValue key=%v,value=%v format error,encrypt value should contain guid message", key, value)
	}

	//need to decode
	afterDecode, err := AesDePassword(values[1], seed, values[0])
	if err != nil {
		log.Logger.Error("GetRawKey fail,decode password error", log.Error(err))
	}
	return key, afterDecode, err
}

func GetInputVariableMap(variable string, seed string, specialList []string) (map[string]string, error) {
	inputMap := make(map[string]string)
	kvs := strings.Split(variable, VARIABLE_KEY_SEPERATOR)
	if len(kvs) != 2 {
		return inputMap, fmt.Errorf("VarialbeList(%v) format error,can not find '='", variable)
	}

	keys := strings.Split(kvs[0], VARIABLE_VARIABLE_SEPERATOR)
	values := strings.Split(kvs[1], KEY_KEY_SEPERATOR)

	if len(keys) != len(values) {
		return inputMap, fmt.Errorf("VarialbeList(%v) format error,keys num != value num", variable)
	}

	for i, _ := range keys {
		key, value, err := getRawKeyValue(keys[i], values[i], seed)
		if err != nil {
			return inputMap, err
		}
		for _,v := range specialList {
			if strings.HasPrefix(key, v) {
				key = key[len(v):]
				break
			}
		}
		key = strings.ToLower(key)
		inputMap[key] = value
	}
	return inputMap, nil
}

func CheckVariableIsAllReady(input map[string]string, variablelist []string) (err error) {
	for _, va := range variablelist {
		toLowerV := strings.ToLower(va)
		if _, ok := input[toLowerV]; !ok {
			return fmt.Errorf("Variable %s prefix fetch,but not input ", va)
		}
	}

	return nil
}

func PathExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		if f.IsDir() {
			return false,fmt.Errorf("path:%s is dir", path)
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("path:%s is not exist", path)
	}
	return false, err
}

func isKeyNeedEncrypt(key string, prefix []string) bool {
	isNeed := false
	for _,v := range prefix {
		if v == "" {
			continue
		}
		if strings.HasPrefix(key, v) {
			isNeed = true
			break
		}
	}
	return isNeed
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
	tmpWorkspace := fmt.Sprintf("enc-%d", time.Now().UnixNano())
	args := []string{
		"enc",
		publicKeyFile,
		privateKeyFile,
		rawDataFile,
		encrpyDataFile,
		tmpWorkspace,
	}
	out, err := runBashScript("/home/app/wecube-plugins-saltstack/scripts/rsautil.sh", args)
	if err != nil {
		log.Logger.Error("Encrypt variable data fail", log.String("output", out), log.Error(err))
		return "", fmt.Errorf("Run encrypt bash shell fail,output:%s, err:%s ", out, err.Error())
	}

	encryptData, err := readStringFromFile(encrpyDataFile)
	if err != nil {
		return "", err
	}

	return encryptData, nil
}

func getVariableValue(key string, value string, publicKey string, privateKey string, prefix []string) (string, error) {
	needEncrypt := isKeyNeedEncrypt(key, prefix)
	if !needEncrypt {
		return value, nil
	}

	if publicKey == "" {
		return "", errors.New("GetVariableValue publicKey is empty")
	}
	if privateKey == "" {
		return "", errors.New("GetVariableValue privateKey is empty")
	}

	publicKey = replaceLF(publicKey)
	privateKey = replaceLF(privateKey)

	encryptValue,err := encrpytSenstiveData(value, publicKey, privateKey)
	if err != nil {
		err = fmt.Errorf("Try to encrypt key %s fail,%s ", key, err.Error())
	}
	return encryptValue,err
}

func replaceFileVar(keyMap map[string]string, filepath, seed, publicKey, privateKey, decompressDirName string,specialReplaceList, prefix, fileReplacePrefix []string) error {
	bf, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Open file %s fail,%s ", filepath, err.Error())
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
		return fmt.Errorf("Try to create tmp file fail,%s ", err.Error())
	}
	defer f.Close()
	fileReplaceMap := make(map[string]string)
	tmpLineCount := 0
	br := bufio.NewReader(bf)
	for {
		tmpLineCount = tmpLineCount + 1
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Read file %s line %d error:%s ", filepath, tmpLineCount, err.Error())
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

				for _, specialFlag := range specialReplaceList {
					if specialFlag == "" {
						continue
					}
					if strings.HasPrefix(key, specialFlag) {
						s := strings.Split(key, specialFlag)
						if s[1] == "" {
							return fmt.Errorf("File %s have empty variable [%s] in line %d ", filepath, key, tmpLineCount)
						}
						if strings.Contains(s[1], " ") {
							continue
						}
						toLowerKey := strings.ToLower(s[1])
						if _,b := keyMap[toLowerKey]; !b {
							continue
						}
						oldStr := "[" + key + "]"
						variableValue, err := getVariableValue(key, keyMap[toLowerKey], publicKey, privateKey, prefix)
						if err != nil {
							return err
						}
						isFileReplaceFlag := false
						for _,frPrefix := range fileReplacePrefix {
							if specialFlag == frPrefix {
								isFileReplaceFlag = true
								break
							}
						}
						if isFileReplaceFlag {
							fileReplaceMap[key[len(specialFlag):]] = variableValue
						}else {
							newLine = strings.Replace(newLine, oldStr, variableValue, -1)
						}
					}
				}
			}
		}
		_, err = f.WriteString(newLine + "\n")
		if err != nil {
			return fmt.Errorf("Try to write new line to tmp file fail,%s ", err.Error())
		}
	}
	err = os.Rename(newfilePath, filepath)
	if err != nil {
		return fmt.Errorf("Replace file %s with tmp file fail,%s ", filepath, err.Error())
	}
	if len(fileReplaceMap) > 0 {
		var tmpOut []byte
		for k,v := range fileReplaceMap {
			if k == "" || v == "" {
				continue
			}
			sourceFile, err := downloadS3File(v, DefaultS3Key, DefaultS3Password, false, "")
			if err != nil {
				log.Logger.Error(fmt.Sprintf("VariableReplaceAction downloadS3File get replace source file s3Path=%s,fullPath=%s,err=%v", v, sourceFile, err))
				return err
			}
			moveCmd := fmt.Sprintf("mv -f %s %s/%s", sourceFile, decompressDirName, k)
			if k[:1] == "/" {
				moveCmd = fmt.Sprintf("mv -f %s %s%s", sourceFile, decompressDirName, k)
			}
			log.Logger.Debug(fmt.Sprintf("File replace ,source: %s -> dist: %s command: %s \n", sourceFile, k, moveCmd))
			tmpOut,err = exec.Command("/bin/bash", "-c", moveCmd).Output()
			if err != nil {
				log.Logger.Error("File replace fail", log.String("command", moveCmd), log.String("output", string(tmpOut)), log.Error(err))
				return fmt.Errorf("Try to replace file fail,output:%s,error:%s ", string(tmpOut), err.Error())
			}
		}
	}

	return nil
}

func compressDir(decompressDirName string, suffix string, newPackageName string) (error,string) {
	sh := ""
	if suffix != ".zip" && suffix != ".tgz" && suffix != ".tar.gz" {
		return fmt.Errorf("%s is invalid suffix", suffix),newPackageName
	}

	if suffix == ".zip" {
		sh = "cd " + decompressDirName + " && " + "zip -r " + UPLOADS3FILE_DIR + newPackageName + " * .[^.]*"
	}
	if suffix == ".tgz" || suffix == ".tar.gz" {
		sh = "cd " + decompressDirName + " && " + "tar czf  " + UPLOADS3FILE_DIR + newPackageName + " * .[^.]*"
	}

	cmd := exec.Command("/bin/sh", "-c", sh)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger.Error("Can not obtain stdout pipe", log.String("command", sh), log.Error(err))
		return err,newPackageName
	}
	if err := cmd.Start(); err != nil {
		log.Logger.Error("Command start error", log.Error(err))
		return err,newPackageName
	}
	_, err = LogReadLine(cmd, stdout)
	if err != nil {
		return err,newPackageName
	}
	newPackagePath := UPLOADS3FILE_DIR + newPackageName
	md5Value,err := GetFileMD5Value(newPackagePath)
	if err != nil {
		return err,newPackageName
	}
	if strings.Contains(newPackageName, "_") {
		tmpOldMd5Value := strings.Split(newPackageName, "_")[0]
		if len(tmpOldMd5Value) == 32 {
			newPackageName = newPackageName[33:]
		}
	}
	newPackageName = fmt.Sprintf("%s_%s", md5Value, newPackageName)
	output,err := exec.Command("/bin/bash", "-c", fmt.Sprintf("'mv %s %s%s'", newPackagePath, UPLOADS3FILE_DIR, newPackageName)).Output()
	if err != nil {
		return fmt.Errorf("Try to rename package name fail,output=%s,error=%s ", string(output), err.Error()), newPackageName
	}

	return nil,newPackageName
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
				log.Logger.Warn("Read line error", log.Error(err))
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

func GetFileMD5Value(filePath string) (string, error) {
	output,err := exec.Command("/bin/bash", "-c", fmt.Sprintf("'md5sum %s'", filePath)).Output()
	if err != nil {
		log.Logger.Error("Get md5 value fail", log.String("file", filePath), log.Error(err))
		return "",fmt.Errorf("Try to get md5 value fail,output=%s,error=%s ", string(output), err.Error())
	}
	outputSplit := strings.Split(string(output), " ")
	return outputSplit[0], nil
}

func checkIsUniqueList(aList []string) bool {
	tmpMap := make(map[string]int)
	for _,v := range aList {
		if _,b:=tmpMap[v];b {
			return false
		}
		tmpMap[v] = 1
	}
	return true
}