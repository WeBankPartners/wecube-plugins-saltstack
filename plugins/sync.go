package plugins

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const ClusterDeletePath = `/salt/cluster/delete`

var (
	coreRequestToken string
	requestCoreNonce = "salt"
)

type coreHostDto struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func SyncClusterList() {
	defer func() {
		cpErr := exec.Command("bash","-c", "/bin/cp -f /srv/salt/minions/conf/minion /var/www/html/salt-minion/conf/").Run()
		if cpErr != nil {
			log.Logger.Error(fmt.Sprintf("copy /srv/salt/minions/conf/minion /var/www/html/salt-minion/conf/minion error %v \n ", cpErr))
		}
	}()
	if CoreUrl == "" || MasterHostIp == "" {
		log.Logger.Warn("sync cluster quit,core url or master ip is empty \n")
		return
	}
	if SubSystemCode == "" || SubSystemKey == "" {
		log.Logger.Warn("Sync cluster list fail,subSystemCore or subSystemKey is empty", log.String("SubSystemCode", SubSystemCode), log.String("SubSystemKey", SubSystemKey))
		return
	}
	err := initCoreToken(SubSystemKey)
	if err != nil || coreRequestToken == "" {
		log.Logger.Warn("Sync cluster list fail,init core token fail", log.Error(err))
		return
	}
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/available-container-hosts", CoreUrl), strings.NewReader(""))
	request.Header.Set("Authorization", coreRequestToken)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("sync cluster list,get response error %v \n", err))
		return
	}
	var coreResult coreHostDto
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &coreResult)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("sync cluster list,response json unmarshal error %v \n", err))
		return
	}
	if len(coreResult.Data) == 0 {
		log.Logger.Info("sync cluster list,get plugin host list empty \n")
		return
	}
	var tmpClusterIps string
	var problemCluster bool
	for _, v := range coreResult.Data {
		if v != MasterHostIp {
			//tmpClusterIps += fmt.Sprintf("  - %s\n", v)
			ClusterList = append(ClusterList, v)
		}else{
			problemCluster = true
		}
	}
	if !problemCluster {
		log.Logger.Error(fmt.Sprintf("sync cluster error,can not find masterIp:%s in cluster:%s ", MasterHostIp, ClusterList))
		ClusterList = []string{}
		return
	}
	minionByte, err := ioutil.ReadFile("/srv/salt/minions/conf/minion")
	if err != nil {
		log.Logger.Error(fmt.Sprintf("read minion file error : %v \n", err))
		return
	}
	sb := string(minionByte)
	if strings.Contains(sb, MasterHostIp) {
		for _,v := range ClusterList {
			if v != MasterHostIp && v != "" {
				if !strings.Contains(sb, v) {
					tmpClusterIps += fmt.Sprintf("  - %s\n", v)
				}
			}
		}
		if tmpClusterIps == "" {
			log.Logger.Info("sync cluster abort with the right config ")
			return
		}
		sb = strings.Replace(sb, MasterHostIp, fmt.Sprintf("%s\n%s", MasterHostIp, tmpClusterIps), -1)
	} else {
		log.Logger.Info(fmt.Sprintf("read minion file,can not find master ip %s \n", MasterHostIp))
		return
	}
	err = ioutil.WriteFile("/srv/salt/minions/conf/minion", []byte(sb), 0644)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("write minion file error : %v \n", err))
		return
	}
	log.Logger.Info("replace minion file success !! \n")
}

func StartClusterServer() {
	http.Handle(ClusterDeletePath, http.HandlerFunc(handleHostDelete))
	http.ListenAndServe(":4507", nil)
}

type ClusterHostDeleteResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleHostDelete(w http.ResponseWriter, r *http.Request) {
	result := ClusterHostDeleteResult{Code: 0, Message: "success"}
	hostIp := r.FormValue("hosts")
	if hostIp == "" {
		result.Code = 1
		result.Message = "Param hosts empty"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		for _, v := range strings.Split(hostIp, ",") {
			if strings.Count(v, ".") != 3 {
				result.Code = 2
				result.Message += fmt.Sprintf(" ip:%s illegal ", v)
				continue
			}
			removeSaltKeys(v)
		}
	}
	rb, _ := json.Marshal(result)
	w.Write(rb)
}

func SendHostDelete(hosts []string) {
	if len(hosts) == 0 {
		return
	}
	urlParam := fmt.Sprintf("hosts=%s", strings.Join(hosts, ","))
	for _, v := range ClusterList {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:4507%s?%s", v, ClusterDeletePath, urlParam), strings.NewReader(""))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("send host delete to %s error %v \n", v, err))
		} else {
			var resultBody ClusterHostDeleteResult
			b, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(b, &resultBody)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("send host delete to %s fail, json unmarshal error %v \n", v, err))
			} else {
				if resultBody.Code > 0 {
					log.Logger.Error(fmt.Sprintf("send host delete to %s fail, code: %d, message: %s", v, resultBody.Code, resultBody.Message))
				} else {
					log.Logger.Error(fmt.Sprintf("send host delete to %s success ", v))
				}
			}
			resp.Body.Close()
		}
	}
}

type requestCoreToken struct {
	Password  string  `json:"password"`
	Username  string  `json:"username"`
	Nonce     string  `json:"nonce"`
	ClientType string  `json:"clientType"`
}

type responseCoreObj struct {
	Status  string  `json:"status"`
	Message  string  `json:"message"`
	Data  []*responseCoreDataObj  `json:"data"`
}

type responseCoreDataObj struct {
	Expiration  string  `json:"expiration"`
	Token  string  `json:"token"`
	TokenType  string  `json:"tokenType"`
}

func initCoreToken(rsaKey string) error {
	encryptBytes,err := RSAEncryptByPrivate([]byte(fmt.Sprintf("%s:%s", SubSystemCode, requestCoreNonce)), rsaKey)
	encryptString := base64.StdEncoding.EncodeToString(encryptBytes)
	if err != nil {
		return err
	}
	postParam := requestCoreToken{Username: SubSystemCode, Nonce: requestCoreNonce, ClientType: "SUB_SYSTEM", Password: encryptString}
	postBytes,_ := json.Marshal(postParam)
	fmt.Printf("param: %s \n", string(postBytes))
	req,err := http.NewRequest(http.MethodPost, CoreUrl + "/auth/v1/api/login", bytes.NewReader(postBytes))
	if err != nil {
		return fmt.Errorf("http new request fail,%s ", err.Error())
	}
	resp,err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http response fail, %s ", err.Error())
	}
	var respObj responseCoreObj
	bodyBytes,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(bodyBytes, &respObj)
	if err != nil {
		return fmt.Errorf("http response body read fail,%s ", err.Error())
	}
	for _,v := range respObj.Data {
		if v.TokenType == "accessToken" {
			coreRequestToken = v.Token
		}
	}
	return nil
}