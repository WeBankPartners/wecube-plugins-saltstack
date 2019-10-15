# Build SaltStack Plugin Development Environment

English / [中文](wecube-plugins-saltstack_build_dev_env.md)

- [Linux Development Environment](#Linux)  
- [Windows Development Environment](#Windows)
- [Mac Development Environment](#Mac)

>  Notice: Pay attention to the user execution permission.

## <span id="Linux">Build Linux Development Environment</span>

1. Firstly, make the directory on the linux host.

```
mkdir -p /data/gowork/src/github.com/WeBankPartners/
```

2. Download and decompress Golang binary package.

```
cd /data/
wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz 
tar xzvf go1.12.9.linux-amd64.tar.gz 
```

3. Set Golang environmental variables. Make file golang.sh under directory /data/ with the fllowing content:

```
export GOROOT=/data/go
export GOPATH=/data/gowork
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
```

4. Run the command to make it effective.

```
source /data/golang.sh
```

5. You can run **go version** to make sure successful.

```
go version
```

6. Git clone SaltStack Plugin's source code.

```
cd /data/gowork/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-saltstack.git
```

7. Build the source code.

```
cd /data/gowork/src/github.com/WeBankPartners/wecube-plugins-saltstack
go build 
```

## <span id="Windows">Build Windows Development Environment</span>

1. Make directory D:\gowork\src\github.com\WeBankPartners.

2. If your computer does not install git client, please install [Git for Windows](https://www.git-scm.com/download/win).

3. Install [Golang for Windows](https://dl.google.com/go/go1.13.1.windows-amd64.msi). During installation, please change the install directory to D:\go\. 

4. After installation, input the command **go version** to view the Golang version.

5. Set environmental variables GOROOT and GOPATH:

```
GOROOT=D:\go
GOPATH=D:\gowork
```

6. Git clone SaltStack plugin source code. In the cmd, change the directory to  D:\gowork\src\github.com\WeBankPartners and run the following command. 

```
git clone https://github.com/WeBankPartners/wecube-plugins-saltstack.git
```

7. Build SaltStack plugin source code. Please go into D:\gowork\src\github.com\WeBankPartners\wecube-plugins-saltstack and run **go build** in the cmd command line.

```
go build 
```

## <span id="Mac">Build Mac Development Environment</span>

1. The first need Golang installed. Use brew to install Golang, as follow.

```
brew install go
```

2. You can use **go env** to view the Golang version. The GOROOT displayed at this time is the installation directory where you installed the Golang using brew.

3. Set Golang environmental variables.

```
vim ~/.bash_profile
```

Set values for GOROOT、GOPATH、GOBIN、PATH

```
GOROOT=/usr/local/go
export GOROOT
export GOPATH=/Users/gowork/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin
```

Make it effective

```
source ~/.bash_profile
```

4. Please make the directory, as follow.

```
mkdir -p /Users/gowork/go/src/github.com/WeBankPartners/
```

5. Git clone SaltStack Plugin source code.

```
cd /Users/gowork/go/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-saltstack.git
```

6. Build SaltStack Plugin source code.

```
cd /Users/gowork/go/src/github.com/WeBankPartners/wecube-plugins-saltstack/
go build
```