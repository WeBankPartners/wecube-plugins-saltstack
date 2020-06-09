package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"os/exec"
)

const TmpCoreToken = `Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJTWVNfU0FMVFNUQUNLIiwiaWF0IjoxNTkwMTE4MjYxLCJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJjbGllbnRUeXBlIjoiU1VCX1NZU1RFTSIsImV4cCI6MTc0NTYzODI2MSwiYXV0aG9yaXR5IjoiW1NVQl9TWVNURU1dIn0.N2sD9F4TKh1yaatRfr-sqRqlP7fiSqZ1znmr7AtQanr2ZmbldZt2ICeuUnIUcpGGK3YZKKqOPic2JNeECblgnw`
const ClusterDeletePath = `/salt/cluster/delete`

type coreHostDto struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func SyncClusterList() {
	defer func() {
		cpErr := exec.Command("bash","-c", "/bin/cp -f /srv/salt/minions/conf/minion /var/www/html/salt-minion/conf/").Run()
		if cpErr != nil {
			logrus.Errorf("copy /srv/salt/minions/conf/minion /var/www/html/salt-minion/conf/minion error %v \n ", cpErr)
		}
	}()
	if CoreUrl == "" || MasterHostIp == "" {
		logrus.Infof("sync cluster quit,core url or master ip is empty \n")
		return
	}
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/available-container-hosts", CoreUrl), strings.NewReader(""))
	request.Header.Set("Authorization", TmpCoreToken)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Errorf("sync cluster list,get response error %v \n", err)
		return
	}
	var coreResult coreHostDto
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
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
		logrus.Errorf("sync cluster error,can not find masterIp:%s in cluster:%s ", MasterHostIp, ClusterList)
		ClusterList = []string{}
		return
	}
	minionByte, err := ioutil.ReadFile("/srv/salt/minions/conf/minion")
	if err != nil {
		logrus.Errorf("read minion file error : %v \n", err)
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
			logrus.Infof("sync cluster abort with the right config ")
			return
		}
		sb = strings.Replace(sb, MasterHostIp, fmt.Sprintf("%s\n%s", MasterHostIp, tmpClusterIps), -1)
	} else {
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
			logrus.Errorf("send host delete to %s error %v \n", v, err)
		} else {
			var resultBody ClusterHostDeleteResult
			b, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(b, &resultBody)
			if err != nil {
				logrus.Errorf("send host delete to %s fail, json unmarshal error %v \n", v, err)
			} else {
				if resultBody.Code > 0 {
					logrus.Errorf("send host delete to %s fail, code: %d, message: %s", v, resultBody.Code, resultBody.Message)
				} else {
					logrus.Infof("send host delete to %s success ", v)
				}
			}
			resp.Body.Close()
		}
	}
}
