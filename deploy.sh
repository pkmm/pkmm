#!/usr/bin/env bash

#停止应用
killall pkmm
#更新
git pull

#编译
source /etc/profile
go build

#启动
nohup ./pkmm 2>&1 >> info.log 2>&1 /dev/null &
