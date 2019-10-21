<?xml version="1.0" encoding="UTF-8"?>
<package name="salt-stack-deployment" version="{{PLUGIN_VERSION}}">
    <docker-image-file>wecube-plugins-saltstack.tar</docker-image-file>
    <docker-image-repository>wecube-plugins-saltstack</docker-image-repository>
    <docker-image-tag>{{IMAGE_TAG}}</docker-image-tag>
    <container-port>8082</container-port>
    <container-config-directory>/home/app/wecube-plugins-saltstack/conf</container-config-directory>
    <container-log-directory>/home/app/wecube-plugins-saltstack/logs</container-log-directory>
    <container-start-param>-v /etc/localtime:/etc/localtime -e minion_master_ip={{HOST_IP}} -e minion_passwd=Ab888888 -e minion_port=22 -p 9099:80 -p 9090:8080 -p 4505:4505 -p 4506:4506 --privileged=true  -v /home/app/data/minions_pki:/etc/salt/pki/master/minions -v /home/app/wecube-plugins-saltstack/logs:/home/app/wecube-plugins-saltstack/logs -v /home/app/data:/home/app/data </container-start-param>
    <plugin id="file" name="File Operation" >
        <interface name="copy" path="/v1/deploy/file/copy">
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
                <!--
                <parameter datatype="string">resultCode</parameter>
                <parameter datatype="string">resultMessage</parameter>
                 -->
            </output-parameters>
        </interface>
    </plugin>
	<plugin id="agent" name="Salt-Stack Agent">
        <interface name="install" path="/v1/deploy/agent/install">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">seed</parameter>
                 <parameter datatype="string">password</parameter>
				<parameter datatype="string">host</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">guid</parameter>
                <!--
                <parameter datatype="string">resultCode</parameter>
                <parameter datatype="string">resultMessage</parameter>
                -->
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="variable" name="Variable Operation">
        <interface name="copy" path="/v1/deploy/variable/replace">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">confFiles</parameter>
                <parameter datatype="string">endpoint</parameter>
                <parameter datatype="string">variableList</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">guid</parameter>
                 <parameter datatype="string">s3PkgPath</parameter>
                 <!--
                <parameter datatype="string">detail</parameter>
                -->
            </output-parameters>
        </interface>
    </plugin>
   <plugin id="script" name="Script Operation">
        <interface name="run" path="/v1/deploy/script/run">
            <input-parameters>
	        <parameter datatype="string">guid</parameter>
                <parameter datatype="string">endpointType</parameter>
                <parameter datatype="string">endpoint</parameter>
                <!-- <parameter datatype="string">accessKey</parameter>
                <parameter datatype="string">secretKey</parameter> -->
                <parameter datatype="string">target</parameter>
                <parameter datatype="string">runAs</parameter>
                <parameter datatype="string">args</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">guid</parameter>
            <!--
                <parameter datatype="string">detail</parameter>
            -->
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="user" name="User Management">
        <interface name="add" path="/v1/deploy/user/add">
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
             <!--
                <parameter datatype="string">detail</parameter>
             -->
            </output-parameters>
        </interface>
        <interface name="remove" path="/v1/deploy/user/remove">
            <input-parameters>
                 <parameter datatype="string">guid</parameter>
                 <parameter datatype="string">target</parameter>
                <parameter datatype="string">userName</parameter>
            </input-parameters>
            <output-parameters>
              <!--
                <parameter datatype="string">detail</parameter>
              -->
                <parameter datatype="string">guid</parameter>
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="database" name="Database Operation">
        <interface name="runScript" path="/v1/deploy/database/runScript">
            <input-parameters>
                <parameter datatype="string">endpoint</parameter>
                <!-- <parameter datatype="string">accessKey</parameter>
                <parameter datatype="string">secretKey</parameter> -->
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
                <!--
                <parameter datatype="string">detail</parameter>
                 -->
            </output-parameters>
        </interface>
    </plugin>
    <plugin id="disk" name="Storage Disk Operation" >
        <interface name="getUnformatedDisk" path="/v1/deploy/disk/getUnformatedDisk">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">target</parameter>
            </input-parameters>
            <output-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">unformatedDisks</parameter>
            </output-parameters>
        </interface>
        <interface name="formatAndMountDisk" path="/v1/deploy/disk/formatAndMountDisk">
            <input-parameters>
                <parameter datatype="string">guid</parameter>
                <parameter datatype="string">target</parameter>
                <parameter datatype="string">diskName</parameter>
                <parameter datatype="string">fileSystemType</parameter>
                <parameter datatype="string">mountDir</parameter>
            </input-parameters>

            <output-parameters>
              <parameter datatype="string">guid</parameter>
            <!--
                <parameter datatype="string">detail</parameter>
             -->
            </output-parameters>
        </interface>
    </plugin>
     <plugin id="apply-deployment" name="Apply Deployment Operation" >
        <interface name="new" path="/v1/deploy/apply-deployment/new">
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
        <interface name="update" path="/v1/deploy/apply-deployment/update">
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
</package>
