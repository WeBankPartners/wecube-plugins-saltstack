package plugins

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	WEB_CONSOLE_DEFAULT_USER_NAME = "root"
	WEB_CONSOLE_DEFAULT_PORT      = 22
	WEB_CONSOLE_DEFAULT_COLS      = 800
	WEB_CONSOLE_DEFUALT_ROWS      = 600
)

var (
	sshTokenMap          = make(map[string]*ssh)
	webConsoleTokenMutex sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ssh struct {
	user    string
	pwd     string
	addr    string
	client  *gossh.Client
	session *gossh.Session
}

func (s *ssh) Connect() (*ssh, error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = s.user
	config.Auth = []gossh.AuthMethod{gossh.Password(s.pwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	client, err := gossh.Dial("tcp", s.addr, config)
	if nil != err {
		return nil, err
	}
	s.client = client
	return s, nil
}

func (s *ssh) Exec(cmd string) (string, error) {
	var buf bytes.Buffer
	session, err := s.client.NewSession()
	if nil != err {
		return "", err
	}
	session.Stdout = &buf
	session.Stderr = &buf
	err = session.Run(cmd)
	if err != nil {
		return "", err
	}
	defer session.Close()
	stdout := buf.String()
	fmt.Printf("Stdout:%v\n", stdout)
	return stdout, nil
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type RunWebConsoleParam struct {
	Token   string `json:"token,omitempty"`
	Rows    uint32 `json:"rows,omitempty"`
	Columns uint32 `json:"columns,omitempty"`
}

type RunWebConsoleErr struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMessage"`
}

func getRunWebConsoleBytes(err error) []byte {
	consoleErr := RunWebConsoleErr{
		ResultCode: "-1",
		ResultMsg:  err.Error(),
	}
	b, _ := json.Marshal(consoleErr)
	return b
}

func WebConsoleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var param RunWebConsoleParam

	defer func() {
		if err != nil {
			w.Header().Set("content-type", "application/json")
			w.Write(getRunWebConsoleBytes(err))
		}

	}()

	if err = UnmarshalJson(r.Body, &param); err != nil {
		return
	}
	sh, err := getSshInfoByTimeStamp(param.Token)
	if err != nil {
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrader failed err=%v\n", err)
		return
	}
	defer ws.Close()
	sh, err = sh.Connect()
	if err != nil {
		fmt.Printf("ssh connect failed,err=%v\n", err)
		return
	}
	defer func() {
		sh.client.Close()
	}()

	channel, incomingRequests, err := sh.client.Conn.OpenChannel("session", nil)
	if err != nil {
		fmt.Printf("ssh open channel failed,err=%v\n", err)
		return
	}
	defer func() {
		channel.Close()
	}()

	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()

	modes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, gossh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)

	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  param.Columns,
		Rows:     param.Rows,
		Width:    param.Columns * 8,
		Height:   param.Rows * 8,
		Modelist: string(modeList),
	}

	ok, err := channel.SendRequest("pty-req", true, gossh.Marshal(&req))
	if !ok || err != nil {
		fmt.Printf("channel send request failed,ok=%v,err=%v\n", ok, err)
		return
	}

	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		fmt.Printf("channel send request failed,ok=%v,err=%v\n", ok, err)
		return
	}

	done := make(chan bool, 2)
	go func() {
		defer func() {
			done <- true
		}()

		for {
			//从前端读取数据写给ssh client
			m, p, err := ws.ReadMessage()
			if err != nil {
				fmt.Printf("websocket read message meet err=%v\n", err)
				return
			}

			if m == websocket.TextMessage {
				fmt.Printf("command=%v\n", p)
				if _, err := channel.Write(p); nil != err {
					fmt.Printf("channel.Write meet err=%v\n", err)
					return
				}
			}
		}
	}()

	go func() {
		defer func() {
			done <- true
		}()

		br := bufio.NewReader(channel)
		buf := []byte{}

		t := time.NewTimer(time.Millisecond * 100)
		defer t.Stop()
		r := make(chan rune)

		go func() {
			for {
				//从ssh读取数据，写到r
				x, size, err := br.ReadRune()
				if err != nil {
					fmt.Printf("readRune meet error=%v", err)
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		for {
			select {
			case <-t.C:
				if len(buf) != 0 {
					err = ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						fmt.Printf("writeMessage meet error=%v\n", err)
						return
					}
				}
				t.Reset(time.Millisecond * 100)
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	<-done
}

//-----------get web console url plugin--------------------//
var WebConsoleActions = make(map[string]Action)

func init() {
	WebConsoleActions["get_webconsole_url"] = new(GetWebConsoleUrlAction)
}

type WebConsolePlugin struct {
}

func (plugin *WebConsolePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := WebConsoleActions[actionName]
	if !found {
		return nil, fmt.Errorf("webConsole plugin,action = %s not found", actionName)
	}
	return action, nil
}

type WebConsoleUrlInputs struct {
	Inputs []WebConsoleUrlInput `json:"inputs,omitempty"`
}

type WebConsoleUrlInput struct {
	Guid      string `json:"guid,omitempty"`
	HostIp    string `json:"host_ip,omitempty"`
	ShellPort uint   `json:"shell_port,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Seed      string `json:"seed,omitempty"`
	Password  string `json:"password,omitempty"`
}

type WebConsoleOutputs struct {
	Outputs []WebConsoleOutput `json:"outputs,omitempty"`
}

type WebConsoleOutput struct {
	Guid        string `json:"guid,omitempty"`
	Token       string `json:"token,omitempty"`
	ReDirectUrl string `json:"redirect_url,omitempty"`
}

type GetWebConsoleUrlAction struct {
}

func (action *GetWebConsoleUrlAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs WebConsoleUrlInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *GetWebConsoleUrlAction) CheckParam(input interface{}) error {
	inputs, ok := input.(WebConsoleUrlInputs)
	if !ok {
		return fmt.Errorf("AddUserAction:input type=%T not right", input)
	}
	for _, input := range inputs.Inputs {
		if input.HostIp == "" {
			return errors.New("host_ip is empty")
		}
		if input.Seed == "" {
			return errors.New("seed is empty")
		}
		if input.Password == "" {
			return errors.New("password is empty")
		}
	}
	return nil
}

type SshConnectResult struct {
	Err      error
	Input    *WebConsoleUrlInput
	Password string
}

func trySshConnect(input *WebConsoleUrlInput, password string, ch chan SshConnectResult) {
	result := SshConnectResult{
		Input:    input,
		Password: password,
	}

	sh := &ssh{
		user: input.UserName,
		pwd:  password,
		addr: fmt.Sprintf("%s:%v", input.HostIp, input.ShellPort),
	}

	sh, err := sh.Connect()
	if err != nil {
		result.Err = err
	} else {
		sh.client.Close()
	}

	ch <- result
}

func getSshInfoByTimeStamp(timeStamp string) (*ssh, error) {
	webConsoleTokenMutex.Lock()
	defer webConsoleTokenMutex.Unlock()
	sh, ok := sshTokenMap[timeStamp]
	if !ok {
		return sh, fmt.Errorf("can't found token in map")
	}
	return sh, nil
}

func addSshInfoToMap(token string, password string, input *WebConsoleUrlInput) {
	webConsoleTokenMutex.Lock()
	defer webConsoleTokenMutex.Unlock()
	delKeys := []string{}

	for timestamp, _ := range sshTokenMap {
		k, _ := strconv.ParseInt(timestamp, 10, 64)
		tm := time.Unix(0, k)
		passedTime := time.Since(tm)
		if passedTime.Minutes() > 5 {
			delKeys = append(delKeys, timestamp)
		}
	}

	for _, key := range delKeys {
		delete(sshTokenMap, key)
	}

	sh := &ssh{
		user: input.UserName,
		pwd:  password,
		addr: fmt.Sprintf("%s:%v", input.HostIp, input.ShellPort),
	}

	sshTokenMap[token] = sh
}

func (action *GetWebConsoleUrlAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(WebConsoleUrlInputs)
	chResult := make(chan SshConnectResult)
	outputs := WebConsoleOutputs{}

	for i := 0; i < len(inputs.Inputs); i++ {
		if inputs.Inputs[i].ShellPort == 0 {
			inputs.Inputs[i].ShellPort = WEB_CONSOLE_DEFAULT_PORT
		}
		if inputs.Inputs[i].UserName == "" {
			inputs.Inputs[i].UserName = WEB_CONSOLE_DEFAULT_USER_NAME
		}
	}

	for _, input := range inputs.Inputs {
		md5sum := Md5Encode(input.Guid + input.Seed)
		password, err := AesDecode(md5sum[0:16], input.Password)
		if err != nil {
			return outputs, err
		}

		go trySshConnect(&input, password, chResult)
	}

	for _, _ = range inputs.Inputs {
		result := <-chResult
		if result.Err != nil {
			return outputs, fmt.Errorf("host(%v) ssh connect failed,err=%v", result.Input.HostIp, result.Err)
		}
		token := fmt.Sprintf("%v", time.Now().UnixNano())
		output := WebConsoleOutput{
			Guid:        result.Input.Guid,
			Token:       token,
			ReDirectUrl: "/v1/deploy/webconsole",
		}
		addSshInfoToMap(token, result.Password, result.Input)
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, nil
}
