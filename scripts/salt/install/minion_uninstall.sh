#!/bin/bash

psout=`ps aux|grep salt-minion|grep -v 'grep'`
if [ ! -n "$psout" ]
then
  sudo systemctl stop salt-minion
fi
sudo rpm -evh salt-minion
sudo rpm -evh salt
sudo rm -rf /etc/salt/*
echo "remove salt-minion done"