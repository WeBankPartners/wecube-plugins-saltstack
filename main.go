package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/models"
	"github.com/WeBankPartners/wecube-plugins-saltstack/plugins"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	go plugins.StartCleanInterval()

	if err := http.ListenAndServe(":"+models.Config.Http.Port, nil); err != nil {
		log.Logger.Fatal("Start listening error", log.Error(err))
	} else {
		log.Logger.Info(fmt.Sprintf("Listening %s ...", models.Config.Http.Port))
	}
}

func initRouter() {
	//path should be defined as "/[package]/[version]/[plugin]/[action]"
	models.CoreJwtKey = plugins.DecryptRsa(models.CoreJwtKey)
	http.HandleFunc("/saltstack/v1/", routeDispatcher)
	http.HandleFunc("/v1/deploy/webconsole", plugins.WebConsoleHandler)
	http.HandleFunc("/v1/deploy/webconsoleStaticPage", plugins.WebConsoleStaticPageHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func initConfig() {
	cfgFile := flag.String("c", "conf/default.json", "config file")
	flag.Parse()
	err := models.InitConfig(*cfgFile)
	if err != nil {
		fmt.Printf("Init config fail,%s \n", err.Error())
		os.Exit(1)
	}
}

func routeDispatcher(w http.ResponseWriter, r *http.Request) {
	if authCore(r.Header.Get("Authorization")) {
		start := time.Now()
		pluginRequest := parsePluginRequest(r)
		pluginResponse, _ := plugins.Process(pluginRequest)
		if pluginResponse.ResultCode == "1" {
			log.Logger.Error("Handle error", log.JsonObj("response", pluginResponse))
		} else {
			log.Logger.Debug("Handle success", log.JsonObj("response", pluginResponse))
		}
		log.Logger.Info("Request end ----------------<<", log.String("url", r.RequestURI), log.String("method", r.Method), log.String("ip", strings.Split(r.RemoteAddr, ":")[0]), log.Float64("cost_second", time.Now().Sub(start).Seconds()))
		write(w, pluginResponse)
	} else {
		log.Logger.Info("Request token illegal ----------------!!", log.String("url", r.RequestURI), log.String("method", r.Method), log.String("ip", strings.Split(r.RemoteAddr, ":")[0]))
		pluginResponse := plugins.PluginResponse{ResultCode: "1", ResultMsg: "Token illegal"}
		write(w, &pluginResponse)
	}
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

func authCore(coreToken string) bool {
	tokenObj, err := decodeCoreToken(coreToken, models.CoreJwtKey)
	if err == nil {
		isSystemCall := false
		for _, v := range tokenObj.Roles {
			if v == plugins.SystemRole {
				isSystemCall = true
				break
			}
		}
		if !isSystemCall {
			log.Logger.Warn("token illegal", log.JsonObj("token", tokenObj))
		}
		return isSystemCall
	}
	return false
}

func decodeCoreToken(token, key string) (result models.CoreJwtToken, err error) {
	if strings.HasPrefix(token, "Bearer") {
		token = token[7:]
	}
	if key == "" || strings.HasPrefix(key, "{{") {
		key = "Platform+Auth+Server+Secret"
	}
	keyBytes, err := ioutil.ReadAll(base64.NewDecoder(base64.RawStdEncoding, bytes.NewBufferString(key)))
	if err != nil {
		log.Logger.Error("Decode core token fail,base64 decode error", log.Error(err))
		return result, err
	}
	pToken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return keyBytes, nil
	})
	if err != nil {
		log.Logger.Error("Decode core token fail,jwt parse error", log.Error(err))
		return result, err
	}
	claimMap, ok := pToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Logger.Error("Decode core token fail,claims to map error", log.Error(err))
		return result, err
	}
	result.User = fmt.Sprintf("%s", claimMap["sub"])
	result.Expire, err = strconv.ParseInt(fmt.Sprintf("%.0f", claimMap["exp"]), 10, 64)
	if err != nil {
		log.Logger.Error("Decode core token fail,parse expire to int64 error", log.Error(err))
		return result, err
	}
	roleListString := fmt.Sprintf("%s", claimMap["authority"])
	roleListString = roleListString[1 : len(roleListString)-1]
	result.Roles = strings.Split(roleListString, ",")
	return result, nil
}
