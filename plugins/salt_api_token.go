package plugins

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	ENV_NAME_SALT_API_HOST = "SALT_API_HOST"
	ENV_NAME_SALT_API_USER = "SALT_API_USER"
	ENV_NAME_SALT_API_PWD  = "SALT_API_PWD"
	SALT_TOKEN_VALID_TIME =  30*60   //second
)

type SaltApiToken struct {
	Token      string
	CreateTime time.Time
	Expire     float64 //unit second
}

func (token *SaltApiToken) isSaltApiTokenValid() bool {
	if token.Token == "" {
		return false
	}

	if elapse := time.Since(token.CreateTime).Seconds();  elapse < float64(token.Expire - 10)  {
		return true
	}
	return false
}

var (
	tokenMutex sync.Mutex
	saltToken  SaltApiToken
)

//{"return": [{"perms": [".*"], "start": 1557300494.752942, "token": "f37858551ee3cb4f9d7f1653545e627215e1aaa5", "expire": 1557343694.752943, "user": "saltapi", "eauth": "pam"}]}

type SaltApiTokenResult struct {
	Token  string  `json:"token"`
	Start  float64 `json:"start"`
	Expire float64 `json:"expire"`
}
type NewSaltApiTokenRsp struct {
	Result []SaltApiTokenResult `json:"return"`
}

func newSaltApiToken() error {
	urlPath := "https://127.0.0.1:8080"
	userName := "saltapi"
	passwd := "saltapi"

	if urlPath == "" || userName == "" || passwd == "" {
		return fmt.Errorf("newSaltApiToken:meet empty env param")
	}

	data := url.Values{}
	data.Set("eauth", "pam")
	data.Set("username", userName)
	data.Set("password", passwd)

	request, err := http.NewRequest("POST", urlPath+"/login", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("newRequest meet error=%v,url=%v", err, urlPath)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("newRequest meet error=%v,url=%v", err, urlPath)
		return fmt.Errorf("newRequest meet error=%v,url=%v", err, urlPath)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logrus.Errorf("newSaltApiToken StatusCode != 200,statusCode=%v,url=%v", resp.StatusCode, urlPath)
		return fmt.Errorf("newSaltApiToken StatusCode != 200,statusCode=%v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("newSaltApiToken,readAll meet err=%v", err)
		return fmt.Errorf("newSaltApiToken,readAll meet err=%v", err)
	}

	logrus.Infof("newSaltApiToken http result=%s", []byte(body))
	rtnData := NewSaltApiTokenRsp{}
	if err = json.Unmarshal(body, &rtnData); err != nil {
		logrus.Errorf("newSaltApiToken unmarshal meet err=%v", err)
		return fmt.Errorf("newSaltApiToken unmarshal meet err=%v", err)
	}
	if len(rtnData.Result) != 1 {
		logrus.Errorf("newSaltApiToken len(result)=%v", len(rtnData.Result))
		return fmt.Errorf("newSaltApiToken len(result)=%v", len(rtnData.Result))
	}

	saltToken.Token = rtnData.Result[0].Token
	saltToken.CreateTime = time.Now()
	saltToken.Expire = SALT_TOKEN_VALID_TIME
	//saltToken.Expire = rtnData.Result[0].Expire - rtnData.Result[0].Start

	logrus.Infof("newSaltApiToken ok,token=%++v", saltToken)

	return nil
}

func getSaltApiToken() (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if saltToken.isSaltApiTokenValid() {
		return saltToken.Token, nil
	}

	err := newSaltApiToken()
	if err != nil {
		return "", err
	}

	return saltToken.Token, nil
}
