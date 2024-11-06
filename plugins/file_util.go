package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"path/filepath"
	"strings"
)

type FindResults struct {
	Results []map[string][]string `json:"return,omitempty"`
}

// FindGlobFiles salt '*' file.find path=/var/log name=*yum.log print=path
func FindGlobFiles(dest, name string, target string) (string, error) {
	// parse name to dir and file name

	// name is always relative to dest
	// if !filepath.IsAbs(name) {
	name = filepath.Join(dest, name)
	dirName := filepath.Dir(name)
	fileName := filepath.Base(name)

	// build salt api params
	request := SaltApiRequest{
		Client:   "local",
		Function: "file.find",
		Target:   target,
		Args: []string{
			fmt.Sprintf("path=%s", dirName),
			fmt.Sprintf("name=%s", fileName),
			"print=path",
		},
	}

	output, err := CallSaltApi("https://127.0.0.1:8080", request, "")
	log.Logger.Debug("call salt file.find result",
		log.String("name", name), log.String("target", target), log.String("output", output))
	if err != nil {
		log.Logger.Error("call salt file.find error", log.Error(err), log.JsonObj("request", request))
		return output, err
	}

	// {"return": [{"127.0.0.1": ["/var/log/xx_yum.log", "/var/log/yum.log"]}]}
	var ret FindResults
	if err = json.Unmarshal([]byte(output), &ret); err != nil {
		log.Logger.Error("unmarshal salt file.find error", log.Error(err), log.JsonObj("output", output))
		return output, fmt.Errorf("unmarshal salt file.find result error: %s", err.Error())
	}

	if len(ret.Results) == 0 {
		return output, fmt.Errorf("'%s' not found on %s", name, target)
	}
	return strings.Join(ret.Results[0][target], ","), nil
}
