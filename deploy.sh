#!/usr/bin/env bash

#停止应用
killall pkmm
sleep 1
echo "关闭了进程 pkmm"
#更新
git pull
echo "更新代码"

#编译
echo "现在开始编译啦~~"
source /etc/profile
go build
echo "编译完成~~"

#启动
nohup ./pkmm 2>&1 >> info.log 2>&1 /dev/null &
echo "已经启动了！"