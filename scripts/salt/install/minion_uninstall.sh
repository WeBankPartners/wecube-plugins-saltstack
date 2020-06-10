#!/bin/bash

psout=`ps aux|grep salt-minion|grep -v 'grep'`
if [ ! -n "$psout" ]
then
  systemctl stop salt-minion
fi
rpm -evh salt-minion
rpm -evh salt
rm -rf /etc/salt/*
echo "remove salt-minion done"