package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"flag"
	"fmt"
	"time"
	"github.com/WeBankPartners/wecube-plugins-saltstack/plugins"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/models"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

func init() {
	initConfig()
	log.InitZapLogger()
	initRouter()
}

func main() {
	plugins.InitErrorMessageList()
	plugins.InitEnvParam()
	plugins.SyncClusterList()
	go plugins.StartClusterServer()

	if err := http.ListenAndServe(":"+models.Config.Http.Port, nil); err != nil {
		log.Logger.Fatal("Start listening error", log.Error(err))
	}else{
		log.Logger.Info(fmt.Sprintf("Listening %s ...", models.Config.Http.Port))
	}
}

func initRouter() {
	//path should be defined as "/[package]/[version]/[plugin]/[action]"
	http.HandleFunc("/saltstack/v1/", routeDispatcher)
	http.HandleFunc("/v1/deploy/webconsole", plugins.WebConsoleHandler)
	http.HandleFunc("/v1/deploy/webconsoleStaticPage", plugins.WebConsoleStaticPageHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func initConfig()  {
	cfgFile := flag.String("c", "conf/default.json", "config file")
	flag.Parse()
	err := models.InitConfig(*cfgFile)
	if err != nil {
		fmt.Printf("Init config fail,%s \n", err.Error())
		os.Exit(1)
	}
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	pluginRequest := parsePluginRequest(r)
	pluginResponse, _ := plugins.Process(pluginRequest)
	if pluginResponse.ResultCode == "1" {
		log.Logger.Error("Handle error", log.JsonObj("response", pluginResponse))
	}else{
		log.Logger.Debug("Handle success", log.JsonObj("response", pluginResponse))
	}
	log.Logger.Info("Request",log.String("url", r.RequestURI), log.String("method",r.Method), log.String("ip",strings.Split(r.RemoteAddr,":")[0]), log.Float64("cost_second",time.Now().Sub(start).Seconds()))
	write(w, pluginResponse)
}

func write(w http.ResponseWriter, output *plugins.PluginResponse) {
	w.Header().Set("content-type", "application/json")
	b, _ := json.Marshal(output)
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