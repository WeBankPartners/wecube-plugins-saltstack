package plugins

import (
	"os/exec"
	"strings"
	"time"
	"io/ioutil"
	"os"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

func StartCleanInterval()  {
	cmdOut,_ := exec.Command("/bin/sh", "-c", "date|awk '{print $5}'").Output()
	cmdOutString := strings.TrimSpace(string(cmdOut))
	localTimeZone := "CST"
	log.Logger.Info("timezone", log.String("out", cmdOutString))
	if cmdOutString != "" && len(cmdOutString) <= 5 {
		localTimeZone = cmdOutString
	}
	intervalSecond := 86400
	timeStartValue,_ := time.Parse("2006-01-02 15:04:05 MST", fmt.Sprintf("%s 00:00:00 "+localTimeZone, time.Now().Format("2006-01-02")))
	time.Sleep(time.Duration(timeStartValue.Unix()+86400-time.Now().Unix())*time.Second)
	t := time.NewTicker(time.Duration(intervalSecond)*time.Second).C
	for {
		go cleanLocalPackage(UNCOMPRESSED_DIR)
		go cleanLocalPackage(UPLOADS3FILE_DIR)
		<- t
	}
}

func cleanLocalPackage(dirPath string)  {
	log.Logger.Info("Start clean job", log.String("dir", dirPath))
	files,err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Logger.Error("Clean local package job fail", log.String("dir", dirPath), log.Error(err))
		return
	}
	minUnixTime := time.Now().Unix() - 86400
	for _,f := range files {
		if f.ModTime().Unix() <= minUnixTime {
			log.Logger.Info("Start to clean package", log.String("dir", dirPath), log.String("name", f.Name()))
			err = os.RemoveAll(fmt.Sprintf("%s/%s", dirPath, f.Name()))
			if err != nil {
				log.Logger.Error("Remove package fail", log.String("dir", dirPath), log.String("name", f.Name()), log.Error(err))
			}
		}
	}
}
