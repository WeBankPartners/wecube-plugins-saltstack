#!/bin/bash

if [[ $# < 2 ]]; then
  echo "please input the host and password to installing salt-minion"
  exit
fi

targetFile=/etc/salt/roster
rm -rf ${targetFile}

echo "minion_$1:">> ${targetFile}
echo "  port: {{minion_port}} ">> ${targetFile}
echo "  host: $1" >> ${targetFile}
echo "  user: root" >> ${targetFile}
echo "  passwd: $2" >> ${targetFile}
echo "  sudo: True" >> ${targetFile}
echo "  timeout: 10" >> ${targetFile}

salt-ssh '*' -i state.sls minions.install

exit 0

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

salt-ssh '*' -i state.sls minions.install
