<?xml version="1.0" encoding="UTF-8"?>
<package name="saltstack" version="{{PLUGIN_VERSION}}">
    <!-- 1.依赖分析 - 描述运行本插件包需要的其他插件包 -->
    <packageDependencies>
        <packageDependency name="wecmdb" version="v1.5.0"/>
        <packageDependency name="platform" version="v2.9.0"/>
    </packageDependencies>

    <!-- 2.菜单注入 - 描述运行本插件包需要注入的菜单 -->
    <menus>
    </menus>

    <!-- 3.数据模型 - 描述本插件包的数据模型,并且描述和Framework数据模型的关系 -->
    <dataModel>
    </dataModel>

    <!-- 4.系统参数 - 描述运行本插件包需要的系统参数 -->
    <systemParameters>
        <systemParameter name="SALTSTACK_SCRIPT_LOCAL" scopeType="global" defaultValue="LOCAL"/>
        <systemParameter name="SALTSTACK_SCRIPT_S3" scopeType="global" defaultValue="S3"/>
        <systemParameter name="SALTSTACK_SCRIPT_USER_PARAM" scopeType="global" defaultValue="USER_PARAM"/>
        <systemParameter name="SALTSTACK_ENCRYPT_VARIBLE_PREFIX" scopeType="global" defaultValue="!,%"/>
        <systemParameter name="SALTSTACK_SINGLE_ENCRYPT_VARIBLE_PREFIX" scopeType="global" defaultValue="+"/>
        <systemParameter name="SALTSTACK_ENCRYPT_ESCAPE_PREFIX" scopeType="global" defaultValue=""/>
        <systemParameter name="SALTSTACK_FILE_VARIBLE_PREFIX" scopeType="global" defaultValue="^"/>
        <systemParameter name="SALTSTACK_SYSTEM_PRIVATE_KEY" scopeType="global" defaultValue=""/>
        <systemParameter name="SALTSTACK_AGENT_USER" scopeType="global" defaultValue="root"/>
        <systemParameter name="SALTSTACK_AGENT_PORT" scopeType="global" defaultValue="22"/>
        <systemParameter name="SALTSTACK_DEFAULT_SPECIAL_REPLACE" scopeType="global" defaultValue="@,#"/>
        <systemParameter name="SALTSTACK_AGENT_INSTALL_YUM" scopeType="global" defaultValue="yum"/>
        <systemParameter name="SALTSTACK_AGENT_INSTALL_LOCAL" scopeType="global" defaultValue="local"/>
        <systemParameter name="SALTSTACK_LOG_LEVEL" scopeType="global" defaultValue="info"/>
        <systemParameter name="SALTSTACK_RESET_ENV" scopeType="global" defaultValue="Y"/>
        <systemParameter name="SALTSTACK_API_CONCURRENT_NUM" scopeType="global" defaultValue="5"/>
        <systemParameter name="SALTSTACK_VARIABLE_NULL_CHECK" scopeType="global" defaultValue="N"/>
        <systemParameter name="SALTSTACK_AUTH_TYPE_KEY" scopeType="global" defaultValue="key"/>
        <systemParameter name="SALTSTACK_AUTH_KEY_NAME" scopeType="global" defaultValue=""/>
    </systemParameters>

    <!-- 5.权限设定 -->
    <authorities>
    </authorities>

    <!-- 6.运行资源 - 描述部署运行本插件包需要的基础资源(如主机、虚拟机、容器、数据库等) -->
    <resourceDependencies>
        <docker imageName="{{IMAGENAME}}" containerName="{{CONTAINERNAME}}" portBindings="9099:80,4505:4505,4506:4506,4507:4507,{{PORTBINDING}}" volumeBindings="/etc/localtime:/etc/localtime,{{BASE_MOUNT_PATH}}/data/minions_pki:/etc/salt/pki/master/minions,{{BASE_MOUNT_PATH}}/saltstack/logs:/home/app/wecube-plugins-saltstack/logs,{{BASE_MOUNT_PATH}}/data:/home/app/data,{{BASE_MOUNT_PATH}}/certs:/data/certs" envVariables="minion_master_ip={{ALLOCATE_HOST}},minion_port={{SALTSTACK_AGENT_PORT}},DEFAULT_S3_KEY={{S3_ACCESS_KEY}},DEFAULT_S3_PASSWORD={{S3_SECRET_KEY}},SALTSTACK_DEFAULT_SPECIAL_REPLACE={{SALTSTACK_DEFAULT_SPECIAL_REPLACE}},CORE_ADDR={{CORE_ADDR}},GATEWAY_URL={{GATEWAY_URL}},SALTSTACK_ENCRYPT_VARIBLE_PREFIX={{SALTSTACK_ENCRYPT_VARIBLE_PREFIX}},SALTSTACK_SINGLE_ENCRYPT_VARIBLE_PREFIX={{SALTSTACK_SINGLE_ENCRYPT_VARIBLE_PREFIX}},SALTSTACK_ENCRYPT_ESCAPE_PREFIX={{SALTSTACK_ENCRYPT_ESCAPE_PREFIX}},SALTSTACK_FILE_VARIBLE_PREFIX={{SALTSTACK_FILE_VARIBLE_PREFIX}},SALTSTACK_LOG_LEVEL={{SALTSTACK_LOG_LEVEL}},JWT_SIGNING_KEY={{JWT_SIGNING_KEY}},S3_SERVER_URL={{S3_SERVER_URL}},SUB_SYSTEM_CODE={{SUB_SYSTEM_CODE}},SUB_SYSTEM_KEY={{SUB_SYSTEM_KEY}},SALTSTACK_RESET_ENV={{SALTSTACK_RESET_ENV}},API_CONCURRENT_NUM={{SALTSTACK_API_CONCURRENT_NUM}},SALTSTACK_VARIABLE_NULL_CHECK={{SALTSTACK_VARIABLE_NULL_CHECK}},ENCRYPT_SEED={{ENCRYPT_SEED}}"/>
        <s3 bucketName="salt-tmp"/>
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
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">command</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_AGENT_INSTALL_LOCAL" >method</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">passwordType</parameter>
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
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">user</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">passwordType</parameter>
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
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="FILE_VARIBLE_PREFIX">fileReplacePrefix</parameter>
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
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_SCRIPT_LOCAL">endpointType</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">workDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
                </outputParameters>
            </interface>
            <interface action="ssh-run" path="/saltstack/v1/host-script/ssh-run" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_SCRIPT_LOCAL">endpointType</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">workDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
                </outputParameters>
            </interface>
            <interface action="run-local" path="/saltstack/v1/host-script/run-local" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_SCRIPT_LOCAL">endpointType</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">workDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
                </outputParameters>
            </interface>
            <interface action="run-s3" path="/saltstack/v1/host-script/run-s3" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_SCRIPT_S3">endpointType</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">endpoint</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">workDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
                </outputParameters>
            </interface>
            <interface action="run-user-param" path="/saltstack/v1/host-script/run-user-param" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="system_variable" mappingSystemVariableName="SALTSTACK_SCRIPT_USER_PARAM">endpointType</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">scriptContent</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">runAs</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">args</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">workDir</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
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
            <interface action="create" path="/saltstack/v1/host-file/create" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">destinationPath</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">fileContent</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">fileOwner</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="find" path="/saltstack/v1/host-file/find" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">filePath</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">filePattern</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">files</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/host-file/delete" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">baseDir</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">targetPath</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="host-info" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="query" path="/saltstack/v1/host-info/query" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="int" sensitiveData="N" mappingType="entity" mappingEntityExpression="">numCpus</parameter>
                    <parameter datatype="int" sensitiveData="N" mappingType="entity" mappingEntityExpression="">memTotal</parameter>
                    <parameter datatype="int" sensitiveData="N" mappingType="entity" mappingEntityExpression="">diskTotal</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">os</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">kernel</parameter>
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
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">rwDir</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">rwFile</parameter>
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
            <interface action="password" path="/saltstack/v1/host-user/password" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">target</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="constant">userName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="constant">password</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="context">password</parameter>
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
            <interface action="run-custom-sql" path="/saltstack/v1/mysql-script/run-custom-sql" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">sql</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">detail</parameter>
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
            <interface action="change-password" path="/saltstack/v1/mysql-user/change-password" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">databaseUserGuid</parameter>
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
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">appPublicKey</parameter>
                    <parameter datatype="string" required="N" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="SYSTEM_PRIVATE_KEY" >sysPrivateKey</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">rwDir</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">rwFile</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileTrade</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileTrace</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileMetric</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileKeyword</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">signFileSrc</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">signFileDst</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">s3PkgPath</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">fileDetail</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileTrade</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileTrace</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileMetric</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileKeyword</parameter>
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
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileTrade</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileTrace</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileMetric</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="constant">logFileKeyword</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">signFileSrc</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">signFileDst</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">clearPath</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">s3PkgPath</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">fileDetail</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileTrade</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileTrace</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileMetric</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">logFileKeyword</parameter>
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
        <plugin name="password" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="encode" path="/saltstack/v1/password/encode" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="decode" path="/saltstack/v1/password/decode" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">password</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">password</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
             <interface action="sshkeygen" path="/saltstack/v1/password/sshkeygen" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">keyName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">keyName</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="context">privateKey</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">publicKey</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
        <plugin name="redis-user" targetPackage="" targetEntity="" registerName="" targetEntityFilterRule="">
            <interface action="add" path="/saltstack/v1/redis-user/add" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">adminUserName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">adminPassword</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userGuid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">userPassword</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userReadKeyPattern</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userWriteKeyPattern</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userGuid</parameter>
                    <parameter datatype="string" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">userPassword</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
            <interface action="delete" path="/saltstack/v1/redis-user/delete" filterRule="">
                <inputParameters>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">guid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="system_variable" mappingSystemVariableName="ENCRYPT_SEED" >seed</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">host</parameter>
                    <parameter datatype="string" required="N" sensitiveData="N" mappingType="entity" mappingEntityExpression="">adminUserName</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="Y" mappingType="entity" mappingEntityExpression="">adminPassword</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">port</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userGuid</parameter>
                    <parameter datatype="string" required="Y" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userName</parameter>
                </inputParameters>
                <outputParameters>
                    <parameter datatype="string" sensitiveData="N" mappingType="entity" mappingEntityExpression="">userGuid</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorCode</parameter>
                    <parameter datatype="string" sensitiveData="N" mappingType="context">errorMessage</parameter>
                </outputParameters>
            </interface>
        </plugin>
    </plugins>
</package>
