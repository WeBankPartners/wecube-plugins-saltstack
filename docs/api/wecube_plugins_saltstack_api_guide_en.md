# SaltStack Plugin API Guide

English / [中文](wecube_plugins_saltstack_api_guide.md)

Provide a unified interface definition which is clear and convenient for users to use.

## API Resources

**Agent Action**

- [Install Agent](#agent-install)

**File Action**

- [Copy File](#file-copy) 

**Variable Action**

- [Replace Variable](#variable-replace)

**Script Action**

- [Run Script](#script-run)

**User Management Action**

- [Add Linux User](#user-add)  
- [Remove Linux User](#user-remove)  

**Database Action**

- [Run Database Script](#database-runScript)

**Disk Action**

- [Get Unformated Disk](#disk-getUnformatedDisk)  
- [Format and Mount Disk](#disk-formatAndMountDisk)

**Application Deployment Action**

- [Deploy Application](#deploy-install)  
- [Upgrade Application](#deploy-upgrade)

## API and Examples

### Agent Action

#### <span id="agent-install">Install Agent</span>
[POST] /v1/deploy/agent/install

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
host|string|Yes|Target host IP
password|string|Yes|Password of target host ROOT user
seed|string|Yes|Secret key seed of target host ROOT user

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|more output information

##### Example
Input:
```
{
    "inputs":[{
        "guid":"0012_0000000067",
        "host":"127.0.0.1",
        "seed":"sample@123456",
        "password": "3399474758b7b0565600dikw29b92c02"
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "guid": "0012_0000000067",
                "detail": "127.0.0.1:\n----------\n          ID: salt_repo\n    Function: file.recurse\n        Name: /etc/yum.repos.d\n      Result: True\n     Comment: Recursively updated /etc/yum.repos.d\n     Started: 14:39:07.162408\n    Duration: 318.305 ms\n     Changes:   \n              ----------\n              /etc/yum.repos.d/salt-repo.repo:\n                  ----------\n                  diff:\n                      New file\n                  mode:\n                      0644\n----------\n          ID: salt_minion_purge\n    Function: pkg.purged\n      Result: True\n     Comment: All specified packages are already absent\n     Started: 14:39:10.798416\n    Duration: 893.87 ms\n     Changes:   \n----------\n          ID: salt_minion_install\n    Function: pkg.installed\n      Result: True\n     Comment: The following packages were installed/updated: salt-minion\n     Started: 14:39:11.709269\n    Duration: 27070.787 ms\n     Changes:   \n              ----------\n              gpg-pubkey.(none):\n                  ----------\n                  new:\n                      352c64e5-52ae6884,de57bfbe-53a9be98,f4a80eb5-53a7ff4b\n                  old:\n                      352c64e5-52ae6884,f4a80eb5-53a7ff4b\n              libsodium:\n                  ----------\n                  new:\n                      1.0.18-1.el7\n                  old:\n              libtomcrypt:\n                  ----------\n                  new:\n                      1.17-26.el7\n                  old:\n              libtommath:\n                  ----------\n                  new:\n                      0.42.0-6.el7\n                  old:\n              openpgm:\n                  ----------\n                  new:\n                      5.2.122-2.el7\n                  old:\n              python-tornado:\n                  ----------\n                  new:\n                      4.2.1-5.el7\n                  old:\n              python-zmq:\n                  ----------\n                  new:\n                      15.3.0-3.el7\n                  old:\n              python2-crypto:\n                  ----------\n                  new:\n                      2.6.1-16.el7\n                  old:\n              python2-futures:\n                  ----------\n                  new:\n                      3.1.1-5.el7\n                  old:\n              python2-msgpack:\n                  ----------\n                  new:\n                      0.5.6-5.el7\n                  old:\n              python2-psutil:\n                  ----------\n                  new:\n                      2.2.1-5.el7\n                  old:\n              salt:\n                  ----------\n                  new:\n                      2019.2.0-1.el7\n                  old:\n              salt-minion:\n                  ----------\n                  new:\n                      2019.2.0-1.el7\n                  old:\n              zeromq:\n                  ----------\n                  new:\n                      4.1.4-7.el7\n                  old:\n----------\n          ID: salt_minion_conf\n    Function: file.managed\n        Name: /etc/salt/minion\n      Result: True\n     Comment: File /etc/salt/minion updated\n     Started: 14:39:38.791480\n    Duration: 88.582 ms\n     Changes:   \n              ----------\n              diff:\n                  --- \n                  +++ \n                  @@ -13,7 +13,8 @@\n                   \n                   # Set the location of the salt master server. If the master server cannot be\n                   # resolved, then the minion will fail to start.\n                  -#master: salt\n                  +master: \n                  +  - 127.0.0.1\n                   \n                   # Set http proxy information for the minion when doing requests\n                   #proxy_host:\n                  @@ -76,7 +77,7 @@\n                   # retry_dns_count: 3\n                   \n                   # Set the port used by the master reply and authentication server.\n                  -#master_port: 4506\n                  +master_port: 4506\n                   \n                   # The user to run salt.\n                   #user: root\n                  @@ -110,6 +111,7 @@\n                   # same machine but with different ids, this can be useful for salt compute\n                   # clusters.\n                   #id:\n                  +id: 127.0.0.1\n                   \n                   # Cache the minion id to a file when the minion's id is not statically defined\n                   # in the minion config. Defaults to \"True\". This setting prevents potential\n                  @@ -243,7 +245,7 @@\n                   # authorization from it. master_tries will still cycle through all\n                   # the masters in a given try, so it is appropriate if you expect\n                   # occasional downtime from the master(s).\n                  -#master_tries: 1\n                  +master_tries: -1\n                   \n                   # If authentication fails due to SaltReqTimeoutError during a ping_interval,\n                   # cause sub minion process to restart.\n                  @@ -858,12 +860,12 @@\n                   \n                   # Overall state of TCP Keepalives, enable (1 or True), disable (0 or False)\n                   # or leave to the OS defaults (-1), on Linux, typically disabled. Default True, enabled.\n                  -#tcp_keepalive: True\n                  +tcp_keepalive: True\n                   \n                   # How long before the first keepalive should be sent in seconds. Default 300\n                   # to send the first keepalive after 5 minutes, OS default (-1) is typically 7200 seconds\n                   # on Linux see /proc/sys/net/ipv4/tcp_keepalive_time.\n                  -#tcp_keepalive_idle: 300\n                  +tcp_keepalive_idle: 60\n                   \n                   # How many lost probes are needed to consider the connection lost. Default -1\n                   # to use OS defaults, typically 9 on Linux, see /proc/sys/net/ipv4/tcp_keepalive_probes.\n----------\n          ID: salt_minion_service\n    Function: service.running\n        Name: salt-minion\n      Result: True\n     Comment: Service salt-minion has been enabled, and is running\n     Started: 14:39:40.306011\n    Duration: 810.671 ms\n     Changes:   \n              ----------\n              salt-minion:\n                  True\n\nSummary for 127.0.0.1\n------------\nSucceeded: 5 (changed=4)\nFailed:    0\n------------\nTotal states run:     5\nTotal run time:  29.182 s\n"
            }
        ]
    }
}
```

### File Action

#### <span id="file-copy">Copy File</span>
[POST] /v1/deploy/file/copy

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
endpoint|string|Yes|The full path where the file is stored in
target|string|Yes|Target host IP
destinationPath|string|Yes|The destination where the file copied will be stored in

##### Output Parameters
参数名称|类型|描述
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More information

##### Example

Input:
```
{
  "inputs":[{
        "guid":"10002_000000001",
        "endpoint":"http://127.0.0.1:9000/brankbao/unpack-demo.tar",
        "target":"127.0.0.1",
        "destinationPath":"/data/app/scripts/unpack-demo.tar"
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
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

### Variable Action

#### <span id="variable-replace">Replace Variable</span>
[POST] /v1/deploy/variable/replace

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
endpoint|string|Yes|The full path where the application package is stored in
confFiles|string|Yes|The full path where files needing to replace variables are stored in, and using "\|" to distinguish multi files
variableList|string|Yes|Variable List, format: "Name=tom, Age=10, Dog = test1, Cat = tet2"

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
s3PkgPath|string|The absolute path of the application package after the variables are replaced

##### Example
Input:
```
{
    "inputs": [{
        "guid":"10003_000000001",
        "endpoint": "http://127.0.0.1:9000/brankbao/wecube-demo_v2.0.zip",
        "confFiles": "beego-demo/conf/app.conf",
        "variableList":"env=prod"
	}]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "guid": "10003_000000001",
                "s3PkgPath": "http://127.0.0.1:9000/brankbao/wecube-demo_v2.0-201909231548.zip"
            }
        ]
    }
}
```

### Script Action

#### <span id="script-run">Run Script</span>
[POST] /v1/deploy/script/run

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
endpointType|string|Yes|Type of script, options: "S3": script which is stored in S3 server;"LOCAL" : Script which is stored in local
endpoint|string|Yes|The absolute path where script is stored in
target|string|Yes|Target host IP
runAs|string|No|User running the script
args|string|No|The running parameters

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More information
target|string|Target Host' IP
retCode|string|Return code

##### Example
Input:
```
{
  "inputs":[{
        "guid":"10004_000000001",
        "endpointType":"LOCAL",
        "endpoint":"/data/app/scripts/test.sh",
        "target":"127.0.0.1",
        "runAs":"app",
        "args":""
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "target": "127.0.0.1",
                "retCode": 0,
                "detail": "127.0.0.1:",
                "guid": "10004_000000001"
            }
        ]
    }
}
```

### User Management Action

#### <span id="user-add">Add Linux User</span>
[POST] /v1/deploy/user/add

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
target|string|Yes|Target host' IP
userId|string|No|User ID
userName|string|Yes|User name
password|string|Yes|Secret
userGroup|string|No|User group
groupId|string|No|Group ID
homeDir|string|No|User home diretory

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More Information

##### Example
Input:
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

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
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

#### <span id="user-remove">Remove Linux User</span>
[POST] /v1/deploy/user/remove

##### Input Parameters
name|type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
target|string|Yes|Target host IP
userName|string|Yes|User name

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More information

##### Example:
Input:
```
{
  "inputs":[{
        "guid":"10005_000000001",
        "target":"127.0.0.1",
        "userName":"app"
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
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

### Database Action

#### <span id="database-runScript">Run Database Script</span>
[POST] /v1/deploy/database/runScript

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
host|string|Yes|Database IP
port|string|Yes|Database port
userName|string|Yes|Name of database user
password|string|Yes|Password of database user 
seed|string|Yes|Secret seed of database user 
databaseName|string|Yes|Database name
endpoint|string|Yes|The full path where the script is stored in

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More information

##### Example
Input:
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
        "endpoint":"http://127.0.0.1:9000/scripts/createDB.sql"
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "guid": "10006_000000001",
                "detail": "127.0.0.1:"
            }
        ]
    }
}
```

### Disk Action

#### <span id="disk-getUnformatedDisk">Get Unformated Disk</span>
[POST] /v1/deploy/disk/getUnformatedDisk

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
endpoint|string|Yes|Target host IP

##### Output parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
unformatedDisks|array|Unformated disk list

##### Example
Input:
```
{
    "inputs":[{
		"guid":"10007_000000001",
		"target":"127.0.0.1"
    }]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
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

#### <span id="disk-formatAndMountDisk">Formate and Mount Disk</span>
[POST] /v1/deploy/disk/formatAndMountDisk

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
target|string|Yes|Target host IP
diskName|string|Yes|Disk name
fileSystemType|string|Yes|File system type
mountDir|string|Yes|The directory to mount 

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
detail|string|More information

##### Example
Input:
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

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
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

### Application Deployment Action

#### <span id="deploy-install">Deploy Application</span>
[POST] /v1/deploy/apply-deployment/new

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
target|string|Yes|Target host IP
endpoint|string|Yes|The full path where the file is stored in
userName|string|Yes|User name
destinationPath|string|Yes|The destination where the file copied will be stored in
confFiles|string|No|The full path where files needing to replace variables are stored in, and using "\|" to distinguish multi files
variableList|string|No|Variable List, format: "Name=tom, Age=10, Dog = test1, Cat = tet2"
startScript|string|Yes|The full path where the start script is stored in
args|string|No|The running parameters

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
userDetail|string|More information about creating user
fileDetail|string|More information about copying application package
s3PkgPath|string|The absolute path of the application package after the variables are replaced
target|string|Target host IP
retCode|string|Return code
runScriptDetail|string|More information about running the start script

##### Example
Input:
```
{
    "inputs": [
        {
            "guid": "0015_0000000079",
            "userName": "app",
            "endpoint": "http://127.0.0.1:9000/wecube-artifact/775c3b239d8a76b1914377deb346a8a1_edp-core-app_v2.6.tgz",
            "target": "127.0.0.1",
            "destinationPath": "/data/app/",
            "confFiles": "edp-core-app_v2.6/conf/app.conf|/edp-core-app_v2.6/controllers/test/test.conf",
            "variableList": "env=GZ3-SUBSYSTEM,appname=demo-app,diskaa=50,versionbb=STGi,httpport=2019,vip=127.0.0.1",
            "startScript": "/data/app/edp-core-app_v2.6/wecube-demo.sh"
        }
    ]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "guid": "0015_0000000079",
                "userDetail": "{\"return\": [{\"127.0.0.1\": {\"pid\": 19231, \"retcode\": 0, \"stderr\": \"\", \"stdout\": \"\"}}]}",
                "fileDetail": "{\"return\": [{\"127.0.0.1\": \"909d22a818257c502557b7abe9ec636d\"}]}",
                "s3PkgPath": "http://127.0.0.1:9000/wecube-artifact/775c3b239d8a76b1914377deb346a8a1_edp-core-app_v2.6-201910121000.tgz",
                "target": "127.0.0.1",
                "retCode": 0,
                "runScriptDetail": "127.0.0.1:sudo: no tty present and no askpass program specified"
            }
        ]
    }
}
```

#### <span id="deploy-upgrade">Upgrade Application</span>
[POST] /v1/deploy/apply-deployment/update

##### Input Parameters
Name|Type|Required|Description
:--|:--|:--|:-- 
guid|string|Yes|Globally unique CI type ID
target|string|Yes|Target host IP
endpoint|string|Yes|The full path where the file is stored in
userName|string|Yes|User name
destinationPath|string|Yes|The destination where the file copied will be stored in
confFiles|string|No|The full path where files needing to replace variables are stored in, and using "\|" to distinguish multi files
variableList|string|No|Variable List, format: "Name=tom, Age=10, Dog = test1, Cat = tet2"
startScript|string|Yes|The full path where the start script is stored in
stopScript|string|Yes|The full path where the stop script is stored in
args|string|No|The running parameters

##### Output Parameters
Name|Type|Description
:--|:--|:--    
guid|string|Globally unique CI type ID
fileDetail|string|More information about copying application package
s3PkgPath|string|The absolute path of the application package after the variables are replaced
target|string|Target host IP
retCode|string|Return code
runStartScriptDetail|string|More information about running the start script
runStopScriptDetail|string|More information about running the stop script

##### Example
Input:
```
{
    "inputs": [
        {
            "guid": "0015_0000000079",
            "userName": "app",
            "endpoint": "http://127.0.0.1:9000/wecube-artifact/775c3b239d8a76b1914377deb346a8a1_edp-core-app_v2.6.tgz",
            "target": "127.0.0.1",
            "destinationPath": "/data/app/",
            "confFiles": "edp-core-app_v2.6/conf/app.conf|/edp-core-app_v2.6/controllers/test/test.conf",
            "variableList": "env=GZ3-SUBSYSTEM,appname=demo-app,diskaa=50,versionbb=STGi,httpport=2019,vip=127.0.0.1",
            "startScript": "/data/app/edp-core-app_v2.6/wecube-demo.sh",
            "stopScript": "/data/app/edp-core-app_v2.6/stop.sh"
        }
    ]
}
```

Output:
```
{
    "resultCode": "0",
    "resultMessage": "success",
    "results": {
        "outputs": [
            {
                "guid": "0015_0000000079",
                "fileDetail": "{\"return\": [{\"127.0.0.1\": \"2e1538f77758e9e026fdcaa9ed4ad388\"}]}",
                "s3PkgPath": "http://127.0.0.1:9000/wecube-artifact/775c3b239d8a76b1914377deb346a8a1_edp-core-app_v2.6-201910121517.tgz",
                "target": "127.0.0.1",
                "retCode": 0,
                "runStartScriptDetail": "127.0.0.1:sudo: no tty present and no askpass program specified",
                "runStopScriptDetail": "127.0.0.1:"
            }
        ]
    }
}
```