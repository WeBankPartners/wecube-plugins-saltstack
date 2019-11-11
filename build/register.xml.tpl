<?xml version="1.0" encoding="UTF-8"?>
<package name="wecube-plugins-saltstack" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包
    <packageDependencies>
        <packageDependency name='xxx' version='1.0'/>
        <packageDependency name='xxx233' version='1.5'/>
    </packageDependencies> -->

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单
    <menus>
        <menu code='JOBS_SERVICE_CATALOG_MANAGEMENT' cat='JOBS' displayName="Servive Catalog Management">/service-catalog</menu>
        <menu code='JOBS_TASK_MANAGEMENT' cat='JOBS' displayName="Task Management">/task-management</menu>
    </menus> -->

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系
    <dataModel>
        <entity name="service_catalogue" displayName="服务目录" description="服务目录模型">
            <attribute name="id" datatype="int" description="唯一ID"/>
            <attribute name="name" datatype="string" description="名字"/>
            <attribute name="status" datatype="string" description="状态"/>
        </entity>
    </dataModel> -->

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数
    <systemParameters>
        <systemParameter name="xxx" defaultValue='xxxx' scopeType='global'/>
        <systemParameter name="xxx" defaultValue='xxxx' scopeType='plugin-package'/>
    </systemParameters> -->

    <!-- 5.权限设定
    <authorities>
        <authority systemRoleName="admin" >
            <menu code="JOBS_SERVICE_CATALOG_MANAGEMENT" />
            <menu code="JOBS_TASK_MANAGEMENT" />
        </authority >
        <authority systemRoleName="wecube_operator" >
            <menu code="JOBS_TASK_MANAGEMENT" />
        </authority >
    </authorities> -->

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{IMAGENAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,/home/app/data/minions_pki:/etc/salt/pki/master/minions,/home/app/wecube-plugins-saltstack/logs:/home/app/wecube-plugins-saltstack/logs,/home/app/data:/home/app/data" envVariables="minion_master_ip={{minion_master_ip}},minion_passwd=Ab888888,minion_port=22"/>
        <!-- <mysql schema="service_management" initFileName="init.sql" upgradeFileName="upgrade.sql"/>
        <s3 bucketName="service_management"/> -->
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="file">
            <interface name="copy" path="/wecube-plugins-saltstack/v1/deploy/file/copy">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">destinationPath</parameter>
                    <parameter datatype="string">unpack</parameter>
                    <parameter datatype="string">fileOwner</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="agent">
            <interface name="install" path="/wecube-plugins-saltstack/v1/deploy/agent/install">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">host</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="variable">
            <interface name="copy" path="/wecube-plugins-saltstack/v1/deploy/variable/replace">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">confFiles</parameter>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">variableList</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">s3PkgPath</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="script">
            <interface name="run" path="/wecube-plugins-saltstack/v1/deploy/script/run">
                <input-parameters>
                <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">endpointType</parameter>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">runAs</parameter>
                    <parameter datatype="string">args</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="user">
            <interface name="add" path="/wecube-plugins-saltstack/v1/deploy/user/add">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">userName</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">userGroup</parameter>
                    <parameter datatype="string">userId</parameter>
                    <parameter datatype="string">groupId</parameter>
                    <parameter datatype="string">homeDir</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="remove" path="/wecube-plugins-saltstack/v1/deploy/user/remove">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">userName</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="database">
            <interface name="runScript" path="/wecube-plugins-saltstack/v1/deploy/database/runScript">
                <input-parameters>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">seed</parameter>
                    <parameter datatype="string">host</parameter>
                    <parameter datatype="string">userName</parameter>
                    <parameter datatype="string">password</parameter>
                    <parameter datatype="string">port</parameter>
                    <parameter datatype="string">databaseName</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="disk">
            <interface name="getUnformatedDisk" path="/wecube-plugins-saltstack/v1/deploy/disk/getUnformatedDisk">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">unformatedDisks</parameter>
                </output-parameters>
            </interface>
            <interface name="formatAndMountDisk" path="/wecube-plugins-saltstack/v1/deploy/disk/formatAndMountDisk">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">diskName</parameter>
                    <parameter datatype="string">fileSystemType</parameter>
                    <parameter datatype="string">mountDir</parameter>
                </input-parameters>

                <output-parameters>
                <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface name="new" path="/wecube-plugins-saltstack/v1/deploy/apply-deployment/new">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">userName</parameter>
                    <parameter datatype="string">destinationPath</parameter>
                    <parameter datatype="string">confFiles</parameter>
                    <parameter datatype="string">variableList</parameter>
                    <parameter datatype="string">startScript</parameter>
                    <parameter datatype="string">args</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface name="update" path="/wecube-plugins-saltstack/v1/deploy/apply-deployment/update">
                <input-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">target</parameter>
                    <parameter datatype="string">userName</parameter>
                    <parameter datatype="string">endpoint</parameter>
                    <parameter datatype="string">confFiles</parameter>
                    <parameter datatype="string">destinationPath</parameter>
                    <parameter datatype="string">variableList</parameter>
                    <parameter datatype="string">stopScript</parameter>
                    <parameter datatype="string">startScript</parameter>
                    <parameter datatype="string">args</parameter>
                </input-parameters>
                <output-parameters>
                <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
    </plugins>
</package>