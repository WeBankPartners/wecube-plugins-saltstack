package plugins

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"os/exec"
	"github.com/WeBankPartners/wecube-plugins-saltstack/common/log"
)

const (
	UNCOMPRESSED_DIR = "/data/decompressed/"
	UPLOADS3FILE_DIR = "/data/minio/"
)

type decompressFunc func(comporessedFileFullPath string, uncompressDir string) error

func decompressZipFile(comporessedFileFullPath string, uncompressDir string) error {
	r, err := zip.OpenReader(comporessedFileFullPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(uncompressDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func decompressTgzFile(comporessedFileFullPath string, uncompressDir string) error {
	fr, err := os.Open(comporessedFileFullPath)
	if err != nil {
		return err
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if hdr.Typeflag != tar.TypeDir {
			os.MkdirAll(uncompressDir+"/"+path.Dir(hdr.Name), os.ModePerm)

			fw, _ := os.OpenFile(uncompressDir+"/"+hdr.Name, os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getDeCompressorFunc(comporessedFileFullPath string) (decompressFunc, error) {
	fileName := path.Base(comporessedFileFullPath)
	funcMap := map[string]decompressFunc{
		".zip":    decompressZipFile,
		".tgz":    decompressTgzFile,
		".tar.gz": decompressTgzFile,
	}

	for ext, function := range funcMap {
		if strings.HasSuffix(fileName, ext) {
			return function, nil
		}
	}

	return nil, fmt.Errorf("Unsupported compress type,fileName=%s ", fileName)
}

func validateCompressedFile(comporessedFileFullPath string) error {
	fileName := path.Base(comporessedFileFullPath)
	validExts := []string{".zip", ".tgz", ".tar.gz"}

	for _, validExt := range validExts {
		if strings.HasSuffix(fileName, validExt) {
			return nil
		}
	}

	return fmt.Errorf("Unsupported compress type,fileName=%s ", fileName)
}

func getDecompressDirName(comporessedFileName string) string {
	fileName := path.Base(comporessedFileName)
	validExts := []string{".zip", ".tgz", ".tar.gz"}

	for _, validExt := range validExts {
		if strings.HasSuffix(fileName, validExt) {
			return UNCOMPRESSED_DIR + fileName[0:len(fileName)-len(validExt)]
		}
	}

	return ""
}

func getCompressFileSuffix(comporessedFileName string) (string, error) {
	validExts := []string{".zip", ".tgz", ".tar.gz"}
	for _, validExt := range validExts {
		if strings.HasSuffix(comporessedFileName, validExt) {
			return validExt, nil
		}
	}

	return "", fmt.Errorf("%s have invalid compressFileSuffix", comporessedFileName)
}

func decompressFile(comporessedFileFullPath string, decompressDir string) error {
	var err error

	decompressFunc, err := getDeCompressorFunc(comporessedFileFullPath)
	if err != nil {
		return err
	}

	if err = ensureDirExist(decompressDir); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			os.RemoveAll(decompressDir)
		}
	}()

	if err = decompressFunc(comporessedFileFullPath, decompressDir); err != nil {
		return err
	}

	return err
}

func bashDecompressFunc(filePath,distDir string) error {
	var bashScript string
	if !strings.HasSuffix(distDir, "/") {
		distDir += "/"
	}
	if strings.HasSuffix(filePath, ".tar.gz") || strings.HasSuffix(filePath, ".tgz") {
		bashScript = fmt.Sprintf("tar zxf %s -C %s", filePath, distDir)
	}
	if strings.HasSuffix(filePath, ".zip") {
		bashScript = fmt.Sprintf("unzip %s -d %s", filePath, distDir)
	}
	if bashScript == "" {
		return fmt.Errorf("Unsupported compress type,fileName=%s ", filePath)
	}
	output,err := exec.Command("/bin/bash", "-c", bashScript).Output()
	if err != nil {
		log.Logger.Error("Try to bash decompress fail", log.String("cmd", bashScript), log.String("output", string(output)), log.Error(err))
		return fmt.Errorf("Try to bash decompress fail,output=%s,error=%s ", string(output), err.Error())
	}
	return nil
}