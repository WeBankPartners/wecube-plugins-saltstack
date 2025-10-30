package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
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

// 文件缓存服务
type DownloadFileParam struct {
	Endpoint         string
	AccessKey        string
	SecretKey        string
	RequestLanguage  string
	RandName         bool
	Md5              string
	WithCopyToMaster bool
}

type FileCacheObj struct {
	Md5 string
	// Status         int // 0->notStart 1->downloading 2->done 3->error
	FilePath       string
	LastUsedTime   int64
	SaltMasterPath string
	ErrorMsg       string
	Lock           *sync.RWMutex
}

func (f *FileCacheObj) UpdateUsedTime() (localFilePath string, err error) {
	f.Lock.RLock()
	f.LastUsedTime = time.Now().Unix()
	localFilePath = f.FilePath
	if f.ErrorMsg != "" {
		err = fmt.Errorf(f.ErrorMsg)
	}
	f.Lock.RUnlock()
	return
}

func (f *FileCacheObj) Expired() (ok bool) {
	f.Lock.RLock()
	if (time.Now().Unix() - f.LastUsedTime) > 3600 {
		ok = true
	}
	f.Lock.RUnlock()
	return
}

func (f *FileCacheObj) GetMasterBasepPath() (saltMasterPath string) {
	f.Lock.RLock()
	saltMasterPath = f.SaltMasterPath
	f.Lock.RUnlock()
	return
}

var (
	GlobalFileCacheMap = new(sync.Map) // key->fileMd5 value->FileCacheObj
)

func GetGlobalCacheFile(param *DownloadFileParam) (localFilePath, saltMasterPath string, err error) {
	fileCache := &FileCacheObj{Md5: param.Md5, Lock: new(sync.RWMutex)}
	existFileCache, loadOk := GlobalFileCacheMap.LoadOrStore(param.Md5, fileCache)
	if loadOk {
		fileCache = existFileCache.(*FileCacheObj)
		_, checkFileErr := os.Stat(fileCache.FilePath)
		if os.IsNotExist(checkFileErr) {
			GlobalFileCacheMap.Delete(param.Md5)
			log.Logger.Warn("GetGlobalCacheFile check cache file not exist", log.String("md5", param.Md5), log.String("localFilePath", localFilePath))
		} else {
			if localFilePath, err = fileCache.UpdateUsedTime(); err != nil {
				err = fmt.Errorf("get cache file fail,%s ", err.Error())
				return
			}
			if param.WithCopyToMaster {
				saltMasterPath = fileCache.GetMasterBasepPath()
				if saltMasterPath == "" {
					savePath, tmpErr := saveFileToSaltMasterBaseDir(localFilePath)
					if tmpErr != nil {
						err = fmt.Errorf("save cache file to master base dir fail,%s ", tmpErr.Error())
					} else {
						fileCache.SaltMasterPath = savePath
						saltMasterPath = savePath
					}
				}
			}
			log.Logger.Info("GetGlobalCacheFile get cache", log.String("md5", param.Md5), log.String("localFilePath", localFilePath), log.String("saltMasterPath", saltMasterPath))
			return
		}
	}
	fileCache.Lock.Lock()
	localFilePath, downloadErr := downloadS3File(param.Endpoint, param.AccessKey, param.SecretKey, param.RandName, param.RequestLanguage)
	if downloadErr != nil {
		err = fmt.Errorf("download cache file fail,%s ", downloadErr.Error())
		fileCache.ErrorMsg = err.Error()
	} else {
		localFileMd5 := countMd5WithCmd(localFilePath)
		if localFileMd5 != param.Md5 {
			err = fmt.Errorf("download cache file fail with md5 check,expectMd5:%s realMd5:%s ", param.Md5, localFileMd5)
			os.Remove(localFilePath)
			fileCache.ErrorMsg = err.Error()
		} else {
			fileCache.FilePath = localFilePath
			fileCache.LastUsedTime = time.Now().Unix()
		}
	}
	if param.WithCopyToMaster && err == nil {
		savePath, tmpErr := saveFileToSaltMasterBaseDir(localFilePath)
		if tmpErr != nil {
			err = fmt.Errorf("save cache file to master base dir fail,%s ", tmpErr.Error())
		} else {
			fileCache.SaltMasterPath = savePath
			saltMasterPath = savePath
		}
	}
	fileCache.Lock.Unlock()
	if err != nil {
		GlobalFileCacheMap.Delete(param.Md5)
		log.Logger.Error("GetGlobalCacheFile error", log.String("md5", param.Md5), log.String("endpoint", param.Endpoint), log.Error(err))
	} else {
		log.Logger.Info("GetGlobalCacheFile set cache", log.String("md5", param.Md5), log.String("localFilePath", localFilePath), log.String("saltMasterPath", saltMasterPath))
	}
	return
}

func cleanCacheFiles() {
	cleanKeyList := []string{}
	GlobalFileCacheMap.Range(func(key, value any) bool {
		fileCache := value.(*FileCacheObj)
		if fileCache.Expired() {
			cleanKeyList = append(cleanKeyList, key.(string))
			os.Remove(fileCache.FilePath)
			if fileCache.SaltMasterPath != "" {
				os.Remove(fileCache.SaltMasterPath)
			}
		}
		return true
	})
	for _, v := range cleanKeyList {
		GlobalFileCacheMap.Delete(v)
		log.Logger.Info("GetGlobalCacheFile clean cache", log.String("md5", v))
	}
}
