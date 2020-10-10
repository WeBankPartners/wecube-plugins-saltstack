package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"os/exec"
	"strings"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

type FileNode struct {
	Name      string      `json:"name"`
	Path      string      `json:"-"`
	FileNodes []*FileNode `json:"-"`
	IsDir     bool        `json:"isDir"`
	Md5       string      `json:"md5"`
}

func listCurrentDirectory(dirname string) ([]FileNode, error) {
	fileNodes := []FileNode{}

	if err := isDirExist(dirname); err != nil {
		return fileNodes, fmt.Errorf("exist check error:%s ", err.Error())
	}

	files := listFiles(dirname)
	for _, filename := range files {
		fileNode := FileNode{
			Name:  filename,
			IsDir: false,
			Md5:   "",
		}
		fpath := filepath.Join(dirname, filename)
		fio, _ := os.Lstat(fpath)
		if fio.IsDir() {
			fileNode.IsDir = true
		}else{
			fileNode.Md5 = countMd5WithCmd(fpath)
		}
		fileNodes = append(fileNodes, fileNode)
	}

	return fileNodes, nil
}

func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)
	names, _ := f.Readdirnames(-1)
	f.Close()

	sort.Strings(names)

	return names
}

func walk(path string, info os.FileInfo, node *FileNode) {
	files := listFiles(path)

	for _, filename := range files {
		fpath := filepath.Join(path, filename)
		fio, _ := os.Lstat(fpath)

		child := FileNode{
			Name:      filename,
			Path:      fpath,
			FileNodes: []*FileNode{},
			IsDir:     fio.IsDir(),
		}

		node.FileNodes = append(node.FileNodes, &child)

		if fio.IsDir() {
			walk(fpath, fio, &child)
		}
	}

	return
}

func isDirExist(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dir(%v) not exist", dir)
		} else {
			return err
		}
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is exist,but not dir", dir)
	}

	return nil
}

func getDirTree(dir string) ([]*FileNode, error) {
	emptyFileNode := []*FileNode{}
	if err := isDirExist(dir); err != nil {
		return emptyFileNode, err
	}
	root := FileNode{"projects", dir, []*FileNode{}, true, ""}
	fileInfo, _ := os.Lstat(dir)
	walk(dir, fileInfo, &root)

	return root.FileNodes, nil
}

func ensureDirExist(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, os.FileMode(0755))
		} else {
			return err
		}
	}
	if !f.IsDir() {
		// exited dir
		return fmt.Errorf("path %s is exist,but not dir", dir)
	}

	return nil
}

func countMd5WithCmd(filePath string) string {
	b,err := exec.Command("bash", "-c", fmt.Sprintf("md5sum %s", filePath)).Output()
	if err != nil {
		log.Logger.Error("count md5 value with command fail", log.String("path", filePath), log.Error(err))
		return ""
	}
	md5Value := strings.Split(string(b), " ")[0]
	return md5Value
}