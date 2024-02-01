package plugins

import (
	"io/ioutil"
	"encoding/json"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/models"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
	"strings"
	"fmt"
)

var errorMessageList []*models.ErrorMessageObj

func InitErrorMessageList()  {
	fs,err := ioutil.ReadDir("./conf/i18n")
	if err != nil {
		log.Logger.Error("Init errorMessage fail", log.Error(err))
		return
	}
	if len(fs) == 0 {
		log.Logger.Error("Init errorMessage fail, conf/i18n is empty dir")
		return
	}
	for _,v := range fs {
		tmpFileBytes,tmpErr := ioutil.ReadFile("./conf/i18n/"+v.Name())
		if tmpErr != nil {
			log.Logger.Error("Init errorMessage,read " + v.Name() + " fail", log.Error(tmpErr))
			continue
		}
		var tmpErrorMessageObj models.ErrorMessageObj
		tmpErr = json.Unmarshal(tmpFileBytes, &tmpErrorMessageObj)
		if err != nil {
			log.Logger.Error("Init errorMessage,unmarshal file " + v.Name() + " fail", log.Error(tmpErr))
			continue
		}
		tmpErrorMessageObj.Language = strings.Replace(v.Name(), ".json", "", -1)
		errorMessageList = append(errorMessageList, &tmpErrorMessageObj)
	}
	if len(errorMessageList) == 0 {
		log.Logger.Error("Init errorMessage fail, errorMessageList is empty")
	}else{
		log.Logger.Info("Init errorMessage success")
	}
}

func getMessageMap(acceptLanguage string) *models.ErrorMessageObj {
	if len(errorMessageList) == 0 {
		return &models.ErrorMessageObj{}
	}
	if acceptLanguage != "" {
		acceptLanguage = strings.Replace(acceptLanguage, ";", ",", -1)
		for _, v := range strings.Split(acceptLanguage, ",") {
			if strings.HasPrefix(v, "q=") {
				continue
			}
			lowerV := strings.ToLower(v)
			for _, vv := range errorMessageList {
				if vv.Language == lowerV {
					return vv
				}
			}
		}
	}
	for _,v := range errorMessageList {
		if v.Language == models.Config.DefaultLanguage {
			return v
		}
	}
	return errorMessageList[0]
}

func getParamEmptyError(language,paramName string) error {
	return fmt.Errorf(getMessageMap(language).ParamEmptyError, paramName)
}

func getParamValidateError(language,paramName,message string) error {
	return fmt.Errorf(getMessageMap(language).ParamValidateError, paramName, message)
}

func getSysParamEmptyError(language,paramName string) error {
	return fmt.Errorf(getMessageMap(language).SysParamEmptyError, paramName)
}

func getPasswordDecodeError(language string,err error) error {
	return fmt.Errorf(getMessageMap(language).PasswordDecodeError, err.Error())
}

func getPasswordEncodeError(language string,err error) error {
	return fmt.Errorf(getMessageMap(language).PasswordEncodeError, err.Error())
}

func getRemoteCommandError(language,ip,output string,err error) error {
	return fmt.Errorf(getMessageMap(language).ExecRemoteCommandError, ip, output, err.Error())
}

func getInstallMinionError(language,ip,output string) error {
	return fmt.Errorf(getMessageMap(language).InstallMinionError, ip, output)
}

func getUninstallMinionError(language,ip,output string,err error) error {
	return fmt.Errorf(getMessageMap(language).UninstallMinionError, ip, output, err.Error())
}

func getS3UrlValidateError(language,url string) error {
	return fmt.Errorf(getMessageMap(language).S3UrlValidateError, url)
}

func getS3FileNotFoundError(language,file string) error {
	return fmt.Errorf(getMessageMap(language).S3FileEmptyError, file)
}

func getS3DownloadError(language,file,output string) error {
	return fmt.Errorf(getMessageMap(language).S3DownloadError, file, output)
}

func getS3UploadError(language,file,output string) error {
	return fmt.Errorf(getMessageMap(language).S3UploadError, file, output)
}

func getSaltApiTargetError(language,target string) error {
	return fmt.Errorf(getMessageMap(language).SaltApiTargetError, target)
}

func getSaltApiConnectError(language,target string) error {
	return fmt.Errorf(getMessageMap(language).SaltApiConnectError, target)
}

func getDecompressSuffixError(language,file string) error {
	return fmt.Errorf(getMessageMap(language).DecompressSuffixError, file)
}

func getUnpackFileError(language,file string,err error) error {
	return fmt.Errorf(getMessageMap(language).UnpackFileError, file, err.Error())
}

func getMysqlConnectError(language string,err error) error {
	return fmt.Errorf(getMessageMap(language).MysqlConnectError, err.Error())
}

func getAddMysqlDatabaseError(language,message string) error {
	return fmt.Errorf(getMessageMap(language).AddMysqlDatabaseError, message)
}

func getDeleteMysqlDatabaseError(language,message string) error {
	return fmt.Errorf(getMessageMap(language).DeleteMysqlDatabaseError, message)
}

func getRunMysqlCommnandError(language,command,message string) error {
	return fmt.Errorf(getMessageMap(language).RunMysqlCommandError, command, message)
}

func getFileNotExistError(language,file string) error {
	return fmt.Errorf(getMessageMap(language).FileNotExistError, file)
}

func getRunMysqlScriptError(language,file,host,database,message string) error {
	return fmt.Errorf(getMessageMap(language).RunMysqlScriptError, file, host, database, message)
}

func getMysqlCreateUserError(language,user,message string) error {
	return fmt.Errorf(getMessageMap(language).MysqlCreateUserError, user, message)
}

func getRunRemoteScriptError(language,target,output string, err error) error {
	return fmt.Errorf(getMessageMap(language).RunRemoteScriptError, target, output, err.Error())
}