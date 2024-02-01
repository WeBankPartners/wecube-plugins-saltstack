package plugins

import (
	"fmt"
	"sync"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

var (
	pluginsMutex sync.Mutex
	plugins      = make(map[string]Plugin)
)

type Plugin interface {
	GetActionByName(actionName string) (Action, error)
}

type Action interface {
	ReadParam(param interface{}) (interface{}, error)
	// CheckParam(param interface{}) error
	Do(param interface{}) (interface{}, error)
	// Set accept language
	SetAcceptLanguage(language string)
}

func registerPlugin(name string, plugin Plugin) {
	pluginsMutex.Lock()
	defer pluginsMutex.Unlock()

	if _, found := plugins[name]; found {
		log.Logger.Fatal("deploy plugin twice", log.String("plugin", name))
	}

	plugins[name] = plugin
}

func getPluginByName(name string) (Plugin, error) {
	pluginsMutex.Lock()
	defer pluginsMutex.Unlock()
	plugin, found := plugins[name]
	if !found {
		return nil, fmt.Errorf("plugin[%s] not found", name)
	}
	return plugin, nil
}

func init() {
	registerPlugin("host-file", new(FilePlugin))
	registerPlugin("salt-api", new(SaltApiPlugin))
	registerPlugin("agent", new(AgentPlugin))
	registerPlugin("package-variable", new(VariablePlugin))
	registerPlugin("host-script", new(ScriptPlugin))
	registerPlugin("host-user", new(UserPlugin))
	registerPlugin("mysql-database", new(MysqlDatabasePlugin))
	registerPlugin("mysql-script", new(MysqlScriptPlugin))
	registerPlugin("mysql-user", new(MysqlUserPlugin))
	registerPlugin("released-package", new(ReleasedPackagePlugin))
	registerPlugin("text-processor", new(TextProcessorPlugin))
	//registerPlugin("log", new(LogPlugin))
	registerPlugin("apply-deployment", new(ApplyDeploymentPlugin))
	registerPlugin("web-console", new(WebConsolePlugin))
	registerPlugin("password", new(PasswordPlugin))
}

type PluginRequest struct {
	Version      string
	ProviderName string
	Name         string
	Action       string
	Parameters   interface{}
}

type PluginResponse struct {
	ResultCode string      `json:"resultCode"`
	ResultMsg  string      `json:"resultMessage"`
	Results    interface{} `json:"results"`
}

func Process(pluginRequest *PluginRequest) (*PluginResponse, error) {
	var pluginResponse = PluginResponse{}
	var err error
	defer func() {
		if err != nil {
			pluginResponse.ResultCode = "1"
			pluginResponse.ResultMsg = fmt.Sprint(err)
		} else {
			pluginResponse.ResultCode = "0"
			pluginResponse.ResultMsg = "success"
		}
	}()

	log.Logger.Info("Request start ---------------->>", log.String("plugin", pluginRequest.Name), log.String("action", pluginRequest.Action))

	plugin, err := getPluginByName(pluginRequest.Name)
	if err != nil {
		return &pluginResponse, err
	}

	action, err := plugin.GetActionByName(pluginRequest.Action)
	if err != nil {
		return &pluginResponse, err
	}

	actionParam, err := action.ReadParam(pluginRequest.Parameters)
	if err != nil {
		return &pluginResponse, err
	}

	action.SetAcceptLanguage("")

	// if err = action.CheckParam(actionParam); err != nil {
	// 	return &pluginResponse, err
	// }

	log.Logger.Debug("Request param", log.JsonObj("param", actionParam))
	outputs, err := action.Do(actionParam)
	if err != nil {
		log.Logger.Error("Action handle error", log.Error(err))
	}

	pluginResponse.Results = outputs

	return &pluginResponse, err
}
