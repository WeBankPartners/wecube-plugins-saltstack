package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var PasswordPluginActions = make(map[string]Action)

func init() {
	PasswordPluginActions["encode"] = new(PasswordEncodeAction)
	PasswordPluginActions["decode"] = new(PasswordDecodeAction)
	PasswordPluginActions["sshkeygen"] = new(PasswordSSHKeyGenAction)
	PasswordPluginActions["rsakeygen"] = new(RsaKeyGenAction)
}

type PasswordPlugin struct {
}

func (plugin *PasswordPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := PasswordPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("Password plugin,action = %s not found", actionName)
	}

	return action, nil
}

type PasswordEncodeInputs struct {
	Inputs []PasswordEncodeInput `json:"inputs,omitempty"`
}

type PasswordEncodeInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordEncodeOutputs struct {
	Outputs []PasswordEncodeOutput `json:"outputs,omitempty"`
}

type PasswordEncodeOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordDecodeInputs struct {
	Inputs []PasswordDecodeInput `json:"inputs,omitempty"`
}

type PasswordDecodeInput struct {
	CallBackParameter
	Guid     string `json:"guid,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordDecodeOutputs struct {
	Outputs []PasswordDecodeOutput `json:"outputs,omitempty"`
}

type PasswordDecodeOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	Password string `json:"password,omitempty"`
}

type PasswordSSHKeyGenInputs struct {
	Inputs []*PasswordSSHKeyGenInput `json:"inputs,omitempty"`
}

type PasswordSSHKeyGenInput struct {
	CallBackParameter
	Guid    string `json:"guid,omitempty"`
	Seed    string `json:"seed,omitempty"`
	KeyName string `json:"keyName,omitempty"`
}

type PasswordSSHKeyGenOutputs struct {
	Outputs []*PasswordSSHKeyGenOutput `json:"outputs,omitempty"`
}

type PasswordSSHKeyGenOutput struct {
	CallBackParameter
	Result
	Guid       string `json:"guid,omitempty"`
	KeyName    string `json:"keyName,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	PublicKey  string `json:"publicKey,omitempty"`
}

type RsaKeyGenInputs struct {
	Inputs []*RsaKeyGenInput `json:"inputs,omitempty"`
}

type RsaKeyGenInput struct {
	CallBackParameter
	Guid string `json:"guid,omitempty"`
	Seed string `json:"seed,omitempty"`
}

type RsaKeyGenOutputs struct {
	Outputs []*RsaKeyGenOutput `json:"outputs,omitempty"`
}

type RsaKeyGenOutput struct {
	CallBackParameter
	Result
	Guid       string `json:"guid,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	PublicKey  string `json:"publicKey,omitempty"`
}

type PasswordEncodeAction struct {
	Language string
}

type PasswordDecodeAction struct {
	Language string
}

type PasswordSSHKeyGenAction struct {
	Language string
}

type RsaKeyGenAction struct {
	Language string
}

func (action *PasswordEncodeAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *PasswordEncodeAction) CheckParam(input PasswordEncodeInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}

	return nil
}

func (action *PasswordEncodeAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PasswordEncodeInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *PasswordEncodeAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(PasswordEncodeInputs)
	outputs := PasswordEncodeOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := PasswordEncodeOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		input.Seed = getEncryptSeed(input.Seed)
		encryptPassword, err := AesEnPassword(input.Guid, input.Seed, input.Password, DEFALT_CIPHER)
		if err != nil {
			err = getPasswordEncodeError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		output.Password = encryptPassword
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

func (action *PasswordDecodeAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *PasswordDecodeAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PasswordDecodeInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *PasswordDecodeAction) CheckParam(input PasswordDecodeInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	//if input.Seed == "" {
	//	return getParamEmptyError(action.Language, "seed")
	//}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}

	return nil
}

func (action *PasswordDecodeAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(PasswordDecodeInputs)
	outputs := PasswordDecodeOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := PasswordDecodeOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		input.Seed = getEncryptSeed(input.Seed)
		decodePassword, err := AesDePassword(input.Guid, input.Seed, input.Password)
		if err != nil {
			err = getPasswordDecodeError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		output.Password = decodePassword
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

func (action *PasswordSSHKeyGenAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *PasswordSSHKeyGenAction) CheckParam(input *PasswordSSHKeyGenInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	if input.KeyName == "" {
		return getParamEmptyError(action.Language, "keyName")
	}

	return nil
}

func (action *PasswordSSHKeyGenAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PasswordSSHKeyGenInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *PasswordSSHKeyGenAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(PasswordSSHKeyGenInputs)
	outputs := PasswordSSHKeyGenOutputs{Outputs: []*PasswordSSHKeyGenOutput{}}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := PasswordSSHKeyGenOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		input.Seed = getEncryptSeed(input.Seed)
		privateKey, publicKey, err := action.GenSSHKey(input)
		if err != nil {
			err = getGenSSHKeyError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		output.PrivateKey, err = AesEnPassword(input.Guid, input.Seed, privateKey, DEFALT_CIPHER)
		if err != nil {
			err = getPasswordEncodeError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		output.PublicKey = publicKey
		output.KeyName = input.KeyName
		outputs.Outputs = append(outputs.Outputs, &output)
	}
	return &outputs, finalErr
}

func (action *PasswordSSHKeyGenAction) GenSSHKey(input *PasswordSSHKeyGenInput) (privateKey, publicKey string, err error) {
	workDir := "/tmp/sshgen_" + getRandString()
	if err = os.MkdirAll(workDir, 0755); err != nil {
		err = fmt.Errorf("mkdir for ssh key gen fail,%s ", err.Error())
		return
	}
	defer os.RemoveAll(workDir)
	cmdStr := fmt.Sprintf("ssh-keygen -t rsa -b 2048 -m PEM -N \"\" -q -C \"%s\" -f %s/sshkey", input.KeyName, workDir)
	_, err = exec.Command("/bin/bash", "-c", cmdStr).Output()
	if err != nil {
		err = fmt.Errorf("exec sshkeygen command:%s fail,%s ", cmdStr, err.Error())
		return
	}
	priBytes, readPriErr := os.ReadFile(workDir + "/sshkey")
	if readPriErr != nil {
		err = fmt.Errorf("read private key from %s/sshkey fail,%s ", workDir, readPriErr.Error())
		return
	}
	privateKey = strings.ReplaceAll(string(priBytes), "\n", "\\n")
	pubBytes, readPubErr := os.ReadFile(workDir + "/sshkey.pub")
	if readPubErr != nil {
		err = fmt.Errorf("read public key from %s/sshkey.pub fail,%s ", workDir, readPubErr.Error())
		return
	}
	publicKey = string(pubBytes)
	return
}

func (action *RsaKeyGenAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *RsaKeyGenAction) CheckParam(input *RsaKeyGenInput) error {
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	return nil
}

func (action *RsaKeyGenAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RsaKeyGenInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *RsaKeyGenAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RsaKeyGenInputs)
	outputs := RsaKeyGenOutputs{Outputs: []*RsaKeyGenOutput{}}
	var finalErr error
	for _, input := range inputs.Inputs {
		output := RsaKeyGenOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		if err := action.CheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		input.Seed = getEncryptSeed(input.Seed)
		privateKey, publicKey, err := action.GenSSHKey(input)
		if err != nil {
			err = getGenSSHKeyError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		output.PrivateKey, err = AesEnPassword(input.Guid, input.Seed, privateKey, DEFALT_CIPHER)
		if err != nil {
			err = getPasswordEncodeError(action.Language, err)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, &output)
			continue
		}
		output.PublicKey = publicKey
		outputs.Outputs = append(outputs.Outputs, &output)
	}
	return &outputs, finalErr
}

func (action *RsaKeyGenAction) GenSSHKey(input *RsaKeyGenInput) (privateKey, publicKey string, err error) {
	workDir := "/tmp/rsagen_" + getRandString()
	if err = os.MkdirAll(workDir, 0755); err != nil {
		err = fmt.Errorf("mkdir for ssh key gen fail,%s ", err.Error())
		return
	}
	defer os.RemoveAll(workDir)
	_, err = runBashScript("/home/app/wecube-plugins-saltstack/scripts/rsa_gen.sh", []string{workDir, "app"})
	if err != nil {
		err = fmt.Errorf("exec rsakeygen fail,%s ", err.Error())
		return
	}
	priBytes, readPriErr := os.ReadFile(workDir + "/app.pem")
	if readPriErr != nil {
		err = fmt.Errorf("read private key from %s/app.pem fail,%s ", workDir, readPriErr.Error())
		return
	}
	privateKey = formatRsaKeyToString(priBytes)
	pubBytes, readPubErr := os.ReadFile(workDir + "/app_pub.pem")
	if readPubErr != nil {
		err = fmt.Errorf("read public key from %s/app_pub.pem fail,%s ", workDir, readPubErr.Error())
		return
	}
	publicKey = formatRsaKeyToString(pubBytes)
	return
}

func formatRsaKeyToString(inputBytes []byte) (output string) {
	for _, v := range strings.Split(string(inputBytes), "\n") {
		if strings.HasPrefix(v, "-----") || v == "" {
			continue
		}
		output = output + v
	}
	return
}
