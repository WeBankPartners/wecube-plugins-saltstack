# saltstack插件
saltstack插件里包含salt-master服务，wecube部署完saltstack插件后，对于新创建的机器，可通过该插件里的初始化接口来安装saltstack的agent，一旦安装完agent，可通过saltstack插件让机器执行相关脚本命令。

该插件包的开发语言为golang，开发过程中每加一个新的资源管理接口，同时需要修改build下的register.xm.tpl文件，在里面同步更新相关接口的url、入参和出参。

插件包制作完成后，需通过wecube的插件管理界面进行注册才能使用，运行插件的主机需提前安装好docker。

## 编译插件包的准备工作
1. 准备一台linux主机，建议操作系统为centos7.2以上。
1. 确认已经安装好git命令,如未安装通过如下命令安装
```
yum install -y git
```
2. 确认主机上已经安装好docker命令,docker安装可参考[docker安装指引](https://github.com/WeBankPartners/we-cmdb/blob/master/cmdb-wiki/docs/install/docker_install_guide.md)

3. 确认主机上有make命令，如未安装执行如下命令安装:
```
yum install -y make
```

4. 通过netstat命令确认主机上9099,9090,4505,4506这4个端口未被占用。

## 插件包的制作
1. 使用git命令拉取插件包:
```
git clone https://github.com/WeBankPartners/wecube-plugins-saltstack.git
```

2. 通过如下命令编译和打包插件，其中PLUGIN_VERSION为插件包的版本号，编译完成后将生成一个zip的插件包
```
make package PLUGIN_VERSION=v1.0
```

