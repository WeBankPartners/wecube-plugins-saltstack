package plugins

import (
	"fmt"
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

	filepath := "/Users/tylertang/Desktop/go_test/config2.conf"

	keyInfo, err := GetVariable(filepath)
	if err != nil {
		fmt.Printf("GetVariable err=%v\n", err)
	}

	fmt.Println("keyINfo", keyInfo)

	err = replaceFileVar(keyMap, filepath, "seed", "private_key", "public_key", "{cpher_A}")
	if err != nil {
		fmt.Printf("replaceFileVar err=%v\n", err)
	}

	fmt.Println("done")
}
