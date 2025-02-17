package plugins

import (
	"fmt"
	"os"
	"testing"
)

func TestReplaceFileVar(t *testing.T) {
	InitEnvParam()
	keyMap := map[string]string{
		"Key":       "my_key",
		"AccessKey": "my_access_key",
		"APP":       "wecube",
		"Standard":  "standard1",
		"SecretKey": "my_secret_key",
		"Allow":     "yes",
	}

	dir, err := os.Getwd()
	fmt.Println(dir)

	filepath := "/Users/tylertang/Desktop/go_test/config2.conf"

	keyInfo, err := GetVariable(filepath, DefaultSpecialReplaceList, false)
	if err != nil {
		fmt.Printf("GetVariable err=%v\n", err)
	}

	fmt.Println("keyINfo", keyInfo)

	priKey := ""
	pubKey := ""
	decompressDirName := ""
	specialReplaceList := []string{}
	prefix := []string{}
	if err != nil {
		err = replaceFileVar(keyMap, filepath, "seed", pubKey, priKey, decompressDirName, specialReplaceList, prefix, DefaultSpecialReplaceList)
		fmt.Printf("replaceFileVar err=%v\n", err)
	}

	fmt.Println("done")
}
