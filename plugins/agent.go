package plugins

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"net"
	gossh "golang.org/x/crypto/ssh"
)

var AgentActions = make(map[string]Action)

func init() {
	AgentActions["install_old"] = new(AgentInstallAction)
	AgentActions["uninstall_old"] = new(AgentUninstallAction)
	AgentActions["install"] = new(MinionInstallAction)
	AgentActions["uninstall"] = new(MinionUninstallAction)
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
	Command  string `json:"command,omitempty"`
	Method   string `json:"method,omitempty"`
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

	if input.Command != "" {
		_,err = execRemote(input.User, password, input.Host, input.Command)
		if err != nil {
			logrus.Errorf("To host: %s Exec command: %s error %v ", input.Host, input.Command, err)
			return output, fmt.Errorf("To host: %s Exec command: %s error %v ", input.Host, input.Command, err)
		}
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

func execRemote(user,password,host,command string) (output []byte,err error) {
	var(
		client *gossh.Client
		session *gossh.Session
	)
	auth := make([]gossh.AuthMethod, 0)
	auth = append(auth, gossh.Password(password))
	clientConfig := &gossh.ClientConfig{
		User: user,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key gossh.PublicKey) error {
			return nil
		},
	}
	if client,err = gossh.Dial("tcp", fmt.Sprintf("%s:22", host), clientConfig); err != nil {
		fmt.Printf("ssh dial error %v \n", err)
		return output,err
	}
	session,err = client.NewSession()
	if err != nil {
		fmt.Printf("ssh client new session error %v \n", err)
		return output,err
	}
	output,err = session.Output(command)
	if err != nil {
		fmt.Printf("ssh run command error %v \n", err)
	}
	session.Close()
	return output,err
}

type MinionInstallAction struct {}

type MinionUninstallAction struct {}

func (action *MinionInstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentInstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MinionUninstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentUninstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MinionInstallAction) CheckParam(input AgentInstallInput) error {
	if input.Host == "" {
		return errors.New("Host is empty ")
	}
	if input.Guid == "" {
		return errors.New("Guid is empty ")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty ")
	}

	if input.Password == "" {
		return errors.New("Password is empty ")
	}

	if MasterHostIp == "" {
		return errors.New("Master ip is empty,please check docker env ")
	}

	return nil
}

func (action *MinionInstallAction) installMinion(input *AgentInstallInput) (output AgentInstallOutput, err error) {
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
	var cmdOut []byte
	if input.Command != "" {
		cmdOut,err = execRemote(input.User, password, input.Host, input.Command)
		if err != nil {
			logrus.Errorf("To host: %s Exec command: %s output: %s error %v ", input.Host, input.Command, string(cmdOut), err)
			return output, fmt.Errorf("To host: %s Exec command: %s output: %s error %v ", input.Host, input.Command, string(cmdOut), err)
		}
	}

	cmdOut,err = execRemote(input.User, password, input.Host, fmt.Sprintf("curl http://%s:9099/salt-minion/minion_install.sh | bash /dev/stdin %s %s %s", MasterHostIp, MasterHostIp, input.Host, input.Method))
	logrus.Infof("Install minion:%s with output: %s ", input.Host, string(cmdOut))
	if err != nil {
		logrus.Errorf("Install minion to host: %s  error %v ", input.Host, err)
		return output, fmt.Errorf("Install minion to host: %s  error %v ", input.Host, err)
	}

	return output, err
}

func (action *MinionInstallAction) Do(input interface{}) (interface{}, error) {
	agents, _ := input.(AgentInstallInputs)
	outputs := AgentInstallOutputs{}
	var finalErr error
	for _, agent := range agents.Inputs {
		b,tmpErr := exec.Command("bash","-c", fmt.Sprintf("find /etc/salt/pki/master/minions -name '%s'", agent.Host)).Output()
		if tmpErr != nil || string(b) != "" {
			removeSaltKeys(agent.Host)
		}
		agentInstallOutput, err := action.installMinion(&agent)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, agentInstallOutput)
	}

	logrus.Infof("all agents = %v have been installed", agents)
	return &outputs, finalErr
}

func (action *MinionUninstallAction) agentUninstallCheckParam(input AgentUninstallInput) error {
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

func (action *MinionUninstallAction) agentUninstall(input *AgentUninstallInput) (output AgentUninstallOutput, err error) {
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

	var cmdOut []byte
	cmdOut,err = execRemote("root", password, input.Host, fmt.Sprintf("curl http://%s:9099/salt-minion/minion_uninstall.sh | bash ", MasterHostIp))
	logrus.Infof("Uninstall minion:%s with output: %s ", input.Host, string(cmdOut))
	if err != nil {
		logrus.Errorf("Uninstall minion from host: %s  error %v ", input.Host, err)
		return output, fmt.Errorf("Uninstall minion from host: %s  error %v ", input.Host, err)
	}

	return output, err
}

func (action *MinionUninstallAction) Do(input interface{}) (interface{}, error) {
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
