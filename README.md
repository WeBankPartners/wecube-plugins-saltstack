# SaltStack插件
SaltStack插件里包含salt-master服务。


## 技术实现
WeCube部署完SaltStack插件后，对于新创建的机器，可通过该插件里的初始化接口来安装SaltStack的agent，一旦安装完agent，可通过SaltStack插件让机器执行相关命令。

该插件包的开发语言为golang，开发过程中每加一个新的资源管理接口，同时需要修改build下的register.xm.tpl文件，在里面同步更新相关接口的url、入参和出参。


## 主要功能

- 文件操作：拷贝文件；
- agent操作：安装；
- 变量替换操作：复制替换；
- 脚本操作：执行；
- 用户操作：新增用户、删除用户；
- 数据库操作：执行脚本；
- 物料包操作：目录结构查询、差异化文件查询；
- 数据盘操作：检查未挂载盘、挂载盘；


## 编译打包
插件采用容器化部署。

如何编译插件，请查看以下文档
[SaltStack插件编译文档](docs/compile/wecube-plugins-saltstack_compile_guide.md)


## 插件运行
插件包制作完成后，需要通过WeCube的插件管理界面进行注册才能使用。运行插件的主机需提前安装好docker。

