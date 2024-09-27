package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

// MountInfo object from salt disk.usage --out json
type MountInfo struct {
	Filesystem string `json:"filesystem"`
	KBlocks    string `json:"1K-blocks"`
	Used       string `json:"used"`
	Available  string `json:"available"`
	Capacity   string `json:"capacity"`
}

// MinionInfo object from salt grains.items --out json
type MinionInfo struct {
	Cwd   string `json:"cwd"`
	IpGw  bool   `json:"ip_gw"`
	Ip4Gw string `json:"ip4_gw"`
	Ip6Gw bool   `json:"ip6_gw"`
	Dns   struct {
		Nameservers    []string `json:"nameservers"`
		Ip4Nameservers []string `json:"ip4_nameservers"`
		Ip6Nameservers []string `json:"ip6_nameservers"`
		Sortlist       []string `json:"sortlist"`
		Domain         string   `json:"domain"`
		Search         []string `json:"search"`
		Options        []string `json:"options"`
	} `json:"dns"`
	Fqdns            []string            `json:"fqdns"`
	MachineId        string              `json:"machine_id"`
	Master           string              `json:"master"`
	ServerId         int                 `json:"server_id"`
	Localhost        string              `json:"localhost"`
	Fqdn             string              `json:"fqdn"`
	Host             string              `json:"host"`
	Domain           string              `json:"domain"`
	HwaddrInterfaces map[string]string   `json:"hwaddr_interfaces"`
	Id               string              `json:"id"`
	Ip4Interfaces    map[string][]string `json:"ip4_interfaces"`
	Ip6Interfaces    map[string][]string `json:"ip6_interfaces"`
	Ipv4             []string            `json:"ipv4"`
	Ipv6             []string            `json:"ipv6"`
	FqdnIp4          []string            `json:"fqdn_ip4"`
	FqdnIp6          []string            `json:"fqdn_ip6"`
	IpInterfaces     map[string][]string `json:"ip_interfaces"`
	Kernelparams     [][]string          `json:"kernelparams"`
	LocaleInfo       struct {
		Defaultlanguage  string `json:"defaultlanguage"`
		Defaultencoding  string `json:"defaultencoding"`
		Detectedencoding string `json:"detectedencoding"`
		Timezone         string `json:"timezone"`
	} `json:"locale_info"`
	NumGpus int `json:"num_gpus"`
	Gpus    []struct {
		Vendor string `json:"vendor"`
		Model  string `json:"model"`
	} `json:"gpus"`
	Kernel        string `json:"kernel"`
	Nodename      string `json:"nodename"`
	Kernelrelease string `json:"kernelrelease"`
	Kernelversion string `json:"kernelversion"`
	Cpuarch       string `json:"cpuarch"`
	Selinux       struct {
		Enabled  bool   `json:"enabled"`
		Enforced string `json:"enforced"`
	} `json:"selinux"`
	Systemd struct {
		Version  string `json:"version"`
		Features string `json:"features"`
	} `json:"systemd"`
	Init               string   `json:"init"`
	LsbDistribId       string   `json:"lsb_distrib_id"`
	LsbDistribCodename string   `json:"lsb_distrib_codename"`
	Osfullname         string   `json:"osfullname"`
	Osrelease          string   `json:"osrelease"`
	Oscodename         string   `json:"oscodename"`
	Os                 string   `json:"os"`
	NumCpus            int      `json:"num_cpus"`
	CpuModel           string   `json:"cpu_model"`
	CpuFlags           []string `json:"cpu_flags"`
	OsFamily           string   `json:"os_family"`
	Osarch             string   `json:"osarch"`
	MemTotal           int      `json:"mem_total"`
	SwapTotal          int      `json:"swap_total"`
	Biosversion        string   `json:"biosversion"`
	Productname        string   `json:"productname"`
	Manufacturer       string   `json:"manufacturer"`
	Biosreleasedate    string   `json:"biosreleasedate"`
	Uuid               string   `json:"uuid"`
	Serialnumber       string   `json:"serialnumber"`
	Virtual            string   `json:"virtual"`
	Ps                 string   `json:"ps"`
	OsreleaseInfo      []int    `json:"osrelease_info"`
	Osmajorrelease     int      `json:"osmajorrelease"`
	Osfinger           string   `json:"osfinger"`
	Path               string   `json:"path"`
	Systempath         []string `json:"systempath"`
	Pythonexecutable   string   `json:"pythonexecutable"`
	Pythonpath         []string `json:"pythonpath"`
	Pythonversion      []string `json:"pythonversion"`
	Saltpath           string   `json:"saltpath"`
	Saltversion        string   `json:"saltversion"`
	Saltversioninfo    []int    `json:"saltversioninfo"`
	Zmqversion         string   `json:"zmqversion"`
	Disks              []string `json:"disks"`
	Ssds               []string `json:"ssds"`
	Shell              string   `json:"shell"`
	Transactional      bool     `json:"transactional"`
	Efi                bool     `json:"efi"`
	EfiSecureBoot      bool     `json:"efi-secure-boot"`
	Lvm                struct {
		VolGroup00 []string `json:"VolGroup00"`
	} `json:"lvm"`
	Mdadm           []string `json:"mdadm"`
	Username        string   `json:"username"`
	Groupname       string   `json:"groupname"`
	Pid             int      `json:"pid"`
	Gid             int      `json:"gid"`
	Uid             int      `json:"uid"`
	ZfsSupport      bool     `json:"zfs_support"`
	ZfsFeatureFlags bool     `json:"zfs_feature_flags"`
}

// HostInfo final detail info of host
type HostInfo struct {
	MinionInfo
	DiskMounts MountInfo `json:"disk_mounts"`
	DiskTotal  int       `json:"disk_total"`
}

// MinionDetailResults salt-api return result
type MinionDetailResults struct {
	Results []map[string]MinionInfo `json:"return,omitempty"`
}

// DiskUsageResults salt-api return result
type DiskUsageResults struct {
	Results []map[string]MountInfo `json:"return,omitempty"`
}

var HostCollectorActions = make(map[string]Action)

func init() {
	HostCollectorActions["query-target"] = new(HostCollectorAction)
}

type HostCollectorPlugin struct {
}

func (plugin *HostCollectorPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := HostCollectorActions[actionName]
	if !found {
		return nil, fmt.Errorf("HostCollector plugin,action = %s not found", actionName)
	}

	return action, nil
}

type HostCollectorInputs struct {
	Inputs []HostCollectorInput `json:"inputs,omitempty"`
}

type HostCollectorInput struct {
	CallBackParameter
	Guid   string `json:"guid,omitempty"`
	Target string `json:"target,omitempty"`
}

type HostCollectorOutputs struct {
	Outputs []HostCollectorOutput `json:"outputs,omitempty"`
}

type HostCollectorOutput struct {
	CallBackParameter
	Result
	Guid   string   `json:"guid,omitempty"`
	Detail HostInfo `json:"detail,omitempty"`
}

type HostCollectorAction struct {
	Language string
}

func (action *HostCollectorAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *HostCollectorAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs HostCollectorInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *HostCollectorAction) Do(input interface{}) (interface{}, error) {
	files, _ := input.(HostCollectorInputs)
	outputs := HostCollectorOutputs{}
	var finalErr error
	for _, file := range files.Inputs {
		fileOutput, err := action.collectHostInfo(&file)
		if err != nil {
			log.Logger.Error("Host collector action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, fileOutput)
	}

	return &outputs, finalErr
}

func (action *HostCollectorAction) CheckParam(input HostCollectorInput) error {
	if input.Target == "" {
		return getParamEmptyError(action.Language, "target")
	}
	if checkIllegalParam(input.Target) {
		return getParamValidateError(action.Language, "target", "Contains illegal character")
	}

	return nil
}
func (action *HostCollectorAction) collectHostInfo(input *HostCollectorInput) (output HostCollectorOutput, err error) {
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

	// check input params
	if err = action.CheckParam(*input); err != nil {
		return output, err
	}

	// salt target disk.usage by salt-api
	var result string
	var diskResults DiskUsageResults
	if result, err = CallSaltApi("https://127.0.0.1:8080", SaltApiRequest{
		Client:   "local",
		Function: "disk.usage",
		Target:   input.Target,
	}, action.Language); err != nil {
		return output, err
	} else {
		if err = json.Unmarshal([]byte(result), &diskResults); err != nil {
			return output, fmt.Errorf("unmarshal disk.usage result error: %s", err.Error())
		}
	}

	// salt target grains.items by salt-api
	var minionResults MinionDetailResults
	minionUrl := fmt.Sprintf("https://127.0.0.1:8080/minions/%s/", input.Target)
	if result, err = CallSaltApi(minionUrl, SaltApiRequest{}, action.Language); err != nil {
		return output, err
	} else {
		if err = json.Unmarshal([]byte(result), &minionResults); err != nil {
			return output, fmt.Errorf("unmarshal disk.usage result error: %s", err.Error())
		}
	}

	// merge to host info
	output.Detail = HostInfo{
		MinionInfo: minionResults.Results[0][input.Target],
		DiskMounts: diskResults.Results[0][input.Target],
	}

	return output, err
}
