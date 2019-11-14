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
        <docker imageName="{{IMAGENAME}}" containerName="{{IMAGENAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,/home/app/data/minions_pki:/etc/salt/pki/master/minions,{{base_mount_path}}/wecube-plugins-saltstack/logs:/home/app/wecube-plugins-saltstack/logs,/home/app/data:/home/app/data" envVariables="minion_master_ip={{minion_master_ip}},minion_passwd=Ab888888,minion_port=22"/>
        <!-- <mysql schema="service_management" initFileName="init.sql" upgradeFileName="upgrade.sql"/>
        <s3 bucketName="service_management"/> -->
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="file">
            <interface action="copy" path="/wecube-plugins-saltstack/v1/deploy/file/copy">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">destinationPath</parameter>
                    <parameter datatype="string" required="N">unpack</parameter>
                    <parameter datatype="string" required="N">fileOwner</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="agent">
            <interface action="install" path="/wecube-plugins-saltstack/v1/deploy/agent/install">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">seed</parameter>
                    <parameter datatype="string" required="Y">password</parameter>
                    <parameter datatype="string" required="Y">host</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="variable">
            <interface action="copy" path="/wecube-plugins-saltstack/v1/deploy/variable/replace">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">confFiles</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">variableList</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">s3PkgPath</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="script">
            <interface action="run" path="/wecube-plugins-saltstack/v1/deploy/script/run">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">endpointType</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="N">runAs</parameter>
                    <parameter datatype="string" required="N">args</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="user">
            <interface action="add" path="/wecube-plugins-saltstack/v1/deploy/user/add">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="N">password</parameter>
                    <parameter datatype="string" required="N">userGroup</parameter>
                    <parameter datatype="string" required="N">userId</parameter>
                    <parameter datatype="string" required="N">groupId</parameter>
                    <parameter datatype="string" required="N">homeDir</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface action="remove" path="/wecube-plugins-saltstack/v1/deploy/user/remove">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="database">
            <interface action="runScript" path="/wecube-plugins-saltstack/v1/deploy/database/runScript">
                <input-parameters>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">seed</parameter>
                    <parameter datatype="string" required="Y">host</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="Y">password</parameter>
                    <parameter datatype="string" required="N">port</parameter>
                    <parameter datatype="string" required="N">databaseName</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="disk">
            <interface action="getUnformatedDisk" path="/wecube-plugins-saltstack/v1/deploy/disk/getUnformatedDisk">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">unformatedDisks</parameter>
                </output-parameters>
            </interface>
            <interface action="formatAndMountDisk" path="/wecube-plugins-saltstack/v1/deploy/disk/formatAndMountDisk">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">diskName</parameter>
                    <parameter datatype="string" required="Y">fileSystemType</parameter>
                    <parameter datatype="string" required="Y">mountDir</parameter>
                </input-parameters>

                <output-parameters>
                <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface action="new" path="/wecube-plugins-saltstack/v1/deploy/apply-deployment/new">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="Y">destinationPath</parameter>
                    <parameter datatype="string" required="Y">confFiles</parameter>
                    <parameter datatype="string" required="Y">variableList</parameter>
                    <parameter datatype="string" required="Y">startScript</parameter>
                    <parameter datatype="string" required="N">args</parameter>
                </input-parameters>
                <output-parameters>
                    <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
            <interface action="update" path="/wecube-plugins-saltstack/v1/deploy/apply-deployment/update">
                <input-parameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">confFiles</parameter>
                    <parameter datatype="string" required="Y">destinationPath</parameter>
                    <parameter datatype="string" required="Y">variableList</parameter>
                    <parameter datatype="string" required="Y">stopScript</parameter>
                    <parameter datatype="string" required="Y">startScript</parameter>
                    <parameter datatype="string" required="N">args</parameter>
                </input-parameters>
                <output-parameters>
                <parameter datatype="string">guid</parameter>
                </output-parameters>
            </interface>
        </plugin>
    </plugins>
</package>