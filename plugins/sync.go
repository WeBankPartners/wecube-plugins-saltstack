package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const TmpCoreToken = `Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ3ZHNfc3lzdGVtIiwiaWF0IjoxNTgyMDkyODYyLCJ0eXBlIjoiYWNjZXNzVG9rZW4iLCJjbGllbnRUeXBlIjoiVVNFUiIsImV4cCI6MzQ3NDI1Mjg2MiwiYXV0aG9yaXR5IjoiW2FkbWluLElNUExFTUVOVEFUSU9OX1dPUktGTE9XX0VYRUNVVElPTixJTVBMRU1FTlRBVElPTl9CQVRDSF9FWEVDVVRJT04sSU1QTEVNRU5UQVRJT05fQVJUSUZBQ1RfTUFOQUdFTUVOVCxNT05JVE9SX01BSU5fREFTSEJPQVJELE1PTklUT1JfTUVUUklDX0NPTkZJRyxNT05JVE9SX0NVU1RPTV9EQVNIQk9BUkQsTU9OSVRPUl9BTEFSTV9DT05GSUcsTU9OSVRPUl9BTEFSTV9NQU5BR0VNRU5ULENPTExBQk9SQVRJT05fUExVR0lOX01BTkFHRU1FTlQsQ09MTEFCT1JBVElPTl9XT1JLRkxPV19PUkNIRVNUUkFUSU9OLEFETUlOX1NZU1RFTV9QQVJBTVMsQURNSU5fUkVTT1VSQ0VTX01BTkFHRU1FTlQsQURNSU5fVVNFUl9ST0xFX01BTkFHRU1FTlQsQURNSU5fQ01EQl9NT0RFTF9NQU5BR0VNRU5ULENNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsQURNSU5fUVVFUllfTE9HLE1FTlVfQURNSU5fUEVSTUlTU0lPTl9NQU5BR0VNRU5ULE1FTlVfREVTSUdOSU5HX0NJX0RBVEFfRU5RVUlSWSxNRU5VX0RFU0lHTklOR19DSV9JTlRFR1JBVEVEX1FVRVJZX0VYRUNVVElPTixNRU5VX0NNREJfREVTSUdOSU5HX0VOVU1fRU5RVUlSWSxNRU5VX0RFU0lHTklOR19DSV9EQVRBX01BTkFHRU1FTlQsTUVOVV9ERVNJR05JTkdfQ0lfSU5URUdSQVRFRF9RVUVSWV9NQU5BR0VNRU5ULE1FTlVfQ01EQl9ERVNJR05JTkdfRU5VTV9NQU5BR0VNRU5ULE1FTlVfSURDX1BMQU5OSU5HX0RFU0lHTixNRU5VX0lEQ19SRVNPVVJDRV9QTEFOTklORyxNRU5VX0FQUExJQ0FUSU9OX0FSQ0hJVEVDVFVSRV9ERVNJR04sTUVOVV9BUFBMSUNBVElPTl9ERVBMT1lNRU5UX0RFU0lHTixNRU5VX0FETUlOX0NNREJfTU9ERUxfTUFOQUdFTUVOVCxNRU5VX0NNREJfQURNSU5fQkFTRV9EQVRBX01BTkFHRU1FTlQsTUVOVV9BRE1JTl9RVUVSWV9MT0csSk9CU19TRVJWSUNFX0NBVEFMT0dfTUFOQUdFTUVOVCxKT0JTX1RBU0tfTUFOQUdFTUVOVF0ifQ.XbPpkiS6AG7zSLHYxFacU3gnyQMWcIvxqXbI3MSlxTGQqJDWrdPUCyyvE0lfJrPoG69GC2gI25Ys_WyGA71E8A`
const ClusterDeletePath = `/salt/cluster/delete`

type coreHostDto struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func SyncClusterList() {
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
