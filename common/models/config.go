package models

import (
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
	"strings"
	"fmt"
)

type HttpConfig struct {
	Port  string  `json:"port"`
	Token  string  `json:"token"`
}

type LogConfig struct {
	Level   string  `json:"level"`
	File    string  `json:"file"`
	ArchiveMaxSize int `json:"archive_max_size"`
	ArchiveMaxBackup int `json:"archive_max_backup"`
	ArchiveMaxDay int `json:"archive_max_day"`
	Compress  bool  `json:"compress"`
}

type GlobalConfig struct {
	Http  HttpConfig  `json:"http"`
	Log   LogConfig     `json:"log"`
	DefaultLanguage  string  `json:"default_language"`
	InstallMinionTimeout int `json:"install_minion_timeout"`
	ExecRemoteCommandTimeout int `json:"exec_remote_command_timeout"`
}

var (
	Config  *GlobalConfig
)

func InitConfig(cfg string) error {
	if cfg == "" {
		log.Println("use -c to specify configuration file")
		return fmt.Errorf("config file empty,use -c to specify config file")
	}
	_, err := os.Stat(cfg)
	if os.IsExist(err) {
		return fmt.Errorf("config file not found")
	}
	b,err := ioutil.ReadFile(cfg)
	if err != nil {
		return fmt.Errorf("read file %s error %v", cfg, err)
	}
	configContent := strings.TrimSpace(string(b))
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s error %v", cfg, err)
	}
	Config = &c
	return nil
}