package plugins

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"encoding/base64"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	crand "crypto/rand"
)

const (
	CHARGE_TYPE_PREPAID = "PREPAID"
	RESULT_CODE_SUCCESS = "0"
	RESULT_CODE_ERROR   = "1"
	PASSWORD_LEN        = 12
	DEFALT_CIPHER       = "CIPHER_A"
	ASCII_CODE_LF       = 10
	SystemRole          = `SUB_SYSTEM`
	PlatformUser        = `SYS_PLATFORM`
)

var (
	DefaultS3Key              = "access_key"
	DefaultS3Password         = "secret_key"
	DefaultSpecialReplaceList []string
	DefaultEncryptReplaceList []string
	DefaultFileReplaceList    []string
	ClusterList               []string
	MasterHostIp              string
	CoreUrl                   string
	DefaultS3TmpAddress       string
)

var CIPHER_MAP = map[string]string{
	"CIPHER_A": "{cipher_a}",
}

type CallBackParameter struct {
	Parameter string `json:"callbackParameter,omitempty"`
}

type Result struct {
	Code    string `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func replaceLF(str string) string {
	if str == "" {
		return str
	}
	buf := []byte(str)
	ns := bytes.Replace(buf, []byte{92, 110}, []byte{ASCII_CODE_LF}, -1)

	return string(ns)
}

func Md5Encode(rawData string) string {
	data := []byte(rawData)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}

func AesEncode(key string, rawData string) (string, error) {
	bytesRawKey := []byte(key)
	block, err := aes.NewCipher(bytesRawKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding([]byte(rawData), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, bytesRawKey[:blockSize])
	crypted := make([]byte, len([]byte(origData)))
	blockMode.CryptBlocks(crypted, origData)
	return hex.EncodeToString(crypted), nil
}

func AesDecode(key string, encryptData string) (password string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	bytesRawKey := []byte(key)
	bytesRawData, _ := hex.DecodeString(encryptData)
	block, err := aes.NewCipher(bytesRawKey)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, bytesRawKey[:blockSize])
	origData := make([]byte, len(bytesRawData))
	blockMode.CryptBlocks(origData, bytesRawData)
	origData = PKCS7UnPadding(origData)
	if len(origData) == 0 {
		err = fmt.Errorf("password wrong")
		return
	}

	password = string(origData)
	return
}

func UnmarshalJson(source interface{}, target interface{}) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("the source to be unmarshaled is not a io.reader type")
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("parse http request (%v) meet error (%v)", reader, err)
	}

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("unmarshal http request (%v) meet error (%v)", reader, err)
	}
	return nil
}

func ExtractJsonFromStruct(s interface{}) map[string]string {
	fields := make(map[string]string)
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i).Tag.Get("json")
			fields[strings.Split(field, ",")[0]] = t.Field(i).Type.String()
		}
	}
	return fields
}

type SaltApiRequest struct {
	Client     string   `json:"client,omitempty"`
	TargetType string   `json:"expr_form,omitempty"`
	Target     string   `json:"tgt,omitempty"`
	Function   string   `json:"fun,omitempty"`
	Args       []string `json:"arg,omitempty"`
	FullReturn bool     `json:"full_return,omitempty"`
}

type callSaltApiResults struct {
	Results []map[string]interface{} `json:"return,omitempty"`
}

func CallSaltApi(serviceUrl string, request SaltApiRequest, language string) (string, error) {
	log.Logger.Debug("Call salt api request", log.JsonObj("param", request))

	token, err := getSaltApiToken()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	bytesJson, err := json.Marshal(request)

	req, err := http.NewRequest("POST", serviceUrl, bytes.NewBuffer(bytesJson))
	if err != nil {
		return "", fmt.Errorf("new salt api request meet error = %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call salt api server meet error = %v", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	log.Logger.Debug("Call salt api response", log.String("body", result))

	apiResult := callSaltApiResults{}
	if err := json.Unmarshal([]byte(result), &apiResult); err != nil {
		log.Logger.Error("Call salt api unmarshal result error", log.Error(err))
		return "", err
	}

	if len(apiResult.Results) == 0 || len(apiResult.Results[0]) == 0 {
		return "", getSaltApiTargetError(language, request.Target)
	}
	for _, result := range apiResult.Results {
		for k, v := range result {
			switch v.(type) {
			case bool:
				if v.(bool) == false {
					return "", getSaltApiConnectError(language, k)
				}
			}
		}
	}

	return result, nil
}

func createRandomPassword() string {
	digitals := "0123456789"
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(letters)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < PASSWORD_LEN-4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	bytes = []byte(digitals)
	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

func AesEnPassword(guid, seed, password, cipher string) (string, error) {
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			return password,nil
		}
	}
	if cipher == "" {
		cipher = DEFALT_CIPHER
	}
	md5sum := Md5Encode(guid + seed)
	enPassword, err := AesEncode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return CIPHER_MAP[cipher] + enPassword, nil
}

func AesDePassword(guid, seed, password string) (string, error) {
	var cipher string
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			cipher = _cipher
			break
		}
	}
	if cipher == "" {
		return password, nil
	}
	password = password[len(cipher):]

	md5sum := Md5Encode(guid + seed)
	dePassword, err := AesDecode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return dePassword, nil
}

func getTempFile() (string, error) {
	file, err := ioutil.TempFile("/tmp/", "qcloud_key")
	if err != nil {
		return "", err
	}
	file.Close()
	return file.Name(), nil
}

func writeStringToFile(data string, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(data)
	return err
}

func readStringFromFile(fileName string) (string, error) {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	if f[len(f)-1] == ASCII_CODE_LF {
		f = f[:len(f)-1]
	}
	return string(f), nil
}

func fileExist(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	return !info.IsDir()
}

func listFile(myDir string) ([]string, error) {
	output := []string{}
	files, err := ioutil.ReadDir(myDir)
	if err != nil {
		return output, err
	}
	for _, file := range files {
		if file.IsDir() {
			childOutput, err := listFile(myDir + "/" + file.Name())
			if err != nil {
				return output, err
			}
			output = append(output, childOutput...)
		} else {
			output = append(output, myDir+"/"+file.Name())
		}
	}
	return output, err
}

func deriveUnpackfile(filePath string, desDirPath string, overwrite bool, language string) error {
	name := ""
	args := []string{}
	lowerFilepath := strings.ToLower(filePath)
	unpackToDirPath := ""
	if desDirPath == "" {
		unpackToDirPath = filePath[0:strings.LastIndex(filePath, "/")]
	} else {
		unpackToDirPath = desDirPath
	}

	if strings.HasSuffix(lowerFilepath, ".zip") {
		name = "unzip"
		if overwrite {
			args = append(args, "-o")
		}
		args = append(args, filePath, "-d", unpackToDirPath)
	} else if strings.HasSuffix(lowerFilepath, ".rar") {
		name = "unrar"
		if desDirPath == "" {
			args = append(args, "e", filePath)
		} else {
			args = append(args, "x", filePath, unpackToDirPath)
		}
	} else if strings.HasSuffix(lowerFilepath, ".tar") {
		name = "tar"
		args = append(args, "xf", filePath, "-C", unpackToDirPath)
	} else if strings.HasSuffix(lowerFilepath, ".tar.gz") || strings.HasSuffix(lowerFilepath, ".tgz") {
		name = "tar"
		args = append(args, "zxf", filePath, "-C", unpackToDirPath)
	} else {
		return getDecompressSuffixError(language, lowerFilepath)
	}

	command := exec.Command(name, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		log.Logger.Error("Unpack file", log.String("name", name), log.StringList("args", args), log.String("output", string(out)), log.Error(err))
		return getUnpackFileError(language, lowerFilepath, err)
	}
	return nil
}

func InitEnvParam() {
	tmpKey := os.Getenv("DEFAULT_S3_KEY")
	if tmpKey != "" {
		DefaultS3Key = DecryptRsa(tmpKey)
	}
	tmpPwd := os.Getenv("DEFAULT_S3_PASSWORD")
	if tmpPwd != "" {
		DefaultS3Password = DecryptRsa(tmpPwd)
	}
	log.Logger.Info("S3 config", log.String("key", DefaultS3Key))
	log.Logger.Debug("S3 config", log.String("password", DefaultS3Password))
	tmpSpecialReplace := os.Getenv("SALTSTACK_DEFAULT_SPECIAL_REPLACE")
	if tmpSpecialReplace != "" {
		DefaultSpecialReplaceList = strings.Split(tmpSpecialReplace, ",")
		log.Logger.Info("Variable replace config", log.StringList("special", DefaultSpecialReplaceList))
	} else {
		log.Logger.Warn("Variable replace without any param")
	}
	tmpEncryptReplace := os.Getenv("SALTSTACK_ENCRYPT_VARIBLE_PREFIX")
	if tmpEncryptReplace != "" {
		DefaultEncryptReplaceList = strings.Split(tmpEncryptReplace, ",")
		log.Logger.Info("Variable encrypt", log.StringList("special", DefaultEncryptReplaceList))
	}else{
		log.Logger.Warn("Variable encrypt replace without any param")
	}
	tmpFileReplace := os.Getenv("SALTSTACK_FILE_VARIBLE_PREFIX")
	if tmpFileReplace != "" {
		DefaultFileReplaceList = strings.Split(tmpFileReplace, ",")
		log.Logger.Info("Variable file", log.StringList("special", DefaultFileReplaceList))
	}else{
		log.Logger.Warn("Variable file replace without any param")
	}
	tmpHostIp := os.Getenv("minion_master_ip")
	if tmpHostIp != "" {
		MasterHostIp = tmpHostIp
		log.Logger.Info("Master host", log.String("ip", MasterHostIp))
	} else {
		log.Logger.Warn("Master host ip not found,default null")
	}
	tmpCoreUrl := os.Getenv("GATEWAY_URL")
	if tmpCoreUrl != "" {
		CoreUrl = tmpCoreUrl
		log.Logger.Info("Core url", log.String("url", CoreUrl))
	} else {
		log.Logger.Warn("Core url is empty")
	}
	tmpS3Address := os.Getenv("S3_SERVER_URL")
	if tmpS3Address != "" {
		if strings.HasSuffix(tmpS3Address, "/") {
			tmpS3Address = tmpS3Address[:len(tmpS3Address)-1]
		}
		if !strings.Contains(tmpS3Address, "salt-tmp") {
			tmpS3Address = tmpS3Address + "/salt-tmp"
		}
		DefaultS3TmpAddress = tmpS3Address
		log.Logger.Info("Default s3 address", log.String("address", DefaultS3TmpAddress))
	}else{
		log.Logger.Warn("Default s3 address not found")
	}
}

func checkIllegalParam(input string) bool {
	if strings.Contains(input, "'") {
		return true
	}
	if strings.Contains(input, "\"") {
		return true
	}
	if strings.Contains(input, "`") {
		return true
	}
	return false
}

func DecryptRsa(inputString string) string {
	if !strings.HasPrefix(strings.ToLower(inputString), "rsa@") {
		return inputString
	}
	inputString = inputString[4:]
	result := inputString
	inputBytes,err := base64.RawStdEncoding.DecodeString(inputString)
	if err != nil {
		log.Logger.Error("Input string format to base64 fail", log.Error(err))
		return inputString
	}
	pemPath := "/data/certs/rsa_key"
	fileContent,err := ioutil.ReadFile(pemPath)
	if err != nil {
		log.Logger.Error("Read file fail", log.String("path", pemPath), log.Error(err))
		return result
	}
	block,_ := pem.Decode(fileContent)
	privateKeyInterface,err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Logger.Error("Parse private key fail", log.Error(err))
		return result
	}
	privateKey := privateKeyInterface.(*rsa.PrivateKey)
	decodeBytes,err := rsa.DecryptPKCS1v15(crand.Reader, privateKey, inputBytes)
	if err != nil {
		log.Logger.Error("Decode fail", log.Error(err))
		return result
	}
	result = string(decodeBytes)
	return result
}
