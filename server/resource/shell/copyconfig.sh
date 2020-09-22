#!/usr/bin/expect

#set passwd "bluemoon2016#"
set username [lindex $argv 0]
set host [lindex $argv 1]
set passwd [lindex $argv 2]
set port [lindex $argv 3]
set remoteDir [lindex $argv 4]
set filename [lindex $argv 5]

#spawn scp -P $port $username@$host:$remoteDir/watchDog.yaml /home/appadm/gva/resource/downloads/$filename
spawn scp -P $port $username@$host:$remoteDir/watchDog.yaml /Users/a/myproject/go/src/gin-vue-admin/server/resource/downloads/$filename
#spawn scp -P $port /data/watchDog/web/static/init/init.sh /data/watchDog/web/static/init/restart.sh /data/watchDog/web/static/init/stop.sh /data/watchDog/web/static/init/watch.sh /data/watchDog/watch-dog $username@$host:$remoteDir
expect {
  "password"
	{
	  send "$passwd\n"
	}
  "密码："
        {
          send "$passwd\n"
        }
   "pass"
        {
          send "$passwd\n"
        }
   "yes/no"
        {
          send_user "send yes"
          send "yes\n"
        }
   eof
    {
        send_user "eof\n"
    }
}

set timeout 10000
send "exit\r"
expect eof