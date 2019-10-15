# SaltStack插件
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![](https://img.shields.io/badge/language-golang-orang.svg)
 
中文 / [English](README_EN.md)

## 简介
SaltStack插件包含salt-master,salt-api和httpd等服务，基于这些服务封装了一层对主机进行系统管理和应用管理的API。用户可通过该插件提供的API执行如下操作:
- salt-minion安装:主机安装salt-minion后,后续所有对该主机的操作都可从salt-master发起
- 文件分发：从S3对象存储中下载文件并部署到指定的目录，如果是压缩包还提供解压能力
- 变量替换操作：将安装包指定目录下的配置文件进行变量替换，并重新生成替换后的安装包放到S3对象存储上。
- bash脚本操作：可指定用户在指定主机上执行主机本地或者S3对象存储上的bash脚本。
- 用户管理操作：可在指定主机上新增用户、删除用户；
- 数据库操作：在指定的mysql数据库实例上执行S3对象存储上的sql文件
- 数据盘操作：检查指定主机是否有未挂在的数据库盘；可对指定主机上的数据库盘进行格式化并设置自动挂在到某个主机目录
- 部署操作：可指定主机下发S3上对象存储上的应用安装包，并执行指定的脚本用来启动或者停止应用


<img src="./docs/images/architectrue.png" />

## Salt-Stack插件开发环境搭建
[Salt-Stack插件开发环境搭建指引](docs/compile/wecube-plugins-saltstack_build_dev_env)
Salt-Stack编译完毕后，二进制运行必须确认本机有salt-master、salt-api、mysql client等组件才能运行。因为安装这些组件较繁琐，建议使用docker镜像运行


## Salt-Stack插件docker镜像和插件包制作
[Salt-Stack插件docker镜像包和插件包制作指引](docs/compile/wecube-plugins-saltstack_compile_guide.md)

## 独立运行Salt-Stack插件容器
执行如下命令运行Salt-Stack插件容器,其中变量HOST_IP需要替换为容器所在主机的IP，该ip在执行主机安装salt-minion时使用,TAG_NUM对应代码最后一次提交的commit号

```
docker run -d  --restart=unless-stopped -v /etc/localtime:/etc/localtime -e minion_master_ip={$HOST_IP}} -e minion_passwd=Ab888888 -e minion_port=22 -p 9099:80 -p 9090:8080 -p 4505:4505 -p 4506:4506 -p 8082:8082 --privileged=true  -v /home/app/data/minions_pki:/etc/salt/pki/master/minions -v /home/app/wecube-plugins-saltstack/logs:/home/app/wecube-plugins-saltstack/logs -v /home/app/data:/home/app/data wecube-plugins-saltstack:{$TAG_NUM}}
```

## 验证插件


## API使用说明
关于Salt-Stack插件的API说明，请查看以下文档
[SaltStack插件API手册](docs/api/wecube_plugins_saltstack_api_guide.md)

## License
Salt-Stack插件是基于 Apache License 2.0 协议
