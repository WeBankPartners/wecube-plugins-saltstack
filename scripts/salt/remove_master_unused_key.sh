#!/bin/bash

if [[ $# < 1 ]]; then
  echo "please input the host ip for new added host"
  exit
fi

hosts=$1
hostsArray=(${hosts//,/ })
for host in ${hostsArray[@]} 
do
   salt-key -d $host -y
done

exit 0 

mastersArray=("127.0.0.1")

targetFile=/etc/salt/master_roster
rm -rf ${targetFile}

for host in ${mastersArray[@]}
do
  echo "minion_$host:">> ${targetFile}
  echo "  host: $host" >> ${targetFile}
  echo "  port: 36000" >> ${targetFile}
  echo "  user: root" >> ${targetFile}
  echo "  passwd: Ab888888" >> ${targetFile}
  echo "  sudo: True" >> ${targetFile}
  echo "  timeout: 10" >> ${targetFile}
done


hosts=$1
hostsArray=(${hosts//,/ })

for host in ${hostsArray[@]} 
do
   salt-ssh '*'  --roster-file $targetFile  -i -r "docker exec wecube-plugins-saltstack salt-key -d $host -y"
done




