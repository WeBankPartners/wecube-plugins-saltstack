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

	"github.com/sirupsen/logrus"
)

const (
	CHARGE_TYPE_PREPAID = "PREPAID"
	RESULT_CODE_SUCCESS = "0"
	RESULT_CODE_ERROR   = "1"
	PASSWORD_LEN        = 12
	DEFALT_CIPHER       = "CIPHER_A"
)

var (
	DefaultS3Key = "access_key"
	DefaultS3Password = "secret_key"
	DefaultSpecialReplaceList []string
	ClusterList []string
	MasterHostIp  string
	CoreUrl  string
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

func CallSaltApi(serviceUrl string, request SaltApiRequest) (string, error) {
	logrus.Infof("call salt api request = %v", request)

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
	logrus.Infof("call salt api response = %v", result)

	apiResult := callSaltApiResults{}
	if err := json.Unmarshal([]byte(result), &apiResult); err != nil {
		logrus.Infof("callSaltApi unmarshal result meet error=%v ", err)
		return "", err
	}
	logrus.Infof("apiResult: %++v", apiResult)

	if len(apiResult.Results) == 0 || len(apiResult.Results[0]) == 0 {
		return "", fmt.Errorf("salt api:no target match ,please check if salt-agent installed on target,reqeust=%v", request)
	}
	for _, result := range apiResult.Results {
		for k, v := range result {
			switch v.(type) {
			case bool:
				if v.(bool) == false {
					return "", fmt.Errorf("salt api: can not connect to target[%v]", k)
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

func deriveUnpackfile(filePath string, desDirPath string, overwrite bool) error {
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
		return fmt.Errorf("%s has invalid compressed format", lowerFilepath)
	}

	command := exec.Command(name, args...)
	out, err := command.CombinedOutput()
	logrus.Infof("runDatabaseCommand(%v) output=%v,err=%v\n", command, string(out), err)
	if err != nil {
		return err
	}
	return nil
}

func InitEnvParam()  {
	tmpKey := os.Getenv("DEFAULT_S3_KEY")
	if tmpKey != "" {
		DefaultS3Key = tmpKey
	}
	tmpPwd := os.Getenv("DEFAULT_S3_PASSWORD")
	if tmpPwd != "" {
		DefaultS3Password = tmpPwd
	}
	logrus.Infof("s3  --> key: %s  password: %s \n", DefaultS3Key, DefaultS3Password)
	tmpSpecialReplace := os.Getenv("SALTSTACK_DEFAULT_SPECIAL_REPLACE")
	if tmpSpecialReplace != "" {
		logrus.Infof("variable replace  --> special flag: %s  \n", tmpSpecialReplace)
		DefaultSpecialReplaceList = strings.Split(tmpSpecialReplace, ",")
	}else{
		DefaultSpecialReplaceList = []string{"@","#","!","&"}
		logrus.Infof("variable replace without param,use default @,#,!,&  \n")
	}
	tmpHostIp := os.Getenv("minion_master_ip")
	if tmpHostIp != "" {
		logrus.Infof("master host ip: %s  \n", tmpHostIp)
		MasterHostIp = tmpHostIp
	}else{
		logrus.Infof("master host ip not found,default null!!  \n")
	}
	tmpCoreUrl := os.Getenv("CORE_ADDR")
	if tmpCoreUrl == "" {
		tmpCoreUrl = os.Getenv("GATEWAY_URL")
	}
	if tmpCoreUrl != "" {
		logrus.Infof("core url : %s  \n", tmpCoreUrl)
		CoreUrl = tmpCoreUrl
	}else{
		logrus.Infof("core url is empty!!  \n")
	}
}

const TmpCoreToken = `Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ3ZHNfc3lzdGVtIiwiaWF0IjoxNTgyMDkyODYyLCJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJjbGllbnRUeXBlIjoiVVNFUiIsImV4cCI6MzQ3NDI1Mjg2MiwiYXV0aG9yaXR5IjoiW2FkbWluLElNUExFTUVOVEFUSU9OX1dPUktGTE9XX0VYRUNVVElPTixJTVBMRU1FTlRBVElPTl9CQVRDSF9FWEVDVVRJT04sSU1QTEVNRU5UQVRJT05fQVJUSUZBQ1RfTUFOQUdFTUVOVCxNT05JVE9SX01BSU5fREFTSEJPQVJELE1PTklUT1JfTUVUUklDX0NPTkZJRyxNT05JVE9SX0NVU1RPTV9EQVNIQk9BUkQsTU9OSVRPUl9BTEFSTV9DT05GSUcsTU9OSVRPUl9BTEFSTV9NQU5BR0VNRU5ULENPTExBQk9SQVRJT05fUExVR0lOX01BTkFHRU1FTlQsQ09MTEFCT1JBVElPTl9XT1JLRkxPV19PUkNIRVNUUkFUSU9OLEFETUlOX1NZU1RFTV9QQVJBTVMsQURNSU5fUkVTT1VSQ0VTX01BTkFHRU1FTlQsQURNSU5fVVNFUl9ST0xFX01BTkFHRU1FTlQsQURNSU5fQ01EQl9NT0RFTF9NQU5BR0VNRU5ULENNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsQURNSU5fUVVFUllfTE9HLE1FTlVfQURNSU5fUEVSTUlTU0lPTl9NQU5BR0VNRU5ULE1FTlVfREVTSUdOSU5HX0NJX0RBVEFfRU5RVUlSWSxNRU5VX0RFU0lHTklOR19DSV9JTlRFR1JBVEVEX1FVRVJZX0VYRUNVVElPTixNRU5VX0NNREJfREVTSUdOSU5HX0VOVU1fRU5RVUlSWSxNRU5VX0RFU0lHTklOR19DSV9EQVRBX01BTkFHRU1FTlQsTUVOVV9ERVNJR05JTkdfQ0lfSU5URUdSQVRFRF9RVUVSWV9NQU5BR0VNRU5ULE1FTlVfQ01EQl9ERVNJR05JTkdfRU5VTV9NQU5BR0VNRU5ULE1FTlVfSURDX1BMQU5OSU5HX0RFU0lHTixNRU5VX0lEQ19SRVNPVVJDRV9QTEFOTklORyxNRU5VX0FQUExJQ0FUSU9OX0FSQ0hJVEVDVFVSRV9ERVNJR04sTUVOVV9BUFBMSUNBVElPTl9ERVBMT1lNRU5UX0RFU0lHTixNRU5VX0FETUlOX0NNREJfTU9ERUxfTUFOQUdFTUVOVCxNRU5VX0NNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsTUVOVV9BRE1JTl9RVUVSWV9MT0csSk9CU19TRVJWSUNFX0NBVEFMT0dfTUFOQUdFTUVOVCxKT0JTX1RBU0tfTUFOQUdFTUVOVF0ifQ.XbPpkiS6AG7zSLHYxFacU3gnyQMWcIvxqXbI3MSlxTGQqJDWrdPUCyyvE0lfJrPoG69GC2gI25Ys_WyGA71E8A`

type coreHostDto struct {
	Status  string  `json:"status"`
	Message  string  `json:"message"`
	Data  []string  `json:"data"`
}

func SyncClusterList()  {
	if CoreUrl == "" || MasterHostIp == "" {
		logrus.Infof("sync cluster quit,core url or master ip is empty \n")
		return
	}
	request,_ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/available-container-hosts", CoreUrl), strings.NewReader(""))
	request.Header.Set("Authorization", TmpCoreToken)
	resp,err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Errorf("sync cluster list,get response error %v \n", err)
		return
	}
	var coreResult coreHostDto
	defer resp.Body.Close()
	b,_ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &coreResult)
	if err != nil {
		logrus.Errorf("sync cluster list,response json unmarshal error %v \n", err)
		return
	}
	if len(coreResult.Data) == 0 {
		logrus.Infof("sync cluster list,get plugin host list empty \n")
		return
	}
	var tmpClusterIps string
	for _,v := range coreResult.Data {
		if v != MasterHostIp {
			tmpClusterIps += fmt.Sprintf("  - %s\n", v)
		}
	}
	minionByte,err := ioutil.ReadFile("/srv/salt/minions/conf/minion")
	if err != nil {
		logrus.Errorf("read minion file error : %v \n", err)
		return
	}
	sb := string(minionByte)
	if strings.Contains(sb, MasterHostIp) {
		sb = strings.Replace(sb, MasterHostIp, fmt.Sprintf("%s\n%s", MasterHostIp, tmpClusterIps), -1)
	}else{
		logrus.Infof("read minion file,can not find master ip %s \n", MasterHostIp)
		return
	}
	err = ioutil.WriteFile("/srv/salt/minions/conf/minion", []byte(sb), 0644)
	if err != nil {
		logrus.Errorf("write minion file error : %v \n", err)
		return
	}
	logrus.Infof("replace minion file success !! \n")
}