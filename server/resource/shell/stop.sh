#!/bin/bash
if [ ! -n "$1" ]; then
  echo "请输入进程名, 如 watchDog"
  exit

fi

#dirName = "project_path=$(cd `dirname $0`; pwd)"

#PIDS=$(ps -ef | grep $1 | grep -v grep |grep -v 'restart.sh'| awk '{print $2}' | grep -v awk)
PIDS=`ps -ef | grep $1 | grep -v grep |grep -v 'restart.sh'| awk '{print $2}' | grep -v awk`
#PIDS=$(ps -ef | grep $1 | grep -v grep| awk '{print $2}' | grep -v awk)
#PID=$(ps -ef|grep $1|grep -v "grep"|grep -v "$0"|awk '{printf $2}')
#PIDSTR=$(echo $PIDS)
#
##echo "进程ID: $PIDSTR"
#if [ ! -n "$PIDSTR" ]; then
#  echo "$1 进程不存在"
#else
#  echo "准备kill进程: $PIDSTR"
#  kill -9 $PIDSTR
#fi
#echo "进程: $PIDSTR，kill完成"

for id in $PIDS
do
  kill -9 $id
  echo "killed $id"
done