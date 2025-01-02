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

type ManagedResult struct {
	Name    string      `json:"name"`
	Changes interface{} `json:"changes"`
	Comment string      `json:"comment"`
	Result  bool        `json:"result"`
}

// ManageFileReturn salt file.manage_file return example
//
//	{
//	   "minion": {
//	       "name": "/etc/httpd/test5/httpd.conf",
//	       "changes": {
//	           "diff": "New file",
//	           "mode": "0755"
//	       },
//	       "comment": "File /etc/httpd/test5/httpd.conf updated",
//	       "result": true
//	   }
//	}
type ManageFileReturn map[string]ManagedResult

// SendFile send file by file.manage_file, PlanA chosen by copy
func SendFile(src, dest, owner, target string) (string, error) {
	// PlanA: salt '*' file.manage_file /tmp/test/httpd.conf "" "{}" "salt://base/xxx/api.conf"
	// "{hash_type: 'md5', 'hsum': <md5sum>}" root root "755" "" base "" mkdirs
	// PlanB: salt '*' file.manage_file /tmp/test/httpd.conf "" "{}" ""
	// "{hash_type: 'md5', 'hsum': <md5sum>}" root root "755" "" base "" mkdirs contents="test"
	request := SaltApiRequest{
		Client:   "local",
		Function: "file.manage_file",
		Target:   target,
		Args: []string{
			dest, "", "{}", src, "{hash_type: 'md5', 'hsum': <md5sum>}",
			owner, owner, "755", "", "base", "", "mkdirs",
		},
	}

	output, err := CallSaltApi("https://127.0.0.1:8080", request, "")
	log.Logger.Debug("call salt file.manage_file result",
		log.String("dest", dest), log.String("target", target), log.String("output", output))
	if err != nil {
		log.Logger.Error("call salt file.manage_file error", log.Error(err), log.JsonObj("request", request))
		return output, err
	}

	// parse result only and return md5sum
	var ret ManageFileReturn
	if err = json.Unmarshal([]byte(output), &ret); err != nil {
		log.Logger.Error("unmarshal salt file.manage_file error", log.Error(err), log.JsonObj("output", output))
		return output, fmt.Errorf("unmarshal salt file.manage_file result error: %s", err.Error())
	}
	if !ret[target].Result {
		return output, fmt.Errorf("create file '%s' failed on %s", dest, target)
	}
	md5sum, _ := CallSaltApi("https://127.0.0.1:8080", SaltApiRequest{
		Client:   "local",
		Function: "file.get_hash",
		Target:   target,
		Args:     []string{dest, "md5"},
	}, "")

	return md5sum, nil
}
