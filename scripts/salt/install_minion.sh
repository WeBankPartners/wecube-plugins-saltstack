#!/bin/bash

if [[ $# < 1 ]]; then
  echo "please input the hosts (ip) seperated by comma for installing salt-minion"
  exit
fi

minion_port_default={{minion_port}}
if [ $4 ]
then
minion_port_default=$4
fi

targetFile=/etc/salt/roster
rm -rf ${targetFile}

echo "$1:">> ${targetFile}
echo "  port: ${minion_port_default} ">> ${targetFile}
echo "  host: $1" >> ${targetFile}
echo "  user: $3" >> ${targetFile}
echo "  passwd: $2" >> ${targetFile}
echo "  sudo: True" >> ${targetFile}
echo "  timeout: 10" >> ${targetFile}

salt-ssh '*' -i state.sls minions.install

exit $?
