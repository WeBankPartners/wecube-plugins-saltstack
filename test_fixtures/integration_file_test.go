package test

import (
	"testing"
)

func TestFilePlugin(t *testing.T) {
	//copyFilesFromMasterToAgent()
	//copyFilesFromHttpServerToAgent()
}

func copyFilesFromMasterToAgent() {
	guid_1 := "guid_1"
	copyFileInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"target":"127.0.0.1",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"source_path":"salt://pkgs/hello.sh",
			"destination_path":"/home/app/pkgs/hello.sh"
		}]
	}
	`
	CallPlugin("file", "copy", copyFileInput)
}

func copyFilesFromHttpServerToAgent() {
	guid_1 := "guid_1"
	copyFileInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"target":"127.0.0.1",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"source_path":"http://127.0.0.1:9090/demo.txt",
			"destination_path":"/home/app/pkgs/demo.txt"
		}]
	}
	`
	CallPlugin("file", "copy", copyFileInput)
}
