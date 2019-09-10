package plugins

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"github.com/sirupsen/logrus"
)

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
	return  []byte{}
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
        defer func(){
             	if r:=recover();r!=nil{
                        err=fmt.Errorf("%v",r)
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
	    err=fmt.Errorf("password wrong")
	    return 
	}
	
	password=string(origData)
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

	saltResult,err:=parseSaltApiCallResult(result)
	if err != nil {
		logrus.Infof("parseSaltApiCallResult meet error=%v ", err)
		return "", err
	}
	
	if len(saltResult.Results) == 0  || len(saltResult.Results[0]) == 0 {
		return "",fmt.Errorf("salt api:no target match ,please check if salt-agent installed on target,reqeust=%v",request)
	}

	return result, nil
}
