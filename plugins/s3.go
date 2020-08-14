package plugins

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"sync"
	"time"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

func uploadS3File(endPoint, accessKey, secretKey, language string) (string, error) {
	//get minio url, package name
	if !strings.Contains(endPoint, "//") {
		return "", getS3UrlValidateError(language, endPoint)
	}
	s := strings.Split(endPoint, "//")
	if len(s) < 2 {
		return "", getS3UrlValidateError(language, endPoint)
	}

	Info := strings.Split(s[1], "/")
	if len(Info) < 3 {
		return "", getS3UrlValidateError(language, endPoint)
	}
	if !strings.Contains(Info[len(Info)-1], ".") {
		return "", getS3UploadError(language, endPoint, fmt.Sprintf("package name %s is unvaliable", Info[len(Info)-1]))
	}

	minioStoragePath := ""
	for i := 1; i < len(Info); i++ {
		minioStoragePath += "/" + Info[i]
	}

	pkgInfo := strings.Split(Info[len(Info)-1], ".")
	err := ensureDirExist(UPLOADS3FILE_DIR)
	if err != nil {
		return "", getS3UploadError(language, endPoint, fmt.Sprintf("create upload path error: %s", err.Error()))
	}

	path := UPLOADS3FILE_DIR + pkgInfo[0]
	//check dir exist,need to replace new file TODO
	_, err = os.Stat(path)
	if err == nil {
		log.Logger.Warn("Upload s3 file,path already exist", log.String("path", path))
		return path, nil
	}

	err = fileReplace(endPoint, accessKey, secretKey)
	if err != nil {
		return "", getS3UploadError(language, endPoint, fmt.Sprintf("Prepare s3 template file error: %s ", err))
	}

	sh := "s3cmd -c /home/app/wecube-plugins-saltstack/minioconf put "
	sh += UPLOADS3FILE_DIR + Info[len(Info)-1] + " s3:/" + minioStoragePath
	cmd := exec.Command("/bin/sh", "-c", sh)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", getS3UploadError(language, endPoint, fmt.Sprintf("exec s3cmd to upload fail,output:%s, error:%s", stderr.String(), err.Error()))
	}
	return path, nil
}

func downloadS3File(endPoint, accessKey, secretKey string,randName bool,language string) (string, error) {
	var tmpName string
	if randName {
		tmpName = getWorkspaceName()
	}
	s := strings.Split(endPoint, "//")
	if len(s) < 2 {
		return "", getS3UrlValidateError(language, endPoint)
	}

	Info := strings.Split(s[1], "/")
	if len(Info) < 3 {
		return "", getS3UrlValidateError(language, endPoint)
	}

	//check dir exist
	ensureDirExist(UPLOADS3FILE_DIR)

	path := UPLOADS3FILE_DIR + tmpName + Info[len(Info)-1]
	_, err := os.Stat(path)
	if err == nil {
		log.Logger.Info("Download s3 file stop,already exists", log.String("path", path))
		return path, nil
	}
	//config s3,need to change different workspace TODO
	err = fileReplace(endPoint, accessKey, secretKey)
	if err != nil {
		return "", getS3DownloadError(language, endPoint, fmt.Sprintf("s3 template config error: %s", err.Error()))
	}

	storagePath := ""
	for i := 1; i < len(Info); i++ {
		storagePath += "/" + Info[i]
	}
	sh := "s3cmd -c /home/app/wecube-plugins-saltstack/minioconf get --force "
	sh += " s3:/" + storagePath + " " + path
	log.Logger.Debug("S3 command", log.String("command", sh))
	cmd := exec.Command("/bin/sh", "-c", sh)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		os.Remove(path)
		tmpErrorMsg := fmt.Sprint(err) + ": " + stderr.String()
		if strings.Contains(tmpErrorMsg, "404") {
			return "", getS3FileNotFoundError(language, storagePath)
		}
		return "", getS3DownloadError(language, endPoint, tmpErrorMsg)
	}
	log.Logger.Debug("Download s3 file result", log.String("output", stderr.String()))
	return path, nil
}

//MinioConf .
type MinioConf struct {
	AccessKey string
	MinioURL  string
	BucketURL string
	SecretKey string
}

func fileReplace(endPoint, accessKey, secretKey string) error {
	funcMap := template.FuncMap{}

	s := strings.Split(endPoint, "//")
	Info := strings.Split(s[1], "/")

	test := MinioConf{
		AccessKey: accessKey,
		MinioURL:  "http://" + Info[0],
		BucketURL: "http://" + Info[0] + "/" + Info[1] + "/",
		SecretKey: secretKey,
	}

	tmpl, err := template.New("s3conf").Funcs(funcMap).ParseFiles("/conf/s3conf")
	if err != nil {
		return fmt.Errorf("parsing error: %s", err)
	}

	f, err := os.OpenFile("/home/app/wecube-plugins-saltstack/minioconf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("open file error: %s", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, test)
	if err != nil {
		return fmt.Errorf("execution error: %s", err)
	}
	return nil
}

//GetVariable .
func GetVariable(filepath string,specialList []string) ([]ConfigKeyInfo, error) {
	_, err := PathExists(filepath)
	if err != nil {
		log.Logger.Error("Get variable error", log.Error(err))
		return nil, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		err = fmt.Errorf("Open file %s error,%s ", filepath, err.Error())
		return nil, err
	}
	defer f.Close()

	variableList := []ConfigKeyInfo{}
	br := bufio.NewReader(f)
	lineNumber := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if len(line) == 0 {
			continue
		}

		flysnowRegexp := regexp.MustCompile(`[^\[]*]`)
		params := flysnowRegexp.FindAllString(string(line), -1)
		if len(params) > 0 {
			var configKey ConfigKeyInfo
			n := strconv.Itoa(lineNumber)

			for _, param := range params {
				if false == strings.HasSuffix(param, "]") {
					continue
				}
				param = param[0 : len(param)-1]

				for _, specialFlag := range specialList {
					if specialFlag == "" {
						continue
					}
					if strings.HasPrefix(param, specialFlag) {
						s := strings.Split(param, specialFlag)
						if s[1] == "" {
							return nil, fmt.Errorf("File %s have unvaliable param %s ", filepath, param)
						}
						if strings.Contains(s[1], " ") {
							continue
						}

						configKey.Line = n
						configKey.Key = s[1]
						variableList = append(variableList, configKey)
					}
				}
			}
		}
		lineNumber++
	}

	return variableList, nil
}

var getNameLock = new(sync.RWMutex)

func getWorkspaceName() (name string) {
	getNameLock.Lock()
	name = fmt.Sprintf("%d-", time.Now().UnixNano())
	getNameLock.Unlock()
	return name
}