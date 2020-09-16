package service

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

// 初始化看门狗环境
// 1. scp 相关的脚本至目标机器
// 2. 登录至目标机器，执行初始化脚本
func InitWatchdogEnv(envStruct *request.InitWatchDogEnvStruct) (string, error) {
	username := envStruct.Username
	password := envStruct.Password
	port := envStruct.Port
	ip := envStruct.Ip
	remoteFilePath := envStruct.RemoteFilePath
	resMsg := ""
	session, err := utils.DoSshConnect(username, password, ip, int(port))
	if nil != err {
		resMsg = "请检查账号、密码、端口是否正确"
		global.GVA_LOG.Error("登录失败：", zap.Any("err", err))
		return resMsg, err
	}

	defer session.Close()
	pwd, _ := os.Getwd()
	if nil != err && "" != resMsg {
		shellCmd := fmt.Sprintf("%s%s%s%s%s%s%s%s%d%s%s", pwd, "/resources/shell/scp.sh ", username, " ", ip, " ", password, " ", port, " ", remoteFilePath)
		// 本地运行scp命令
		command := exec.Command("/bin/bash", "-c", shellCmd)
		rep, err := command.CombinedOutput()

		if nil != err {
			resMsg = err.Error()
			global.GVA_LOG.Error("运行scp:" + shellCmd + ", 返回 " + string(rep))
			return "", err
		}
	}
	return "", err

}
