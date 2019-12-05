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
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/data/minions_pki:/etc/salt/pki/master/minions,{{BASE_MOUNT_PATH}}/saltstack/logs:/home/app/saltstack/logs,/{{BASE_MOUNT_PATH}}/data:/home/app/data" envVariables="minion_master_ip={{ALLOCATE_HOST}},minion_passwd=Ab888888,minion_port=22"/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="file">
            <interface action="copy" path="/saltstack/v1/file/copy">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">unpack</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="agent">
            <interface action="install" path="/saltstack/v1/agent/install">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="variable">
            <interface action="copy" path="/saltstack/v1/variable/replace">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">variableList</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>s3PkgPath</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="script">
            <interface action="run" path="/saltstack/v1/script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">runAs</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="user">
            <interface action="add" path="/saltstack/v1/user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">userId</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">groupId</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>password</parameter>
                </outputParameters>
            </interface>
            <interface action="remove" path="/saltstack/v1/user/remove">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="database">
            <interface action="runScript" path="/saltstack/v1/database/runScript">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">port</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>

              <interface action="addDatabase" path="/saltstack/v1/database/addDatabase">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">host</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">password</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">port</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">databaseName</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">databaseOwnerName</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">databaseOwnerPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType='context'>databaseOwnerPassword</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="disk">
            <interface action="getUnformatedDisk" path="/saltstack/v1/disk/getUnformatedDisk">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                    <parameter datatype="string" mappingType='context'>unformatedDisks</parameter>
                </outputParameters>
            </interface>
            <interface action="formatAndMountDisk" path="/saltstack/v1/disk/formatAndMountDisk">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">diskName</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">fileSystemType</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">mountDir</parameter>
                </inputParameters>

                <outputParameters>
                <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface action="new" path="/saltstack/v1/apply-deployment/new">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/apply-deployment/update">
                <inputParameters>
                    <parameter datatype="string" mappingType='entity' required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">target</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">userName</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">stopScript</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                <parameter datatype="string" mappingType='context'>guid</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>