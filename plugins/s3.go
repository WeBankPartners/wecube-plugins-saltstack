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

	"github.com/sirupsen/logrus"
)

func uploadS3File(endPoint, accessKey, secretKey string) (string, error) {
	//get minio url, package name
	if !strings.Contains(endPoint, "//") {
		return "", fmt.Errorf("Endpoint is unvaliable, don't have '//' : %s", endPoint)
	}
	s := strings.Split(endPoint, "//")
	if len(s) < 2 {
		return "", fmt.Errorf("endpoint(%s) is not a valid s3 url", endPoint)
	}

	Info := strings.Split(s[1], "/")
	if len(Info) < 3 {
		return "", fmt.Errorf("Endpoint is unvaliable: %s", endPoint)
	}
	if !strings.Contains(Info[len(Info)-1], ".") {
		return "", fmt.Errorf("package name is unvaliable: %s", Info[len(Info)-1])
	}

	minioStoragePath := ""
	for i := 1; i < len(Info); i++ {
		minioStoragePath += "/" + Info[i]
	}

	pkgInfo := strings.Split(Info[len(Info)-1], ".")
	err := ensureDirExist(UPLOADS3FILE_DIR)
	if err != nil {
		return "", fmt.Errorf("create upload path error : %s", err)
	}

	path := UPLOADS3FILE_DIR + pkgInfo[0]
	//check dir exist
	_, err = os.Stat(path)
	if err == nil {
		logrus.Infof("path %v already exist. ", path)
		return path, nil
	}

	err = fileReplace(endPoint, accessKey, secretKey)
	if err != nil {
		return "", fmt.Errorf("template execution error: %s", err)
	}

	sh := "s3cmd -c /home/app/wecube-plugins-saltstack/minioconf put "
	sh += UPLOADS3FILE_DIR + Info[len(Info)-1] + " s3:/" + minioStoragePath
	cmd := exec.Command("/bin/sh", "-c", sh)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return path, nil
}

func downloadS3File(endPoint, accessKey, secretKey string) (string, error) {
	s := strings.Split(endPoint, "//")
	if len(s) < 2 {
		return "", fmt.Errorf("endpoint(%s) is not a valid s3 url", endPoint)
	}

	Info := strings.Split(s[1], "/")
	if len(Info) < 3 {
		return "", fmt.Errorf("endpoint(%s) is not a valid s3 url", endPoint)
	}

	//check dir exist
	err := ensureDirExist(UPLOADS3FILE_DIR)
	if err != nil {
		logrus.Infof("downloadS3File ensureDirExist meet err=%v", err)
		return "", fmt.Errorf("create upload path error : %s", err)
	}

	path := UPLOADS3FILE_DIR + Info[len(Info)-1]
	_, err = os.Stat(path)
	if err == nil {
		return path, nil
	}
	err = fileReplace(endPoint, accessKey, secretKey)
	if err != nil {
		return "", fmt.Errorf("template execution error: %s", err)
	}

	storagePath := ""
	for i := 1; i < len(Info); i++ {
		storagePath += "/" + Info[i]
	}
	sh := "s3cmd -c /home/app/wecube-plugins-saltstack/minioconf get "
	sh += " s3:/" + storagePath + " " + UPLOADS3FILE_DIR + Info[len(Info)-1]

	cmd := exec.Command("/bin/sh", "-c", sh)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("updown file error: " + fmt.Sprint(err) + ": " + stderr.String())
	}
	logrus.Infof("result=%v", stderr.String())
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
func GetVariable(fullpath, filepath string) ([]ConfigKeyInfo, error) {
	fullPath := fullpath + "/" + filepath
	_, err := os.Stat(fullPath)
	if err != nil {
		logrus.Errorf("path %s not exist. ", fullPath)
		return nil, err
	}

	bf, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file %s fail: %s", fullPath, err)
	}
	defer bf.Close()

	variableList := []ConfigKeyInfo{}
	br := bufio.NewReader(bf)
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

				if strings.Contains(param, "@") {
					s := strings.Split(param, "@")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					configKey.Line = n
					configKey.Key = s[1]
					variableList = append(variableList, configKey)
				}
				if strings.Contains(param, "!") {
					s := strings.Split(param, "!")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					configKey.Line = n
					configKey.Key = s[1]
					variableList = append(variableList, configKey)
				}
				if strings.Contains(param, "&") {
					s := strings.Split(param, "&")
					if s[1] == "" {
						return nil, fmt.Errorf("file %s have unvaliable variable %s", filepath, param)
					}
					configKey.Line = n
					configKey.Key = s[1]
					variableList = append(variableList, configKey)
				}
			}
		}
		lineNumber++
	}

	return variableList, nil
}
