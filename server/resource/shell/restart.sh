#!/bin/bash
if [ ! -n "$1" ]; then
  echo "请输入进程名, 如 watch-dog"
  exit

fi
echo "开始重启...."
./stop.sh $1
./start.sh $1
echo "启动成功"

#/bin/sh $COMMOND
