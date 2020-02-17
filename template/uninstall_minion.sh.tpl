#!/bin/bash

if [[ $# < 1 ]]; then
  echo "please input the hosts (ip) seperated by comma for uninstalling salt-minion"
  exit
fi

targetFile=/etc/salt/roster
rm -rf ${targetFile}

echo "$1:">> ${targetFile}
echo "  port: {{minion_port}} ">> ${targetFile}
echo "  host: $1" >> ${targetFile}
echo "  user: root" >> ${targetFile}
echo "  passwd: $2" >> ${targetFile}
echo "  sudo: True" >> ${targetFile}
echo "  timeout: 10" >> ${targetFile}

salt-ssh '*' -i state.sls minions.uninstall

exit $?

#for salt-msster ha
hosts=$1
hostsArray=(${hosts//,/ })

for host in ${hostsArray[@]} 
do
  echo "minion_$host:">> ${targetFile}
  echo "  port: {{minion_port}} ">> ${targetFile}
  echo "  host: $host" >> ${targetFile}
  echo "  user: root" >> ${targetFile}
  echo "  passwd: {{minion_passwd}}" >> ${targetFile}
  echo "  sudo: True" >> ${targetFile}
  echo "  timeout: 10" >> ${targetFile}
done

salt-ssh '*' -i state.sls minions.uninstall