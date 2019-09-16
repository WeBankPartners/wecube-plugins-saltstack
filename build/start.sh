#!/bin/sh
set -e

#chkconfig: 2345 20 80
#description: wecube-plugins-saltstack
TEMPLATE_DIR="/conf/template/"
function replaceFiles()
{
    dest_paths=($APP_HOME"/scripts/salt/install_minion.sh" $APP_HOME"/scripts/salt/remove_master_unused_key.sh" "/srv/salt/minions/conf/minion" "/srv/salt/minions/yum.repos.d/salt-repo.repo")

    for file in ${dest_paths[@]};
    do
        echo $file
        basename=${file##*/}
        sed -i "s/{{minion_port}}/$minion_port/g" $TEMPLATE_DIR$basename".tpl"
        sed -i "s/{{minion_passwd}}/$minion_passwd/g" $TEMPLATE_DIR$basename".tpl"
        sed -i "s/{{minion_master_ip}}/$minion_master_ip/g" $TEMPLATE_DIR$basename".tpl"
        cp $TEMPLATE_DIR$basename".tpl"  $file
    done
}

minion_port=${minion_port-22}
minion_passwd=${minion_passwd-Ab888888}

if [ ! $minion_master_ip ];then
    echo "environment variable minion_master_ip master be set"
    exit -1
fi

runReplaceOkFile="/etc/runReplace"
if [ ! -f $runReplaceOkFile ];then
    replaceFiles
    touch $runReplaceOkFile
fi

rm -rf /var/run/httpd/httpd.pid 
cd /home/app/wecube-plugins-saltstack
mkdir -p logs
./wecube-plugins-saltstack&
/usr/bin/salt-master&
httpd&
/usr/bin/salt-api 
