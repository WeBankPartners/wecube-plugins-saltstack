<?xml version="1.0" encoding="UTF-8"?>
<package name="saltstack" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
    </packageDependencies>

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单 -->
    <menus>
    </menus>

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系 -->
    <dataModel>
    </dataModel>

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数 -->
    <systemParameters>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{IMAGENAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,{{base_mount_path}}/data/minions_pki:/etc/salt/pki/master/minions,{{base_mount_path}}/saltstack/logs:/home/app/saltstack/logs,/{{base_mount_path}}/data:/home/app/data" envVariables="minion_master_ip={{minion_master_ip}},minion_passwd=Ab888888,minion_port=22"/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="file">
            <interface action="copy" path="/saltstack/v1/deploy/file/copy">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">destinationPath</parameter>
                    <parameter datatype="string" required="N">unpack</parameter>
                    <parameter datatype="string" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="agent">
            <interface action="install" path="/saltstack/v1/deploy/agent/install">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">seed</parameter>
                    <parameter datatype="string" required="Y">password</parameter>
                    <parameter datatype="string" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="variable">
            <interface action="copy" path="/saltstack/v1/deploy/variable/replace">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">confFiles</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">variableList</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">s3PkgPath</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="script">
            <interface action="run" path="/saltstack/v1/deploy/script/run">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">endpointType</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="N">runAs</parameter>
                    <parameter datatype="string" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="user">
            <interface action="add" path="/saltstack/v1/deploy/user/add">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="N">password</parameter>
                    <parameter datatype="string" required="N">userGroup</parameter>
                    <parameter datatype="string" required="N">userId</parameter>
                    <parameter datatype="string" required="N">groupId</parameter>
                    <parameter datatype="string" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
            <interface action="remove" path="/saltstack/v1/deploy/user/remove">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="database">
            <interface action="runScript" path="/saltstack/v1/deploy/database/runScript">
                <inputParameters>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">seed</parameter>
                    <parameter datatype="string" required="Y">host</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="Y">password</parameter>
                    <parameter datatype="string" required="N">port</parameter>
                    <parameter datatype="string" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="disk">
            <interface action="getUnformatedDisk" path="/saltstack/v1/deploy/disk/getUnformatedDisk">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                    <parameter datatype="string">unformatedDisks</parameter>
                </outputParameters>
            </interface>
            <interface action="formatAndMountDisk" path="/saltstack/v1/deploy/disk/formatAndMountDisk">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">diskName</parameter>
                    <parameter datatype="string" required="Y">fileSystemType</parameter>
                    <parameter datatype="string" required="Y">mountDir</parameter>
                </inputParameters>

                <outputParameters>
                <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface action="new" path="/saltstack/v1/deploy/apply-deployment/new">
                <inputParameters>
                    <parameter datatype="string" required="Y">guid</parameter>
                    <parameter datatype="string" required="Y">target</parameter>
                    <parameter datatype="string" required="Y">endpoint</parameter>
                    <parameter datatype="string" required="Y">userName</parameter>
                    <parameter datatype="string" required="Y">destinationPath</parameter>
                    <parameter datatype="string" required="Y">confFiles</parameter>
                    <parameter datatype="string" required="Y">variableList</parameter>
                    <parameter datatype="string" required="Y">startScript</parameter>
                    <parameter datatype="string" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/deploy/apply-deployment/update">
                <inputParameters>
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
                </inputParameters>
                <outputParameters>
                <parameter datatype="string">guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>