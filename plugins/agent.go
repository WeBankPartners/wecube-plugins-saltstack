package plugins

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var AgentActions = make(map[string]Action)

func init() {
	AgentActions["install"] = new(AgentInstallAction)
	AgentActions["uninstall"] = new(AgentUninstallAction)
}

type AgentPlugin struct {
}

func (plugin *AgentPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := AgentActions[actionName]
	if !found {
		return nil, fmt.Errorf("Agent plugin,action = %s not found", actionName)
	}

	return action, nil
}

type AgentInstallInputs struct {
	Inputs []AgentInstallInput `json:"inputs,omitempty"`
}

type AgentInstallInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
}

type AgentInstallOutputs struct {
	Outputs []AgentInstallOutput `json:"outputs,omitempty"`
}

type AgentInstallOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
	//Detail string `json:"detail,omitempty"`
}

type Roster struct {
	Name   string
	Host   string
	User   string
	Passwd string
	Sudo   string
}

type AgentInstallAction struct {
}

func (action *AgentInstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentInstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *AgentInstallAction) CheckParam(input AgentInstallInput) error {
	if input.Host == "" {
		return errors.New("Host is empty")
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty")
	}

	if input.Password == "" {
		return errors.New("Password is empty")
	}

	return nil
}

func runBashScript(shellPath string, args []string) (string, error) {
	cmd := exec.Command(shellPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("failed to runBashScript(%s), err = %v,args=%v,out=%v", shellPath, err, args, string(out))
		return "", err
	}

	logrus.Infof("runBashScript,output=%s", string(out))
	return string(out), nil
}

func removeSaltKeys(host string) {
	runBashScript("./scripts/salt/remove_master_unused_key.sh", []string{host})
	return
}

func (action *AgentInstallAction) installAgent(input *AgentInstallInput) (output AgentInstallOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.CheckParam(*input)
	if err != nil {
		return output, err
	}

	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		logrus.Errorf("AesDePassword meet error(%v)", err)
		return output, err
	}
	if input.User == "" {
		input.User = "root"
	}

	installMinionArgs := []string{
		input.Host,
		password,
		input.User,
	}
	if input.Port != "" {
		installMinionArgs = append(installMinionArgs, input.Port)
	}
	out, er := runBashScript("./scripts/salt/install_minion.sh", installMinionArgs)
	if er != nil {
		err = fmt.Errorf("failed to install salt-minion, err = %v,out=%v", er, string(out))
		return output, err
	}
	logrus.Infof("installAgent run install_minion.sh out=%v", string(out))

	return output, err
}

func (action *AgentInstallAction) Do(input interface{}) (interface{}, error) {
	agents, _ := input.(AgentInstallInputs)
	outputs := AgentInstallOutputs{}
	var finalErr error
	for _, agent := range agents.Inputs {
		removeSaltKeys(agent.Host)
		agentInstallOutput, err := action.installAgent(&agent)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, agentInstallOutput)
	}

	logrus.Infof("all agents = %v have been installed", agents)
	return &outputs, finalErr
}

type AgentUninstallAction struct {
}

type AgentUninstallInputs struct {
	Inputs []AgentUninstallInput `json:"inputs,omitempty"`
}

type AgentUninstallInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
}

type AgentUninstallOutputs struct {
	Outputs []AgentUninstallOutput `json:"outputs,omitempty"`
}

type AgentUninstallOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *AgentUninstallAction) agentUninstallCheckParam(input AgentUninstallInput) error {
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Host == "" {
		return errors.New("Host is empty")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty")
	}
	if input.Password == "" {
		return errors.New("Password is empty")
	}
	return nil
}

func (action *AgentUninstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentUninstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *AgentUninstallAction) agentUninstall(input *AgentUninstallInput) (output AgentUninstallOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	err = action.agentUninstallCheckParam(*input)
	if err != nil {
		return output, err
	}

	// Decrypt Password
	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		logrus.Errorf("AesDePassword meet error(%v)", err)
		return output, err
	}

	// run uninstall_minion.sh
	uninstallMinionArgs := []string{
		input.Host,
		password,
	}
	out, er := runBashScript("./scripts/salt/uninstall_minion.sh", uninstallMinionArgs)
	if er != nil {
		err = fmt.Errorf("failed to uninstall salt-minion, err = %v,out=%v", er, string(out))
		return output, err
	}
	logrus.Infof("agentUninstall run uninstall_minion.sh out=%v", string(out))

	return output, err
}

func (action *AgentUninstallAction) Do(input interface{}) (interface{}, error) {
	agents, _ := input.(AgentUninstallInputs)
	outputs := AgentUninstallOutputs{}
	var finalErr error
	var hosts []string
	for _, agent := range agents.Inputs {
		agentUninstallOutput, err := action.agentUninstall(&agent)
		if err != nil {
			finalErr = err
		}

		// salt-key -d
		removeSaltKeys(agent.Host)
		outputs.Outputs = append(outputs.Outputs, agentUninstallOutput)
		hosts = append(hosts, agent.Host)
	}
	go SendHostDelete(hosts)

	logrus.Infof("all agents = %v have been uninstalled", agents)
	return &outputs, finalErr
}
