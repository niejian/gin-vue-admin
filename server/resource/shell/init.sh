#!/bin/bash
if [ ! -n "$1" ]; then
  echo "请输入路径, 如 /data/watchDog"
  exit

fi
cd $1
echo "开始停止"
./stop.sh watch-dog
#echo "开始删除"
#rm -rf watch-dog && rm -rf stop.sh && rm -rf watch.sh && rm -rf restart.sh
#echo "开始下载安装包"
## 下载安装包
#wget --no-check-certificate https://h5oss.bluemoon.com.cn/monitor/watch.sh
#wget --no-check-certificate https://h5oss.bluemoon.com.cn/monitor/watch-dog
#wget --no-check-certificate https://h5oss.bluemoon.com.cn/monitor/restart.sh
#wget --no-check-certificate https://h5oss.bluemoon.com.cn/monitor/stop.sh
mv watch-dog.tmp watch-dog
chmod a+x watch-dog
chmod a+x init.sh
chmod a+x stop.sh
chmod a+x watch.sh
chmod a+x restart.sh
 # 启动守护进程
$1/watch.sh watch-dog > $1/watch-protect.log &
exit