<?xml version="1.0" encoding="UTF-8"?>
<package name="saltstack" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
        <packageDependency name="wecmdb" version="v1.4.0"/>
        <packageDependency name="qcloud" version="v1.8.0"/>
    </packageDependencies>

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单 -->
    <menus>
    </menus>

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系 -->
    <dataModel>
    </dataModel>

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数 -->
    <systemParameters>
        <systemParameter name="SCRIPT_END_POINT_TYPE_LOCAL" scopeType="global" defaultValue="LOCAL"/>
        <systemParameter name="SCRIPT_END_POINT_TYPE_S3" scopeType="global" defaultValue="S3"/>
        <systemParameter name="SCRIPT_END_POINT_TYPE_USER_PARAM" scopeType="global" defaultValue="USER_PARAM"/>
        <systemParameter name="ENCRYPT_VARIBLE_PREFIX" scopeType="global" defaultValue="!"/>
        <systemParameter name="SYSTEM_PRIVATE_KEY" scopeType="global" defaultValue=""/>
        <systemParameter name="SALTSTACK_SERVER_URL" scopeType="global" defaultValue="http://127.0.0.1:20002"/>
	<systemParameter name="SALTSTACK_AGENT_USER" scopeType="global" defaultValue="root"/>
        <systemParameter name="SALTSTACK_AGENT_PORT" scopeType="global" defaultValue="9000"/>
        <systemParameter name="SALTSTACK_PASSWORD" scopeType="plugins" defaultValue="WB888888"/>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/data/minions_pki:/etc/salt/pki/master/minions,{{BASE_MOUNT_PATH}}/saltstack/logs:/home/app/wecube-plugins-saltstack/logs,/{{BASE_MOUNT_PATH}}/data:/home/app/data" envVariables="minion_master_ip={{ALLOCATE_HOST}},minion_passwd={{SALTSTACK_PASSWORD}},minion_port={{SALTSTACK_AGENT_PORT}}"/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
        <plugin name="agent">
            <interface action="install" path="/saltstack/v1/agent/install">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">user</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
             <interface action="uninstall" path="/saltstack/v1/agent/uninstall">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="package-variable">
            <interface action="replace" path="/saltstack/v1/package-variable/replace">
                <inputParameters>
		            <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">s3PkgPath</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-script">
            <interface action="runDeployScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStartScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStopScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-file">
            <interface action="copy-package" path="/saltstack/v1/host-file/copy">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-user">
            <interface action="add" path="/saltstack/v1/host-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="constant" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/host-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-script">
            <interface action="runDeployScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runUpgradeScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runrRollbackScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-database">
            <interface action="addDatabase" path="/saltstack/v1/mysql-database/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseOwnerName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="Y">databaseOwnerPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">databaseOwnerPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteDatabase" path="/saltstack/v1/mysql-database/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseOwnerGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-user">
            <interface action="addUser" path="/saltstack/v1/mysql-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="Y">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteUser" path="/saltstack/v1/mysql-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">databaseUserGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface action="new" path="/saltstack/v1/apply-deployment/new">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/apply-deployment/update">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">stopScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/apply-deployment/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="" required="Y">stopScript</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>


        
         <!--最佳实践 -->
        <plugin name="agent" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host">
            <interface action="install" path="/saltstack/v1/agent/install" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_AGENT_PORT" required="N">port</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_AGENT_USER" required="N">user</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
             <interface action="uninstall" path="/saltstack/v1/agent/uninstall" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="package-variable" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_instance">
            <interface action="replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code in ['created','changed']}{fixed_date eq ''}">
                <inputParameters>
		            <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package_url">s3PkgPath</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="package-variable" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="rdb_instance">
            <interface action="replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code in ['created','changed']}{fixed_date eq ''}">
                <inputParameters>
		            <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">s3PkgPath</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-script" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy">
            <interface action="runDeployScript" path="/saltstack/v1/host-script/run" filterRule="{state_code in ['created','changed']}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStartScript" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'startup'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.start_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStopScript" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'stoped'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-script" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host">
            <interface action="installMonitorAgent" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runCustomScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_USER_PARAM" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-file" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy">
            <interface action="copy-package" path="/saltstack/v1/host-file/copy" filterRule="{state_code in ['created','changed']}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.is_decompression" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-file" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="monitor_agent">
            <interface action="copy-agent" path="/saltstack/v1/host-file/copy" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_S3_PATH" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-user" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host">
            <interface action="add" path="/saltstack/v1/host-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="constant" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/host-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-user" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy">
            <interface action="add" path="/saltstack/v1/host-user/add" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.system_design>wecmdb:system_design.code" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.subsys_design>wecmdb:subsys_design.subsys_design_id" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.app_system>wecmdb:app_system.system_design>wecmdb:system_design.system_design_id" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove" path="/saltstack/v1/host-user/delete" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-script" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="db_deploy">
            <interface action="runDeployScript" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.deploy_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runUpgradeScript" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'changed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.start_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runrRollbackScript" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'changed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.stop_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-database" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="db_deploy">
            <interface action="addDatabase" path="/saltstack/v1/mysql-database/add" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">databaseOwnerName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" required="N" sensitiveData="Y">databaseOwnerPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" sensitiveData="Y">databaseOwnerPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteDatabase" path="/saltstack/v1/mysql-database/delete" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">databaseOwnerGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-user" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="db_deploy">
            <interface action="addUser" path="/saltstack/v1/mysql-user/add" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" required="N" sensitiveData="Y">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password" sensitiveData="Y">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteUser" path="/saltstack/v1/mysql-user/delete" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid" required="Y">databaseUserGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseUserGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-user" targetPackage="wecmdb" targetEntity="rdb_resource_instance" registerName="mysql_db_monitor">
            <interface action="addMonitorUser" path="/saltstack/v1/mysql-user/add" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_DATABASE" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_PWD" required="N" sensitiveData="Y">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="context">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteMonitorUser" path="/saltstack/v1/mysql-user/delete" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password" required="Y" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid" required="Y">databaseUserGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy">
            <interface action="new" path="/saltstack/v1/apply-deployment/new" filterRule="{state_code eq 'created'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password" required="N" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_script" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password" sensitiveData="Y">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/apply-deployment/update" filterRule="{state_code eq 'changed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script" required="Y">stopScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_script" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/apply-deployment/delete" filterRule="{state_code eq 'destroyed'}{fixed_date eq ''}">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script" required="Y">stopScript</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
