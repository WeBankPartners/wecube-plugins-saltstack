#!/bin/bash

master_ip=''
minion_ip=''
if [ $1 ]
then
  master_ip=$1
else
  echo "param illegal"
  exit 1
fi
if [ $2 ]
then
  minion_ip=$2
else
  echo "param illegal"
  exit 1
fi

mkdir -p /tmp/salt
cd /tmp/salt/
if [ "$3" = "yum" ]
then
  yum install -y salt-minion
else
  curl -O http://$master_ip:9099/salt-minion/minion_install_pkg.tar.gz
  tar zxf minion_install_pkg.tar.gz
  cd minion_install_pkg && ./install_minion.sh
fi
cd /tmp/salt/
curl -O http://$master_ip:9099/salt-minion/conf/minion
sed -i "s~{{ minion_id }}~$minion_ip~g" /tmp/salt/minion
mv /tmp/salt/minion /etc/salt/minion
systemctl enable salt-minion
psout=`ps aux|grep salt-minion|grep -v 'grep'`
if [ -n "$psout" ]
then
  systemctl restart salt-minion
else
  systemctl start salt-minion
fi
sleep 1
is_success=`ps aux|grep salt-minion|grep -v 'grep'`
if [ -n "$is_success" ]
then
  echo "start salt-minion_success"
else
  echo "start salt-minion_fail"
fi