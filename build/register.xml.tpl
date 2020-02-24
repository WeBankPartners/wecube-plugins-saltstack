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
        <plugin name="agent" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="install" path="/saltstack/v1/agent/install">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
             <interface action="uninstall" path="/saltstack/v1/agent/uninstall">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="package-variable" targetPackage="wecmdb" targetEntity="business_app_instance">
            <interface action="replace" path="/saltstack/v1/package-variable/replace">
                <inputParameters>
		            <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url">s3PkgPath</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-script" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="app_deploy">
            <interface action="runDeployScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStartScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.start_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runStopScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.stop_script" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-script" targetPackage="wecmdb" targetEntity="resource_instance">
            <interface action="installMonitorAgent" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_LOCAL" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runCustomScript" path="/saltstack/v1/host-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SCRIPT_END_POINT_TYPE_USER_PARAM" required="Y">endpointType</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">scriptContent</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">runAs</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">args</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-file" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="app_deploy">
            <interface action="copy-package" path="/saltstack/v1/host-file/copy">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.is_decompression" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-file" targetPackage="wecmdb" targetEntity="resource_instance" registerName="monitor_agent">
            <interface action="copy-agent" path="/saltstack/v1/host-file/copy">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="HOST_EXPORTER_S3_PATH" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">unpack</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.NONE" required="N">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-user" targetPackage="wecmdb" targetEntity="resource_instance" registerName="host">
            <interface action="add" path="/saltstack/v1/host-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">password</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="constant" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/host-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="constant" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-user" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="app_deploy">
            <interface action="add" path="/saltstack/v1/host-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="N">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.system>wecmdb:system.system_design>wecmdb:system_design.code" required="N">userGroup</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.subsys_design>wecmdb:subsys_design.subsys_design_id" required="N">userId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.subsys>wecmdb:subsys.system>wecmdb:system.system_design>wecmdb:system_design.system_design_id" required="N">groupId</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.deploy_path" required="N">homeDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">password</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="remove" path="/saltstack/v1/host-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-script" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="db_deploy">
            <interface action="runDeployScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runUpgradeScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.start_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="runrRollbackScript" path="/saltstack/v1/mysql-script/run">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.stop_file_path" required="N">sql_files</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="N">databaseName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-database" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="db_deploy">
            <interface action="addDatabase" path="/saltstack/v1/mysql-database/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">databaseOwnerName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="N">databaseOwnerPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseOwnerGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">databaseOwnerPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteDatabase" path="/saltstack/v1/mysql-database/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">databaseOwnerGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseOwnerGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-user" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="db_deploy">
            <interface action="addUser" path="/saltstack/v1/mysql-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.code" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password" required="N">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user_password">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteUser" path="/saltstack/v1/mysql-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.guid" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">databaseUserGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">databaseUserGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="mysql-user" targetPackage="wecmdb" targetEntity="resource_instance" registerName="mysql_db_monitor">
            <interface action="addMonitorUser" path="/saltstack/v1/mysql-user/add">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_DATABASE" required="Y">databaseName</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_PWD" required="N">databaseUserPassword</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">databaseUserGuid</parameter>
                    <parameter datatype="string" mappingType="context">databaseUserPassword</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="deleteMonitorUser" path="/saltstack/v1/mysql-user/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">host</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_name" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.user_password" required="Y">password</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.login_port" required="N">port</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="MYSQL_MONITOR_USER" required="Y">databaseUserName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id" required="Y">databaseUserGuid</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:resource_instance.id">databaseUserGuid</parameter>
		            <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="apply-deployment" targetPackage="wecmdb" targetEntity="business_app_instance" registerName="app_deploy">
            <interface action="new" path="/saltstack/v1/apply-deployment/new">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_script" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="update" path="/saltstack/v1/apply-deployment/update">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.deploy_package_url" required="Y">endpoint</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_package>wecmdb:deploy_package.diff_conf_file" required="Y">confFiles</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.variable_values" required="Y">variableList</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.stop_script" required="Y">stopScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_script" required="Y">startScript</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.NONE" required="N">args</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_VARIBLE_PREFIX" required="Y">encryptVariblePrefix</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" required="Y">seed</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.public_key" required="Y">appPublicKey</parameter>
                    <parameter datatype="string" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" required="N">sysPrivateKey</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/apply-deployment/delete">
                <inputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id" required="Y">guid</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.resource_instance>wecmdb:resource_instance.intranet_ip>wecmdb:ip_address.code" required="Y">target</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.deploy_user" required="Y">userName</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.unit>wecmdb:unit.deploy_path" required="Y">destinationPath</parameter>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.stop_script" required="Y">stopScript</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" mappingType="entity" mappingEntityExpression="wecmdb:business_app_instance.id">guid</parameter>
                    <parameter datatype="string" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
