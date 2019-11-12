package plugins

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode/utf8"
)

const (
	Zip_gZip                      = "gzip"
	Zip_Deflate                   = "deflate"
	WEB_CONSOLE_DEFAULT_USER_NAME = "root"
	WEB_CONSOLE_DEFAULT_PORT      = 22

	//命令行拦截实现相关
	ENABLE_HIGH_RISK_COMMAND_INTERRUPT = true
	KEY_CR                             = 13
	KEY_CANCEL                         = 3
	KEY_SPACE                          =32
	STATE_WAIT_COMMAND_INPUT           = 0
	STATE_HIGH_RISK_WAIT_CONFIRM       = 1
)

var (
	sshTokenMap          = make(map[string]*ssh)
	webConsoleTokenMutex sync.Mutex
	HighRiskCommands     = []string{
		"rm /* -rf ",
		"rm -rf /* ",
		"rm /* ",
	}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ssh struct {
	user         string
	pwd          string
	addr         string
	client       *gossh.Client
	session      *gossh.Session
	lastCommand  string
	lastInputStr string
	state        int
}

func getHighRiskNotice(command string) string {
	notice := fmt.Sprintf("%c%c%c[0m%c[01;36m(%s) is high risk command,if you want to continue,please press yes.%c[0m", 0x0D, 0x0A, 0x1B, 0x1B, command, 0x1B)
	return notice
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
	Token string `json:"token,omitempty"`
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

func WebConsoleStaticPageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is empty", http.StatusInternalServerError)
		return
	}
	if _, err := getSshInfoByTimeStamp(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//output template
	type WebSocketInfo struct {
		WsAddr string
	}
	wsInfo := WebSocketInfo{
		WsAddr: r.Host + "/v1/deploy/webconsole?token=" + token,
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	rb := bytes.NewBufferString("")
	tmpl, err := template.ParseFiles("/conf/template/console_main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = tmpl.Execute(rb, wsInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//write html to
	Gzip_Html(rb, w, r)
}

func Accept_Encoding(r *http.Request) string {
	ae := r.Header.Get("Accept-Encoding")
	ae = strings.ToLower(ae)
	return ae
}
func Uint32(str string) (uint32, error) {
	v, err := strconv.ParseUint(str, 10, 32)
	return uint32(v), err
}

func Gzip_Html(b io.Reader, w http.ResponseWriter, r *http.Request) {
	ae := Accept_Encoding(r)

	if strings.Contains(ae, Zip_gZip) {
		w.Header().Set("Content-Encoding", Zip_gZip)
		gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer gw.Close()
		b, err := ioutil.ReadAll(b)
		if nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		gw.Write(b)
		return
	} else if strings.Contains(ae, Zip_Deflate) {
		w.Header().Set("Content-Encoding", Zip_Deflate)
		fw, err := flate.NewWriter(w, flate.BestCompression)
		if nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer fw.Close()
		b, err := ioutil.ReadAll(b)
		if nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fw.Write(b)
		return
	} else {
		b, err := ioutil.ReadAll(b)
		if nil != err {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(b))
		return
	}
}

func deleteUnusedSpaces(inputStr string) string {
	result := []byte{}
	lastByte := byte(KEY_SPACE)

	for _, ch := range []byte(inputStr) {
		if ch != KEY_SPACE {
			result = append(result, ch)

		} else {
			if lastByte != KEY_SPACE {
				result = append(result, ch)
			}
		}
		lastByte = ch
	}
	return strings.TrimSpace(string(result))
}

func isHighRiskCommand(inputCommandStr string) bool {
	for _, highRiskCommand := range HighRiskCommands {
		highRiskCmd := deleteUnusedSpaces(highRiskCommand)
		inputCmd := deleteUnusedSpaces(inputCommandStr)
		fmt.Printf("highRiskCmd=%v,inputCmd=%v\n",highRiskCmd,inputCmd)
		if highRiskCmd == inputCmd {
			return true
		}
	}
	return false
}

func highRiskCommandWrite(sh *ssh, p []byte, channel gossh.Channel, r chan rune) error {
	var err error
	writeData := []byte{}
	runes := []rune{}

	if len(p)!=1{
		fmt.Printf("highRiskCommandWrite len(p)=%v\n",len(p))
		return nil 
	}

	if sh.state == STATE_WAIT_COMMAND_INPUT {
		if p[0] == KEY_CR {
			if isHighRiskCommand(sh.lastInputStr) {
				writeData = []byte{KEY_CANCEL}
				notice := getHighRiskNotice(sh.lastInputStr)
				runes = []rune(notice)

				sh.state = STATE_HIGH_RISK_WAIT_CONFIRM
				sh.lastCommand = sh.lastInputStr
				sh.lastInputStr = ""
			} else {
				writeData = p
				sh.lastInputStr = ""
			}
		} else {
			writeData = p
			if p[0] == KEY_CANCEL {
				sh.lastInputStr = ""
			} else {
				sh.lastInputStr += string(p[0])
			}
		}
	} else if sh.state == STATE_HIGH_RISK_WAIT_CONFIRM {
		if p[0] == KEY_CR {
			writeData = []byte{KEY_CANCEL}
			if strings.EqualFold("yes", sh.lastInputStr) {
				writeData = append(writeData, KEY_CR)
				writeData = append(writeData, ([]byte(sh.lastCommand))...)
				writeData = append(writeData, KEY_CR)
			}
			sh.state = STATE_WAIT_COMMAND_INPUT
			sh.lastInputStr = ""
		} else {
			writeData = p
			if p[0] == KEY_CANCEL {
				sh.state = STATE_WAIT_COMMAND_INPUT
				sh.lastInputStr = ""
			} else {
				sh.lastInputStr += string(p[0])
			}
		}
	}

	if len(writeData) > 0 {
		_, err = channel.Write(writeData)
	}
	if len(runes) > 0 {
		time.Sleep(time.Millisecond * 100)
		for _, data := range runes {
			r <- data
		}

	}
	return err
}

func WebConsoleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	chanWebsocketInput := make(chan rune)

	defer func() {
		if err != nil {
			w.Header().Set("content-type", "application/json")
			w.Write(getRunWebConsoleBytes(err))
		}

	}()

	token := r.URL.Query().Get("token")
	colsStr := r.URL.Query().Get("cols")
	rowsStr := r.URL.Query().Get("rows")
	rows, _ := Uint32(rowsStr)
	cols, _ := Uint32(colsStr)

	if token == "" {
		err = errors.New("token is empty")
		return
	}

	sh, err := getSshInfoByTimeStamp(token)
	if err != nil {
		fmt.Printf("getSshbyTimeStamp meet err=%v\n", err)
		return
	}
	sh, err = sh.Connect()
	if err != nil {
		fmt.Printf("ssh connect failed,err=%v\n", err)
		return
	}
	defer func() {
		sh.client.Close()
	}()
	sh.state = STATE_WAIT_COMMAND_INPUT

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrader failed err=%v\n", err)
		return
	}
	defer ws.Close()

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
		Columns:  cols,
		Rows:     rows,
		Width:    cols * 8,
		Height:   rows * 8,
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
				if ENABLE_HIGH_RISK_COMMAND_INTERRUPT {
					if err = highRiskCommandWrite(sh, p, channel, chanWebsocketInput); err != nil {
						fmt.Printf("highRiskCommandWrite meet err=%v\n", err)
						return
					}
					continue
				}
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

		go func() {
			for {
				//从ssh读取数据，写到r
				x, size, err := br.ReadRune()
				if err != nil {
					fmt.Printf("readRune meet error=%v", err)
					return
				}
				if size > 0 {
					chanWebsocketInput <- x
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
			case d := <-chanWebsocketInput:
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
	Method      string `json:"method,omitempty"`
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
		return sh, fmt.Errorf("(can't found token(%v) in map", timeStamp)
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
			ReDirectUrl: "/v1/deploy/webconsoleStaticPage?token=" + token,
			Method:      "GET",
		}
		addSshInfoToMap(token, result.Password, result.Input)
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, nil
}
