#!/usr/bin/env bash

##停止应用
#killall pkmm
#sleep 1
#echo "关闭了进程 pkmm"
##更新
#git pull
#echo "更新代码"
#
##编译
#echo "现在开始编译啦~~"
#source /etc/profile
#go build
#echo "编译完成~~"
#
##启动
##nohup ./pkmm 2>&1 >> info.log 2>&1 /dev/null &
#nohup ./pkmm 1> out_server.out 2> err_server.err &
#echo "已经启动了！"

# ################

# 使用supervisor守护

echo "首先停止pkmm进程"
supervisorctl stop pkmm

echo "同步代码"
git pull
sleep 1

echo "编译最新的代码"
source /etc/profile
go build
echo "编译完成~~"

echo "启动pkmm"
supervisorctl start pkmm

echo "服务已经开启~！"