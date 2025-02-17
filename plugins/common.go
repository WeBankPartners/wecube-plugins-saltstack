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
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
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
	DefaultS3Key                    = "access_key"
	DefaultS3Password               = "secret_key"
	DefaultSpecialReplaceList       []string
	DefaultEncryptReplaceList       []string
	DefaultSingleEncryptReplaceList []string
	DefaultEncryptEscapeList        []string
	DefaultFileReplaceList          []string
	ClusterList                     []string
	MasterHostIp                    string
	CoreUrl                         string
	DefaultS3TmpAddress             string
	SubSystemCode                   string
	SubSystemKey                    string
	SaltResetEnv                    bool
	ApiConcurrentNum                int
	VariableNullCheck               bool
	GlobalEncryptSeed               string
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
	log.Logger.Debug("Call salt api request", log.String("serviceUrl", serviceUrl), log.JsonObj("param", request))

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

	bytesJson, jsonMarshalErr := json.Marshal(request)
	if jsonMarshalErr != nil {
		err = fmt.Errorf("json marshal request data fail,%s ", jsonMarshalErr.Error())
		return "", err
	}

	req, newReqErr := http.NewRequest("POST", serviceUrl, bytes.NewBuffer(bytesJson))
	if newReqErr != nil {
		return "", fmt.Errorf("new salt api request meet error = %s", newReqErr.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, doHttpReqErr := client.Do(req)
	if doHttpReqErr != nil {
		return "", fmt.Errorf("call salt api server meet error = %s", doHttpReqErr.Error())
	}

	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	log.Logger.Debug("Call salt api response", log.String("body", result))

	apiResult := callSaltApiResults{}
	if err = json.Unmarshal([]byte(result), &apiResult); err != nil {
		if strings.Contains(result, "TimeoutError") {
			err = fmt.Errorf("Call Timeout!! ")
			return "", err
		}
		log.Logger.Error("Call salt api unmarshal result error", log.Error(err), log.String("body", result))
		return "", err
	}

	if len(apiResult.Results) == 0 || len(apiResult.Results[0]) == 0 {
		return "", getSaltApiTargetError(language, request.Target)
	}
	for _, resultObj := range apiResult.Results {
		for k, v := range resultObj {
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
	if seed == "" {
		return password, nil
	}
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			return password, nil
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
	var cipherPrefix string
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			cipherPrefix = _cipher
			break
		}
	}
	if cipherPrefix == "" {
		return password, nil
	}
	password = password[len(cipherPrefix):]

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
		if lastIndex := strings.LastIndex(filePath, "/"); lastIndex >= 0 {
			unpackToDirPath = filePath[0:lastIndex]
		} else {
			return fmt.Errorf("filePath:%s illegal with absolute path check ", filePath)
		}
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
	} else {
		log.Logger.Warn("Variable encrypt replace without any param")
	}
	tmpSingleEncryptReplace := os.Getenv("SALTSTACK_SINGLE_ENCRYPT_VARIBLE_PREFIX")
	if tmpEncryptReplace != "" {
		DefaultSingleEncryptReplaceList = strings.Split(tmpSingleEncryptReplace, ",")
		log.Logger.Info("Variable single encrypt", log.StringList("special", DefaultSingleEncryptReplaceList))
	} else {
		log.Logger.Warn("Variable single encrypt replace without any param")
	}
	tmpEncryptEscape := os.Getenv("SALTSTACK_ENCRYPT_ESCAPE_PREFIX")
	if tmpEncryptReplace != "" {
		DefaultEncryptEscapeList = strings.Split(tmpEncryptEscape, ",")
		log.Logger.Info("Variable encrypt escape", log.StringList("special", DefaultEncryptEscapeList))
	} else {
		log.Logger.Warn("Variable encrypt escape without any param")
	}
	tmpFileReplace := os.Getenv("SALTSTACK_FILE_VARIBLE_PREFIX")
	if tmpFileReplace != "" {
		DefaultFileReplaceList = strings.Split(tmpFileReplace, ",")
		log.Logger.Info("Variable file", log.StringList("special", DefaultFileReplaceList))
	} else {
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
	} else {
		log.Logger.Warn("Default s3 address not found")
	}
	SubSystemCode = os.Getenv("SUB_SYSTEM_CODE")
	SubSystemKey = os.Getenv("SUB_SYSTEM_KEY")
	if SubSystemCode == "" {
		log.Logger.Warn("Env SUB_SYSTEM_CODE is empty")
	}
	if SubSystemKey == "" {
		log.Logger.Warn("Env SUB_SYSTEM_KEY is empty")
	}
	saltResetEnvString := strings.ToLower(os.Getenv("SALTSTACK_RESET_ENV"))
	if saltResetEnvString == "y" || saltResetEnvString == "yes" || saltResetEnvString == "true" || saltResetEnvString == "" {
		SaltResetEnv = true
	} else {
		SaltResetEnv = false
	}
	ApiConcurrentNum, _ = strconv.Atoi(os.Getenv("API_CONCURRENT_NUM"))
	if ApiConcurrentNum <= 0 {
		ApiConcurrentNum = 5
	}
	log.Logger.Info("API_CONCURRENT_NUM", log.Int("num", ApiConcurrentNum))
	variableNullCheckEnvString := strings.ToLower(os.Getenv("SALTSTACK_VARIABLE_NULL_CHECK"))
	if variableNullCheckEnvString == "y" || variableNullCheckEnvString == "yes" || variableNullCheckEnvString == "true" {
		VariableNullCheck = true
	} else {
		VariableNullCheck = false
	}
	GlobalEncryptSeed = os.Getenv("ENCRYPT_SEED")
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
	inputBytes, err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		log.Logger.Error("Input string format to base64 fail", log.Error(err))
		return inputString
	}
	pemPath := "/data/certs/rsa_key"
	fileContent, err := ioutil.ReadFile(pemPath)
	if err != nil {
		log.Logger.Error("Read file fail", log.String("path", pemPath), log.Error(err))
		return result
	}
	block, _ := pem.Decode(fileContent)
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Logger.Error("Parse private key fail", log.Error(err))
		return result
	}
	privateKey := privateKeyInterface.(*rsa.PrivateKey)
	decodeBytes, err := rsa.DecryptPKCS1v15(crand.Reader, privateKey, inputBytes)
	if err != nil {
		log.Logger.Error("Decode fail", log.Error(err))
		return result
	}
	result = string(decodeBytes)
	return result
}

func RSAEncryptByPrivate(orgidata []byte, privatekey string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(privatekey)
	if err != nil {
		return nil, fmt.Errorf("RSASign private key is bad")
	}

	privInterface, err := x509.ParsePKCS8PrivateKey(decodeBytes)
	if err != nil {
		return nil, err
	}

	priv := privInterface.(*rsa.PrivateKey)

	k := (priv.N.BitLen() + 7) / 8
	tLen := len(orgidata)
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], orgidata)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(priv.N) > 0 {
		return nil, nil
	}
	var m *big.Int
	var ir *big.Int
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}
	return m.Bytes(), nil
}

func getEncryptSeed(inputSeed string) string {
	if GlobalEncryptSeed == "" {
		return inputSeed
	}
	if inputSeed == "" {
		return GlobalEncryptSeed
	}
	var cipherPrefix string
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(inputSeed, _cipher) {
			cipherPrefix = _cipher
			break
		}
	}
	if cipherPrefix == "" {
		return inputSeed
	}
	return GlobalEncryptSeed
}

func isContains(sList []string, t string) bool {
	for _, s := range sList {
		if s == t {
			return true
		}
	}
	return false
}
