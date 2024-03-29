#!/bin/sh

#chkconfig: 2345 20 80
#description: wecube-plugins-saltstack
TEMPLATE_DIR="/conf/template/"
function replaceFiles()
{
    dest_paths=($APP_HOME"/scripts/salt/install_minion.sh" $APP_HOME"/scripts/salt/uninstall_minion.sh" $APP_HOME"/scripts/salt/remove_master_unused_key.sh" "/srv/salt/minions/conf/minion" "/srv/salt/minions/yum.repos.d/salt-repo.repo")

    for file in ${dest_paths[@]};
    do
        echo $file
        basename=${file##*/}
        sed -i "s/{{minion_port}}/$minion_port/g" $TEMPLATE_DIR$basename".tpl"
        sed -i "s/{{minion_master_ip}}/$minion_master_ip/g" $TEMPLATE_DIR$basename".tpl"
        cp $TEMPLATE_DIR$basename".tpl"  $file
    done
}

minion_port=${minion_port-22}

if [ ! $minion_master_ip ];then
    echo "environment variable minion_master_ip master be set"
    exit -1
fi

runReplaceOkFile="/etc/runReplace"
if [ ! -f $runReplaceOkFile ];then
    replaceFiles
    touch $runReplaceOkFile
fi

cd /home/app/wecube-plugins-saltstack
sed -i "s/{{SALTSTACK_LOG_LEVEL}}/$SALTSTACK_LOG_LEVEL/g" conf/default.json
mkdir -p logs
mkdir -p /data/minio
mkdir -p /tmp
./wecube-plugins-saltstack&
/usr/bin/salt-master&
/usr/bin/salt-api&

while  /bin/true
do
    process=`ps aux | grep httpd| grep -v grep | awk '{print $1}'`
    if [ -z "$process" ];then
        httpd
    else
       rm -rf /var/run/httpd/httpd.pid 
       sleep 5
    fi
done



