package test

import (
	"testing"
)

func TestSaltApiPlugin(t *testing.T) {
	//install_package_tree()
	//getFileMD5Sum()
	//checkFileMD5Sum()
	//checkFileMeta()
	//downloadFileFromUrl()
}

func install_package_tree() {
	guid_1 := "guid_1"
	inputs := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"client":"local",
			"tgt":"*",
			"fun":"cmd.run",
			"arg":["yum install -y tree"]
		}]
	}
	`
	CallPlugin("salt-api", "call", inputs)
}

func getFileMD5Sum() {
	guid_1 := "guid_1"
	inputs := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"client":"local",
			"tgt":"*",
			"fun":"file.get_hash",
			"arg":["/home/app/pkgs/hello.sh","md5"]
		}]
	}
	`
	CallPlugin("salt-api", "call", inputs)
}

func checkFileMD5Sum() {
	guid_1 := "guid_1"
	inputs := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"client":"local",
			"tgt":"*",
			"fun":"file.check_hash",
			"arg":["/home/app/pkgs/hello.sh","md5:3d3bb0ae1af93b239aefaced41602e58"]
		}]
	}
	`
	CallPlugin("salt-api", "call", inputs)
}

func checkFileMeta() {
	guid_1 := "guid_1"
	inputs := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"client":"local",
			"tgt":"*",
			"fun":"file.check_file_meta",
			"arg":["/home/app/pkgs/hello.sh","","salt://pkgs/hello.sh","{hash_type: 'md5', 'hsum': <md5sum>}","root","root","755","base"]
		}]
	}
	`

	//salt '*' file.check_file_meta /home/app/pkgs/hello.sh salt://pkgs/hello.sh '{hash_type: md5, hsum: xxxx}' root, root, '755' base
	CallPlugin("salt-api", "call", inputs)
}

func downloadFileFromUrl() {
	guid_1 := "guid_1"
	inputs := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"service_url":"https://127.0.0.1:8082",
			"token":"aa534397e9cb94d147f5953e94f2baa7d03d2a28",
			"client":"local",
			"tgt":"*",
			"fun":"cp.get_url",
			"arg":["http://127.0.0.1:9090/demo.txt","/home/app/pkgs/demo.txt"]
		}]
	}
	`
	CallPlugin("salt-api", "call", inputs)
}
