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
	filepath := "/Users/tylertang/Desktop/config2.conf"
	err := replaceFileVar(keyMap, filepath, "seed", "private_key", "public_key", "{cpher_A}")
	if err != nil {
		fmt.Printf("err=%v\n", err)
	}
	fmt.Println("done")
}
