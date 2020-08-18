package plugins

import (
	"fmt"
	"net"
	gossh "golang.org/x/crypto/ssh"
	"strings"
	"os/exec"
	"time"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/models"
)

var AgentActions = make(map[string]Action)

func init() {
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


func runBashScript(shellPath string, args []string) (string, error) {
	cmd := exec.Command(shellPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Logger.Error("RunBashScript error", log.String("path", shellPath), log.StringList("args", args), log.String("output", string(out)), log.Error(err))
		return "", err
	}

	log.Logger.Debug("RunBashScript", log.String("output", string(out)))
	return string(out), nil
}

func removeSaltKeys(host string) {
	runBashScript("./scripts/salt/remove_master_unused_key.sh", []string{host})
	return
}


type AgentUninstallInputs struct {
	Inputs []AgentUninstallInput `json:"inputs,omitempty"`
}

type AgentUninstallInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	User     string `json:"user,omitempty"`
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

type ExecRemoteParam struct {
	User  string
	Password  string
	Host  string
	Command  string
	Output  string
	Err  error
	Timeout  int
	DoneChan  chan int
}

func execRemoteWithTimeout(param *ExecRemoteParam)  {
	param.DoneChan = make(chan int)
	go func(gParam *ExecRemoteParam) {
		tmpOutput,tmpError := execRemote(gParam.User,gParam.Password,gParam.Host,gParam.Command)
		gParam.Output = string(tmpOutput)
		gParam.Err = tmpError
		gParam.DoneChan <- 1
	}(param)
	select {
	case <-param.DoneChan: log.Logger.Debug("Exec remote bash done ", log.String("host", param.Host))
	case <-time.After(time.Duration(param.Timeout)*time.Second):
		param.Err = fmt.Errorf("exec remote command timeout %d seconds", param.Timeout)
	}
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
		return output,fmt.Errorf("ssh dial error:%s", err.Error())
	}
	session,err = client.NewSession()
	if err != nil {
		return output,fmt.Errorf("ssh client new session error:%s", err.Error())
	}
	output,err = session.Output(command)
	if err != nil {
		err = fmt.Errorf("ssh run command error:%s", err.Error())
	}
	session.Close()
	return output,err
}

type MinionInstallAction struct { Language string }

type MinionUninstallAction struct { Language string }

func (action *MinionInstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentInstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MinionInstallAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *MinionUninstallAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AgentUninstallInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MinionUninstallAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *MinionInstallAction) CheckParam(input AgentInstallInput) error {
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	if input.Seed == "" {
		return getParamEmptyError(action.Language, "seed")
	}

	if input.User == "" {
		return getParamEmptyError(action.Language, "user")
	}

	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}

	if MasterHostIp == "" {
		return getSysParamEmptyError(action.Language, "minion_master_ip")
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
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	if input.Command != "" {
		tmpParam := ExecRemoteParam{User:input.User,Password:input.Password,Host:input.Host,Command:input.Command,Timeout:models.Config.ExecRemoteCommandTimeout}
		execRemoteWithTimeout(&tmpParam)
		cmdOutString := tmpParam.Output
		err = tmpParam.Err
		if err != nil {
			log.Logger.Error("Exec command", log.String("host", input.Host), log.String("command", input.Command), log.String("output", cmdOutString), log.Error(err))
			err = getRemoteCommandError(action.Language, input.Host, cmdOutString, err)
			return output, err
		}
	}
	execParam := ExecRemoteParam{User:input.User,Password:password,Host:input.Host,Timeout:models.Config.InstallMinionTimeout,Command:fmt.Sprintf("curl http://%s:9099/salt-minion/minion_install.sh | bash /dev/stdin %s %s %s", MasterHostIp, MasterHostIp, input.Host, input.Method)}
	execRemoteWithTimeout(&execParam)
	outString := execParam.Output
	err = execParam.Err
	if err != nil {
		log.Logger.Error("Install minion", log.String("host", input.Host), log.String("output", outString), log.Error(err))
		err = getRemoteCommandError(action.Language, input.Host, outString, err)
		return output, err
	}
	if strings.TrimSpace(outString) == "" {
		err = fmt.Errorf("Remote install salt-minion fail,please check network from target to master with port 9099 ")
		return output, err
	}
	if !strings.Contains(outString, "salt-minion_success") {
		err = getInstallMinionError(action.Language, input.Host, outString)
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
			log.Logger.Error("Install minion action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, agentInstallOutput)
	}

	return &outputs, finalErr
}

func (action *MinionUninstallAction) agentUninstallCheckParam(input AgentUninstallInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if input.Seed == "" {
		return getParamEmptyError(action.Language, "seed")
	}
	if input.User == "" {
		return getParamEmptyError(action.Language, "user")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
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
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	var cmdOut []byte
	cmdOut,err = execRemote(input.User, password, input.Host, fmt.Sprintf("curl http://%s:9099/salt-minion/minion_uninstall.sh | bash ", MasterHostIp))
	log.Logger.Debug("Uninstall minion", log.String("host", input.Host), log.String("output", string(cmdOut)))
	if err != nil {
		err = getUninstallMinionError(action.Language, input.Host, string(cmdOut), err)
		return output, err
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
			log.Logger.Error("Uninstall minion action", log.Error(err))
			finalErr = err
		}

		// salt-key -d
		removeSaltKeys(agent.Host)
		outputs.Outputs = append(outputs.Outputs, agentUninstallOutput)
		hosts = append(hosts, agent.Host)
	}
	go SendHostDelete(hosts)

	return &outputs, finalErr
}
