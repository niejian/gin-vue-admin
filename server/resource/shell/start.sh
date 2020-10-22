#!/bin/bash
if [ ! -n "$1" ]; then
  echo "请输入进程名, 如 watch-dog"
  exit

fi
PIDFILE=watchDog.pid
DirName=$(cd $(dirname $0); pwd)
cd $DirName

if [ -f "$PIDFILE" ]; then
    echo "【看门狗】已启动 ..., 运行stop.sh"
    ./stop.sh
else
  echo "【看门狗】开始启动..."
  nohup $DirName/$1 2>&1 &
  printf '%d' $! > $PIDFILE
  echo "【看门狗】启动成功"
fi
