# WECUBE PLUGINS SALTSTACK API GUIDE
  
提供统一接口定义，为使用者提供清晰明了的使用方法。

## API 操作资源（Resources）:  

**Agent操作**

- [Agent安装](#agent-install)  

**文件操作**

- [文件拷贝](#file-copy)  

**变量操作**

- [变量替换](#variable-replace)  

**脚本操作**

- [脚本执行](#script-run)  

**用户管理操作**

- [Linux用户新增](#user-add)  
- [Linux用户删除](#user-remove)  

**数据库操作**

- [数据库脚本执行](#database-runScript)  

**数据盘操作**

- [查询未挂载数据盘](#disk-getUnformatedDisk)  
- [挂载数据盘](#disk-formatAndMountDisk) 

**部署操作**

- [全量部署](#deploy-install)  
- [增量部署](#deploy-upgrade)  


## API 概览及实例：  

### Agent操作

#### <span id="agent-install">Agent安装</span>
[POST] /v1/deploy/agent/install

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
host|string|是|目标机器IP
password|string|是|目标机器ROOT用户密码
seed|string|是|目标机器ROOT用户密钥种子

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：

```
{
    "inputs":[{
			"guid":"0012_0000000067",
			"host":"127.0.0.1",
			"seed":"sample@123456",
			"password": "3399474758b7b0565600dikw29b92c02"
  }
  ]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0012_0000000067",
                "detail": "127.0.0.1:\n----------\n          ID: salt_repo\n    Function: file.recurse\n        Name: /etc/yum.repos.d\n      Result: True\n     Comment: Recursively updated /etc/yum.repos.d\n     Started: 14:39:07.162408\n    Duration: 318.305 ms\n     Changes:   \n              ----------\n              /etc/yum.repos.d/salt-repo.repo:\n                  ----------\n                  diff:\n                      New file\n                  mode:\n                      0644\n----------\n          ID: salt_minion_purge\n    Function: pkg.purged\n      Result: True\n     Comment: All specified packages are already absent\n     Started: 14:39:10.798416\n    Duration: 893.87 ms\n     Changes:   \n----------\n          ID: salt_minion_install\n    Function: pkg.installed\n      Result: True\n     Comment: The following packages were installed/updated: salt-minion\n     Started: 14:39:11.709269\n    Duration: 27070.787 ms\n     Changes:   \n              ----------\n              gpg-pubkey.(none):\n                  ----------\n                  new:\n                      352c64e5-52ae6884,de57bfbe-53a9be98,f4a80eb5-53a7ff4b\n                  old:\n                      352c64e5-52ae6884,f4a80eb5-53a7ff4b\n              libsodium:\n                  ----------\n                  new:\n                      1.0.18-1.el7\n                  old:\n              libtomcrypt:\n                  ----------\n                  new:\n                      1.17-26.el7\n                  old:\n              libtommath:\n                  ----------\n                  new:\n                      0.42.0-6.el7\n                  old:\n              openpgm:\n                  ----------\n                  new:\n                      5.2.122-2.el7\n                  old:\n              python-tornado:\n                  ----------\n                  new:\n                      4.2.1-5.el7\n                  old:\n              python-zmq:\n                  ----------\n                  new:\n                      15.3.0-3.el7\n                  old:\n              python2-crypto:\n                  ----------\n                  new:\n                      2.6.1-16.el7\n                  old:\n              python2-futures:\n                  ----------\n                  new:\n                      3.1.1-5.el7\n                  old:\n              python2-msgpack:\n                  ----------\n                  new:\n                      0.5.6-5.el7\n                  old:\n              python2-psutil:\n                  ----------\n                  new:\n                      2.2.1-5.el7\n                  old:\n              salt:\n                  ----------\n                  new:\n                      2019.2.0-1.el7\n                  old:\n              salt-minion:\n                  ----------\n                  new:\n                      2019.2.0-1.el7\n                  old:\n              zeromq:\n                  ----------\n                  new:\n                      4.1.4-7.el7\n                  old:\n----------\n          ID: salt_minion_conf\n    Function: file.managed\n        Name: /etc/salt/minion\n      Result: True\n     Comment: File /etc/salt/minion updated\n     Started: 14:39:38.791480\n    Duration: 88.582 ms\n     Changes:   \n              ----------\n              diff:\n                  --- \n                  +++ \n                  @@ -13,7 +13,8 @@\n                   \n                   # Set the location of the salt master server. If the master server cannot be\n                   # resolved, then the minion will fail to start.\n                  -#master: salt\n                  +master: \n                  +  - 10.0.0.8\n                   \n                   # Set http proxy information for the minion when doing requests\n                   #proxy_host:\n                  @@ -76,7 +77,7 @@\n                   # retry_dns_count: 3\n                   \n                   # Set the port used by the master reply and authentication server.\n                  -#master_port: 4506\n                  +master_port: 4506\n                   \n                   # The user to run salt.\n                   #user: root\n                  @@ -110,6 +111,7 @@\n                   # same machine but with different ids, this can be useful for salt compute\n                   # clusters.\n                   #id:\n                  +id: 127.0.0.1\n                   \n                   # Cache the minion id to a file when the minion's id is not statically defined\n                   # in the minion config. Defaults to \"True\". This setting prevents potential\n                  @@ -243,7 +245,7 @@\n                   # authorization from it. master_tries will still cycle through all\n                   # the masters in a given try, so it is appropriate if you expect\n                   # occasional downtime from the master(s).\n                  -#master_tries: 1\n                  +master_tries: -1\n                   \n                   # If authentication fails due to SaltReqTimeoutError during a ping_interval,\n                   # cause sub minion process to restart.\n                  @@ -858,12 +860,12 @@\n                   \n                   # Overall state of TCP Keepalives, enable (1 or True), disable (0 or False)\n                   # or leave to the OS defaults (-1), on Linux, typically disabled. Default True, enabled.\n                  -#tcp_keepalive: True\n                  +tcp_keepalive: True\n                   \n                   # How long before the first keepalive should be sent in seconds. Default 300\n                   # to send the first keepalive after 5 minutes, OS default (-1) is typically 7200 seconds\n                   # on Linux see /proc/sys/net/ipv4/tcp_keepalive_time.\n                  -#tcp_keepalive_idle: 300\n                  +tcp_keepalive_idle: 60\n                   \n                   # How many lost probes are needed to consider the connection lost. Default -1\n                   # to use OS defaults, typically 9 on Linux, see /proc/sys/net/ipv4/tcp_keepalive_probes.\n----------\n          ID: salt_minion_service\n    Function: service.running\n        Name: salt-minion\n      Result: True\n     Comment: Service salt-minion has been enabled, and is running\n     Started: 14:39:40.306011\n    Duration: 810.671 ms\n     Changes:   \n              ----------\n              salt-minion:\n                  True\n\nSummary for 127.0.0.1\n------------\nSucceeded: 5 (changed=4)\nFailed:    0\n------------\nTotal states run:     5\nTotal run time:  29.182 s\n"
            }
        ]
    }
}
```


### 文件操作

#### <span id="file-copy">文件拷贝</span>
[POST] /v1/deploy/file/copy

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
endpoint|string|是|文件存储在对象存储中的地址，全路径
target|string|是|目标机器IP
destination_path|string|是|目标路径，全路径，包括目标文件名

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：
```
{
  "inputs":[{
			"guid":"10002_000000001",
			"endpoint":"http://127.0.0.1:9000/brankbao/unpack-demo.tar",
			"target":"127.0.0.1",
			"destination_path":"/data/app/scripts/unpack-demo.tar"
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10002_000000001",
                "detail": "{\"return\": [{\"127.0.0.1\": \"bae0cbb98dd0f6d346ada2157922f799\"}]}"
            }
        ]
    }
}
```

### 变量操作

#### <span id="variable-replace">变量替换</span>
[POST] /v1/deploy/variable/replace

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
endpoint|string|是|应用包存储在对象存储中的地址 ，全路径
conf_files_path|string|是|差异化变量文件在应用包中的相对路径，多个文件以分号";"分隔
variable_list|string|是|变量列表， 格式："Name=tom, Age=10, Dog = test1, Cat = tet2"

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
s3_pkg_path|string|变量替换后的应用包在对象存储中的绝对路径 

##### 示例：
输入：
```
{
    "inputs": [{
			"guid":"10003_000000001",
			"endpoint": "http://127.0.0.1:9000/brankbao/wecube-demo_v2.0.zip",
			"conf_files": "beego-demo/conf/app.conf",
			"variable_list":"env=prod"
		}]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10003_000000001",
                "s3_pkg_path": "http://127.0.0.1:9000/brankbao/wecube-demo_v2.0-201909231548.zip"
            }
        ]
    }
}
```


### 脚本操作

#### <span id="script-run">脚本执行</span>
[POST] /v1/deploy/script/run

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
end_point_type|string|是|脚本来源类型，可选值："S3"：脚本在对象存储，"LOCAL"：脚本在服务器本地
end_point|string|是|脚本存储路径， 可以是对象存储上的绝对路径，或者服务器的绝对路径
target|string|是|目标机器IP
run_as|string|是|执行用户
args|string|是|执行参数

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
detail|string|详细信息
target|string|目标机器
ret_code|string|返回码

##### 示例：
输入：
```
{
  "inputs":[{
			"guid":"10004_000000001",
			"end_point_type":"LOCAL",
			"end_point":"/data/app/scripts/test.sh",
			"target":"127.0.0.1",
			"run_as":"app",
			"args":""
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "target": "127.0.0.1",
                "ret_code": 0,
                "detail": "127.0.0.1:",
                "guid": "10004_000000001"
            }
        ]
    }
}
```


### 用户管理操作

#### <span id="user-add">Linux用户新增</span>
[POST] /v1/deploy/user/add

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|CI类型全局唯一ID
target|string|是|目标机器IP
userId|string|否|用户ID
userName|string|是|用户名
password|string|是|密码
userGroup|string|是|用户组
groupId|string|否|组ID
homeDir|string|否|用户home目录


##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：
```
{
  "inputs":[{
			"guid":"10005_000000001",
			"target":"127.0.0.1",
			"userName":"app",
			"password":"Apps@2018!",
			"userGroup":"apps"
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10005_000000001",
                "detail": "{\"return\": [{\"127.0.0.1\": {\"pid\": 13765, \"retcode\": 0, \"stderr\": \"\", \"stdout\": \"\"}}]}"
            }
        ]
    }
}
```


#### <span id="user-remove">Linux用户删除</span>
[POST] /v1/deploy/user/remove

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
target|string|是|目标机器IP
userName|string|是|用户名

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：
```
{
  "inputs":[{
			"guid":"10005_000000001",
			"target":"127.0.0.1",
			"userName":"app"
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "detail": "{\"return\": [{\"127.0.0.1\": {\"pid\": 13965, \"retcode\": 0, \"stderr\": \"\", \"stdout\": \"\"}}]}",
                "guid": "10005_000000001"
            }
        ]
    }
}
```


### 数据库操作

#### <span id="database-runScript">数据库脚本执行</span>
[POST] /v1/deploy/database/runScript

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
host|string|是|数据库IP
port|string|是|数据库端口
userName|string|是|数据库用户名
password|string|是|数据库用户密码
seed|string|是|数据库用户密钥种子
databaseName|string|是|数据库名
endpoint|string|是|文件存储在对象存储中的地址 ，全路径


##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|是|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：
```
{
  "inputs":[{
			"guid":"10006_000000001",
			"host":"127.0.0.1",
			"port":"3306",
			"userName":"mariadb",
			"password":"3dfdecb7281a498e362c03987fdd0dd9",
			"seed":"abc@12345",
			"databaseName":"we-cmdb",
			"endpoint":"http://127.0.0.1:9000/scripts/createDB.sql",
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10006_000000001"
                "detail": "127.0.0.1:",
            }
        ]
    }
}
```


### 数据盘操作

#### <span id="disk-getUnformatedDisk">查询未挂载数据盘</span>
[POST] /v1/deploy/disk/getUnformatedDisk

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
endpoint|string|是|目标机器IP


##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|CI类型全局唯一ID
unformatedDisks|string|未挂载数据盘清单

##### 示例：
输入：
```
{
  "inputs":[{
		"guid":"10007_000000001",
		"target":"127.0.0.1"
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10007_000000001",
                "unformatedDisks": [
                    "/dev/vdb"
                ]
            }
        ]
    }
}
```



#### <span id="disk-formatAndMountDisk">挂载数据盘</span>
[POST] /v1/deploy/disk/formatAndMountDisk

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
target|string|是|目标机器IP
diskName|string|是|数据盘名称
fileSystemType|string|是|文件系统
mountDir|string|是|挂载目录名

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
guid|string|是|CI类型全局唯一ID
detail|string|详细信息

##### 示例：
输入：
```
{
  "inputs":[{
		"guid":"10008_000000001",
		"target":"127.0.0.1",
		"diskName":"/dev/vdb",
		"fileSystemType":"ext4",
		"mountDir":"/data1"
  }]
}
```

输出：
```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "10008_000000001",
                "detail": "Filesystem label=\nOS type: Linux\nBlock size=4096 (log=2)\nFragment size=4096 (log=2)\nStride=0 blocks, Stripe width=0 blocks\n1310720 inodes, 5242880 blocks\n262144 blocks (5.00%) reserved for the super user\nFirst data block=0\nMaximum filesystem blocks=2153775104\n160 block groups\n32768 blocks per group, 32768 fragments per group\n8192 inodes per group\nSuperblock backups stored on blocks: \n\t32768, 98304, 163840, 229376, 294912, 819200, 884736, 1605632, 2654208, \n\t4096000\n\nAllocating group tables:   0/160\b\b\b\b\b\b\b       \b\b\b\b\b\b\bdone                            \nWriting inode tables:   0/160\b\b\b\b\b\b\b       \b\b\b\b\b\b\bdone                            \nCreating journal (32768 blocks): done\nWriting superblocks and filesystem accounting information:   0/160\b\b\b\b\b\b\b       \b\b\b\b\b\b\bdone"
            }
        ]
    }
}
```

### 部署操作

#### <span id="agent-install">全量部署</span>
开发中

#### <span id="agent-install">增量部署</span>
开发中