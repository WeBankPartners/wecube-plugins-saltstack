package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"strconv"
)

var MysqlScriptPluginActions = make(map[string]Action)

func init() {
	MysqlScriptPluginActions["run"] = new(RunMysqlScriptAction)
}

type MysqlScriptPlugin struct {
}

func (plugin *MysqlScriptPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MysqlScriptPluginActions[actionName]
	if !found {
		return nil, fmt.Errorf("mysql script plugin,action = %s not found", actionName)
	}

	return action, nil
}

type RunMysqlScriptInputs struct {
	Inputs []RunMysqlScriptInput `json:"inputs,omitempty"`
}

type RunMysqlScriptInput struct {
	CallBackParameter
	EndPoint     string `json:"endpoint,omitempty"`
	SqlFiles     string `json:"sql_files,omitempty"`
	Guid         string `json:"guid,omitempty"`
	Seed         string `json:"seed,omitempty"`
	Host         string `json:"host,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	DatabaseName string `json:"databaseName,omitempty"`
	Port         string `json:"port,omitempty"`
	// AccessKey string `json:"accessKey,omitempty"`
	// SecretKey string `json:"secretKey,omitempty"`
}

type RunMysqlScriptOutputs struct {
	Outputs []RunMysqlScriptOutput `json:"outputs,omitempty"`
}

type RunMysqlScriptOutput struct {
	CallBackParameter
	Result
	Guid   string `json:"guid,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type RunMysqlScriptAction struct { Language string }

func (action *RunMysqlScriptAction) SetAcceptLanguage(language string) {
	action.Language = language
}

func (action *RunMysqlScriptAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RunMysqlScriptInputs
	if err := UnmarshalJson(param, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *RunMysqlScriptAction) runMysqlScriptCheckParam(input RunMysqlScriptInput) error {
	if input.Host == "" {
		return getParamEmptyError(action.Language, "host")
	}
	if checkIllegalParam(input.Host) {
		return getParamValidateError(action.Language, "host", "Contains illegal character")
	}
	if input.Guid == "" {
		return getParamEmptyError(action.Language, "guid")
	}
	if input.Seed == "" {
		return getParamEmptyError(action.Language, "seed")
	}
	if input.UserName == "" {
		return getParamEmptyError(action.Language, "userName")
	}
	if checkIllegalParam(input.UserName) {
		return getParamValidateError(action.Language, "userName", "Contains illegal character")
	}
	if input.Password == "" {
		return getParamEmptyError(action.Language, "password")
	}
	if checkIllegalParam(input.Password) {
		return getParamValidateError(action.Language, "password", "Contains illegal character")
	}
	if input.EndPoint == "" {
		return getParamEmptyError(action.Language, "endpoint")
	}

	if input.Port == "" {
		input.Port = "3306"
	}else{
		_,err := strconv.Atoi(input.Port)
		if err != nil {
			return getParamValidateError(action.Language, "port", "Port is not num")
		}
	}

	return nil
}

func execSqlScript(hostName string, port string, userName string, password string, databaseName string, fileName string) (string, error) {
	argv := []string{
		"-h" + hostName,
		"-u" + userName,
		"-p" + password,
		"-P" + port,
	}

	if databaseName != "" {
		argv = append(argv, "-D"+databaseName)
	}
	log.Logger.Debug("Exec sql script", log.StringList("args", argv))
	cmd := exec.Command("/usr/bin/mysql", argv...)
	f, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("Exec sql script,open script file fail,%s ", err.Error())
	}

	defer f.Close()
	cmd.Stdin = f

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Logger.Error("Exec sql script", log.String("output", string(out)), log.Error(err))
		return "", fmt.Errorf("Exec sql fail,output=%s,err=%s ", string(out), err.Error())
	}

	return string(out), nil
}

func (action *RunMysqlScriptAction) runMysqlScript(input *RunMysqlScriptInput) (output RunMysqlScriptOutput, err error) {
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
	err = action.runMysqlScriptCheckParam(*input)
	if err != nil {
		return output, err
	}

	var fileNameList []string
	for _,v := range splitWithCustomFlag(input.EndPoint) {
		fileName, err := downloadS3File(v, DefaultS3Key, DefaultS3Password, false, action.Language)
		if err != nil {
			return output, err
		}
		fileNameList = append(fileNameList, fileName)
	}
	//fileName, err := downloadS3File(input.EndPoint, DefaultS3Key, DefaultS3Password, false)
	//if err != nil {
	//	return output, err
	//}

	password, err := AesDePassword(input.Guid, input.Seed, input.Password)
	if err != nil {
		err = getPasswordDecodeError(action.Language, err)
		return output, err
	}

	// new dir to place all *.sql
	Info := strings.Split(fileNameList[0], "/")
	newDir := strings.Join(Info[0:len(Info)-2], "/") + "/sql"
	err = ensureDirExist(newDir)
	if err != nil {
		return output, err
	}

	files := []string{}
	// whether the fileName is *.sql or other
	if !strings.HasSuffix(fileNameList[0], ".sql") {
		if input.SqlFiles == "" {
			err = getParamEmptyError(action.Language, "sql_files")
			return output, err
		}

		if len(fileNameList) > 1 {
			return output,fmt.Errorf("Param endpoint validate fail,endpoint must be one when suffix not like *.sql ")
		}

		// unpack file
		err = deriveUnpackfile(fileNameList[0], newDir, true, action.Language)
		if err != nil {
			return output, err
		}

		// split SqlFiles to *.sql
		//sqlFiles := strings.Split(input.SqlFiles, ",")
		sqlFiles := splitWithCustomFlag(input.SqlFiles)
		for _, file := range sqlFiles {
			sqlFile := newDir + "/" + strings.TrimSpace(file)
			if !fileExist(sqlFile) {
				err = getFileNotExistError(action.Language, sqlFile)
				return output, err
			}
			files = append(files, sqlFile)
		}
	} else {
		for _,v := range fileNameList {
			// move the *.sql to newDir directly
			out, tmpErr := exec.Command("/bin/bash", "-c", fmt.Sprintf("mv -f %s %s", v, newDir)).Output()
			log.Logger.Debug("Run move command", log.String("output", string(out)), log.Error(tmpErr))
			if tmpErr != nil {
				err = fmt.Errorf("Move s3 file:%s to %s error %v ", v, newDir, tmpErr)
				return output, err
			}
			tmpNameList := strings.Split(v, "/")
			sqlFile := newDir + "/" + tmpNameList[len(tmpNameList)-1]
			if !fileExist(sqlFile) {
				err = getFileNotExistError(action.Language, sqlFile)
				return output, err
			}
			files = append(files, sqlFile)
		}
	}

	// run sql scripts foreach
	for _, file := range files {
		_, err = execSqlScript(input.Host, input.Port, input.UserName, password, input.DatabaseName, file)
		if err != nil {
			err = getRunMysqlScriptError(action.Language, file, input.Host, input.DatabaseName, err.Error())
			return output, err
		}
	}

	for _,v := range fileNameList {
		err = os.RemoveAll(v)
		if err != nil {
			return output, err
		}
	}
	err = os.RemoveAll(newDir)
	if err != nil {
		return output, err
	}

	return output, err
}

func (action *RunMysqlScriptAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RunMysqlScriptInputs)
	outputs := RunMysqlScriptOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.runMysqlScript(&input)
		if err != nil {
			log.Logger.Error("Run mysql script action", log.Error(err))
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

func splitWithCustomFlag(input string) []string {
	input = strings.Replace(input, ",", "^^^", -1)
	input = strings.Replace(input, "|", "^^^", -1)
	var output []string
	for _,v := range strings.Split(input, "^^^") {
		if v != "" {
			output = append(output, v)
		}
	}
	return output
}