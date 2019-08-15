package plugins

import (
	"fmt"
	"os/exec"
    "errors "
	"github.com/sirupsen/logrus"
)

var AgentActions = make(map[string]Action)

func init() {
	AgentActions["install"] = new(AgentInstallAction)
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
	Guid  string `json:"guid,omitempty"`
	Seed  string `json:"seed,omitempty"`
	Password  string `json:"password,omitempty"`
	Host string `json:"host,omitempty"`
}

type AgentInstallOutputs struct {
	Outputs []AgentInstallOutput `json:"outputs,omitempty"`
}

type AgentInstallOutput struct {
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
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

func (action *AgentInstallAction) CheckParam(input interface{}) error {
	inputs, ok := input.(AgentInstallInputs)
	if !ok {
		return fmt.Errorf("agentInstallAtion:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
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

func (action *AgentInstallAction) removeSaltKeys(host string) {
	runBashScript("./scripts/salt/remove_master_unused_key.sh", []string{host})
	return
}

func (action *AgentInstallAction) installAgent(input *AgentInstallInput) (*AgentInstallOutput, error) {
	md5sum := Md5Encode(input.Guid+input.Seed)
	password,err := AesDecode(md5sum[0:16], input.Password)
	if err != nil {
		logrus.Errorf("AesDecode meet error(%v)", err)
		return nil , err
	}
	installMinionArgs:=[]string{
		input.Host,
		password,
	}
	out, err := runBashScript("./scripts/salt/install_minion.sh", installMinionArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to install salt-minion, err = %v,out=%v", err, string(out))
	}

	output := AgentInstallOutput{}
	output.Guid = input.Guid
	output.Detail = string(out)
	return &output, nil
}

func (action *AgentInstallAction) Do(input interface{}) (interface{}, error) {
	agents, _ := input.(AgentInstallInputs)
	outputs := AgentInstallOutputs{}
	for _, agent := range agents.Inputs {
		action.removeSaltKeys(agent.Host)
		agentInstallOutput, err := action.installAgent(&agent)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *agentInstallOutput)
	}

	logrus.Infof("all agents = %v have been installed", agents)
	return &outputs, nil
}
