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



