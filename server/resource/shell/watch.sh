
#! /bin/sh
#进程名字可修改
#if [ ! -n "$1" ]; then
#  echo "请输入进程名, 如 watchDog"
#  exit
#fi

PRO_NAME=watch-dog
CMD="ps aux | grep ${PRO_NAME} | grep -v grep | grep -v 'watch.sh' |wc -l"
echo $CMD
while true ; do

#    用ps获取$PRO_NAME进程数量
#NUM=$(ps aux|grep ${PRO_NAME}|grep -v grep|wc -l)
NUM=$(ps aux | grep ${PRO_NAME} | grep -v grep | grep -v 'watch.sh' |wc -l)
echo $NUM

#    少于1，重启进程
  if [ "${NUM}" -lt "1" ];then
#  if [ $NUM -lt 1 ];then
    echo "$1 进程异常，正在重启"
    ./restart.sh $PRO_NAME
  fi
  sleep 10
done
