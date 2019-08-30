package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/plugins"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

const (
	PLUGIN_SERVICE_PORT = "8082"
)

func init() {
	initLogger()
	initRouter()
}

func main() {
	logrus.Infof("Start WeCube-Plungins Deploy Service at port %v ... ", PLUGIN_SERVICE_PORT)

	if err := http.ListenAndServe(":"+PLUGIN_SERVICE_PORT, nil); err != nil {
		logrus.Fatalf("ListenAndServe meet err = %v", err)
	}
}

func initLogger() {
	fileName := "logs/wecube-plugins-saltstack.log"
	logrus.SetReportCaller(true)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     7,
		Level:      logrus.InfoLevel,
		Formatter:  &logrus.TextFormatter{DisableTimestamp: false, DisableColors: false},
	})
	logrus.AddHook(rotateFileHook)
}

func initRouter() {
	//path should be defined as "/[version]/[pacakge]/[plugin]/[action]"
	http.HandleFunc("/v1/deploy/file/copy", routeDispatcher)
	http.HandleFunc("/v1/deploy/salt-api/call", routeDispatcher)
	http.HandleFunc("/v1/deploy/agent/install", routeDispatcher)
	http.HandleFunc("/v1/deploy/variable/replace", routeDispatcher)
	http.HandleFunc("/v1/deploy/script/run", routeDispatcher)
	http.HandleFunc("/v1/deploy/user/add", routeDispatcher)
	http.HandleFunc("/v1/deploy/user/remove", routeDispatcher)
	http.HandleFunc("/v1/deploy/database/runScript", routeDispatcher)
	http.HandleFunc("/v1/deploy/released-package/listCurrentDir", routeDispatcher)
	http.HandleFunc("/v1/deploy/released-package/getConfigFileKey", routeDispatcher)
	http.HandleFunc("/v1/deploy/disk/getUnformatedDisk", routeDispatcher)
	http.HandleFunc("/v1/deploy/disk/formatAndMountDisk", routeDispatcher)
	http.HandleFunc("/v1/deploy/text-processor/search", routeDispatcher)
	http.HandleFunc("/v1/deploy/text-processor/getContext", routeDispatcher)
	http.HandleFunc("/v1/qcloud/log/search", routeDispatcher)
	http.HandleFunc("/v1/qcloud/log/searchdetail", routeDispatcher)
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	pluginRequest := parsePluginRequest(r)
	pluginResponse, _ := plugins.Process(pluginRequest)
	write(w, pluginResponse)
}

func write(w http.ResponseWriter, output *plugins.PluginResponse) {
	w.Header().Set("content-type", "application/json")
	b, err := json.Marshal(output)
	if err != nil {
		logrus.Errorf("write http response (%v) meet error (%v)", output, err)
	}
	w.Write(b)
}

func parsePluginRequest(r *http.Request) *plugins.PluginRequest {
	var pluginInput = plugins.PluginRequest{}
	pathStrings := strings.Split(r.URL.Path, "/")
	if len(pathStrings) >= 5 {
		pluginInput.Version = pathStrings[1]
		pluginInput.ProviderName = pathStrings[2]
		pluginInput.Name = pathStrings[3]
		pluginInput.Action = pathStrings[4]
	}
	pluginInput.Parameters = r.Body
	return &pluginInput
}
