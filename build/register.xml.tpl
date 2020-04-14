<?xml version="1.0" encoding="UTF-8"?>
<package name="saltstack" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
        <packageDependency name="wecmdb" version="v1.4.0"/>
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
            <systemParameter name="SALTSTACK_AGENT_PORT" scopeType="global" defaultValue="22"/>
            <systemParameter name="SALTSTACK_PASSWORD" scopeType="plugins" defaultValue="PA888888"/>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="9099:80,9090:8080,4505:4505,4506:4506,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/data/minions_pki:/etc/salt/pki/master/minions,{{BASE_MOUNT_PATH}}/saltstack/logs:/home/app/wecube-plugins-saltstack/logs,/{{BASE_MOUNT_PATH}}/data:/home/app/data" envVariables="minion_master_ip={{ALLOCATE_HOST}},minion_passwd={{SALTSTACK_PASSWORD}},minion_port={{SALTSTACK_AGENT_PORT}},DEFAULT_S3_KEY={{S3_ACCESS_KEY}},DEFAULT_S3_PASSWORD={{S3_SECRET_KEY}}"/>
    </resourceDependencies>

    <!-- 7.插件列表 - 描述插件包中单个插件的输入和输出 -->
    <plugins>
            <plugin name="agent" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="install" path="/saltstack/v1/agent/install" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                 <interface action="uninstall" path="/saltstack/v1/agent/uninstall" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="package-variable" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="replace" path="/saltstack/v1/package-variable/replace" filterRule="">
                    <inputParameters>
    		            <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">s3PkgPath</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-script" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="run" path="/saltstack/v1/host-script/run" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-file" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="copy" path="/saltstack/v1/host-file/copy" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">destinationPath</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">unpack</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">fileOwner</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-user" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/host-user/add" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">userName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="constant">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">userGroup</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">userId</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">groupId</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">homeDir</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="context">password</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/host-user/delete" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">userName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="mysql-script" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="run" path="/saltstack/v1/mysql-script/run" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">sql_files</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="mysql-database" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/mysql-database/add" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseOwnerName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">databaseOwnerPassword</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">databaseOwnerPassword</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/mysql-database/delete" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseOwnerGuid</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="mysql-user" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/mysql-user/add" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">databaseUserPassword</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">databaseUserPassword</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/mysql-user/delete" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="apply-deployment" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
                <interface action="new" path="/saltstack/v1/apply-deployment/new" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">startScript</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="update" path="/saltstack/v1/apply-deployment/update" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">stopScript</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">startScript</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/apply-deployment/delete" filterRule="">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">stopScript</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>



             <!--最佳实践 -->
            <plugin name="agent" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host" targetEntityFilterRule="">
                <interface action="install" path="/saltstack/v1/agent/install" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_name">user</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                 <interface action="uninstall" path="/saltstack/v1/agent/uninstall" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="package-variable" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_instance" targetEntityFilterRule="">
                <interface action="create-replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
    		            <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package_url">s3PkgPath</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="change-replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
    		            <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package_url">s3PkgPath</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="package-variable" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="rdb_instance" targetEntityFilterRule="">
                <interface action="create-replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
    		            <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">s3PkgPath</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                 <interface action="change-replace" path="/saltstack/v1/package-variable/replace" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
    		            <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">s3PkgPath</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-script" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy" targetEntityFilterRule="">
                <interface action="run-start-script" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'startup'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.start_script">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-start-script-for-change" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.start_script">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-stop-script" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'stoped'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-stop-script-for-change" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-script" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host" targetEntityFilterRule="">
                <interface action="install-monitor-agent" path="/saltstack/v1/host-script/run" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_INSTALL_SCRIPT">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-custom-script" path="/saltstack/v1/host-script/run">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_USER_PARAM">endpointType</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">scriptContent</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">runAs</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">args</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-file" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy" targetEntityFilterRule="">
                <interface action="copy" path="/saltstack/v1/host-file/copy" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path">destinationPath</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.is_decompression">unpack</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">fileOwner</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-file" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="monitor_agent" targetEntityFilterRule="">
                <interface action="copy" path="/saltstack/v1/host-file/copy" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_S3_PATH">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_UPLOAD_PATH">destinationPath</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_UNPACK">unpack</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.NONE">fileOwner</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="host-user" targetPackage="wecmdb" targetEntity="host_resource_instance" registerName="host" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/host-user/add">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">userName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="constant">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">userGroup</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">userId</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">groupId</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">homeDir</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="context">password</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/host-user/delete">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">userName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:host_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="mysql-script" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="db_deploy" targetEntityFilterRule="">
                <interface action="run-deploy-script" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.deploy_file_path">sql_files</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.code">databaseName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-upgrade-script" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.start_file_path">sql_files</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.code">databaseName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="run-rollback-script" path="/saltstack/v1/mysql-script/run" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_package>wecmdb:deploy_package.stop_file_path">sql_files</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password">password</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.code">databaseName</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="mysql-database" targetPackage="wecmdb" targetEntity="rdb_instance" registerName="db_deploy" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/mysql-database/add" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.code">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user">databaseOwnerName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password">databaseOwnerPassword</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.deploy_user_password">databaseOwnerPassword</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/mysql-database/delete" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_name">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.rdb_resource_instance>wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.unit>wecmdb:unit.code">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_instance.guid">databaseOwnerGuid</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>

            <plugin name="mysql-user" targetPackage="wecmdb" targetEntity="rdb_resource_instance" registerName="db_monitor" targetEntityFilterRule="">
                <interface action="add" path="/saltstack/v1/mysql-user/add" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_DATABASE">databaseName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER">databaseUserName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_PWD">databaseUserPassword</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="context">databaseUserPassword</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/mysql-user/delete" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.intranet_ip>wecmdb:ip_address.code">host</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_name">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.login_port">port</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER" >databaseUserName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:rdb_resource_instance.guid">databaseUserGuid</parameter>
    		            <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
            <plugin name="apply-deployment" targetPackage="wecmdb" targetEntity="app_instance" registerName="app_deploy" targetEntityFilterRule="">
                <interface action="new" path="/saltstack/v1/apply-deployment/new" filterRule="{state_code eq 'created'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password">password</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_script">startScript</parameter>
                        <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user_password">password</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="update" path="/saltstack/v1/apply-deployment/update" filterRule="{state_code eq 'changed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url">endpoint</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file">confFiles</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.variable_values">variableList</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script">stopScript</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_script">startScript</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.NONE">args</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX">encryptVariblePrefix</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.public_key">appPublicKey</parameter>
                        <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
                <interface action="delete" path="/saltstack/v1/apply-deployment/delete" filterRule="{state_code eq 'destroyed'}{fixed_date is NULL}">
                    <inputParameters>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.host_resource_instance>wecmdb:host_resource_instance.intranet_ip>wecmdb:ip_address.code">target</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.deploy_user">userName</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.unit>wecmdb:unit.deploy_path">destinationPath</parameter>
                        <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.stop_script">stopScript</parameter>
                    </inputParameters>
                    <outputParameters>
                        <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="wecmdb:app_instance.guid">guid</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                        <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    </outputParameters>
                </interface>
            </plugin>
        </plugins>
</package>
