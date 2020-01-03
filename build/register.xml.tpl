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
    	<systemParameter name="ENDPOINT_TYPE_S3" defaultValue="S3" scopeType="saltstack" />
    	<systemParameter name="ENDPOINT_TYPE_LOCAL" defaultValue="LOCAL" scopeType="saltstack" />
    	<systemParameter name="ENDPOINT_TYPE_USER_PARAM" defaultValue="USER_PARAM" scopeType="saltstack" />
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
        <plugin name="file" targetPackage="wecmdb" targetEntity="business_app_instance">
            <interface action="copy" path="/saltstack/v1/file/copy">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.is_decompression" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="agent" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="install" path="/saltstack/v1/agent/install">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="variable" targetPackage="wecmdb" targetEntity="business_app_instance">
            <interface action="copy" path="/saltstack/v1/variable/replace">
                <inputParameters>
		            <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url">s3PkgPath</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="script" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="run" path="/saltstack/v1/script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName='ENDPOINT_TYPE_USER_PARAM' required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.resource_set>wecmdb:resource_set.init_script" required="N">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="user" targetPackage="wecmdb" targetEntity="business_app_instance">
            <interface action="add" path="/saltstack/v1/user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="N">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.group_name" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.user_id" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.group_id" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.home_path" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">password</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove" path="/saltstack/v1/user/remove">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="database" targetPackage="wecmdb" targetEntity="business_app_instance">
            <interface action="runScript" path="/saltstack/v1/database/runScript">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_file_path" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="addDatabase" path="/saltstack/v1/database/addDatabase">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="N">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="N">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="N">databaseOwnerName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="N">databaseOwnerPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">databaseOwnerPassword</parameter>
		            <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
             <interface action="addUser" path="/saltstack/v1/database/addUser">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" required="N">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment">
            <interface action="new" path="/saltstack/v1/apply-deployment/new">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.absolute_deploy_file_path" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">appPublicKey</parameter>
                    <parameter datatype="string" mappingType='entity' required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/apply-deployment/update">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_user>wecmdb:deploy_user.name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.absolute_stop_file_path" required="Y">stopScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.absolute_deploy_file_path" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType='system_variable' mappingSystemVariableName="encrypt_seed" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType='entity' required="Y">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType='context'>errorCode</parameter>
                    <parameter datatype="string" mappingType='context'>errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
