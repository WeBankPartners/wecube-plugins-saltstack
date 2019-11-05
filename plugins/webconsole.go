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
	"time"
	"unicode/utf8"
)

const (
	WEB_CONSOLE_DEFAULT_USER_NAME = "root"
	WEB_CONSOLE_DEFAULT_PORT      = 22
	WEB_CONSOLE_DEFAULT_COLS      = 800
	WEB_CONSOLE_DEFUALT_ROWS      = 600
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
	Guid      string `json:"guid,omitempty"`
	HostIp    string `json:"host_ip,omitempty"`
	ShellPort uint   `json:"shell_port,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Seed      string `json:"seed,omitempty"`
	Password  string `json:"password,omitempty"`
	Rows      uint32 `json:"rows,omitempty"`
	Columns   uint32 `json:"columns,omitempty"`
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

func checkWebConsoleParam(param *RunWebConsoleParam) error {
	if param.HostIp == "" {
		return errors.New("host_ip is empty")
	}
	if param.Guid == "" {
		return errors.New("guid is empty")
	}

	if param.ShellPort == 0 {
		param.ShellPort = WEB_CONSOLE_DEFAULT_PORT
	}

	if param.UserName == "" {
		param.UserName = WEB_CONSOLE_DEFAULT_USER_NAME
	}

	if param.Seed == "" {
		return errors.New("seed is empty")
	}

	if param.Password == "" {
		return errors.New("password is empty")
	}
	if param.Rows == 0 {
		param.Rows = WEB_CONSOLE_DEFUALT_ROWS
	}
	if param.Columns == 0 {
		param.Columns = WEB_CONSOLE_DEFAULT_COLS
	}

	return nil
}

func WebConsoleHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var runWebConsoleParam RunWebConsoleParam
	var password string

	defer func() {
		if err != nil {
			w.Header().Set("content-type", "application/json")
			w.Write(getRunWebConsoleBytes(err))
		}

	}()

	if err = UnmarshalJson(r.Body, &runWebConsoleParam); err != nil {
		return
	}
	if err = checkWebConsoleParam(&runWebConsoleParam); err != nil {
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrader failed err=%v\n", err)
		return
	}
	defer ws.Close()

	md5sum := Md5Encode(runWebConsoleParam.Guid + runWebConsoleParam.Seed)
	password, err = AesDecode(md5sum[0:16], runWebConsoleParam.Password)
	if err != nil {
		return
	}

	sh := &ssh{
		user: runWebConsoleParam.UserName,
		pwd:  password,
		addr: fmt.Sprintf("%s:%v", runWebConsoleParam.HostIp, runWebConsoleParam.ShellPort),
	}

	sh, err = sh.Connect()
	if err != nil {
		fmt.Printf("ssh connect failed,err=%v\n", err)
		return
	}

	channel, incomingRequests, err := sh.client.Conn.OpenChannel("session", nil)
	if err != nil {
		fmt.Printf("ssh open channel failed,err=%v\n", err)
		return
	}

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
		Columns:  runWebConsoleParam.Columns,
		Rows:     runWebConsoleParam.Rows,
		Width:    runWebConsoleParam.Columns * 8,
		Height:   runWebConsoleParam.Rows * 8,
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
