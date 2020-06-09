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
curl -O http://$master_ip:9099/salt-minion/minion_install_pkg.tar.gz /tmp/salt/minion_install_pkg.tar.gz
cd /tmp/salt/
tar zxf minion_install_pkg.tar.gz
cd minion_install_pkg && ./install_minion.sh
curl -O http://$master_ip:9099/salt-minion/conf/minion /tmp/salt/minion
sed -i "s~{{ minion_id }}~$minion_ip~g" /tmp/salt/minion
mv /tmp/salt/minion /etc/salt/minion
systemctl enable salt-minion
systemctl start salt-minion
echo "install salt-minion done"