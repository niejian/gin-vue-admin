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

	shellCmd := fmt.Sprintf("%s%s%s%s%s%s%s%s%d%s%s", pwd, "/resource/shell/scp.sh ", username,
		" ", ip, " ", password, " ", port, " ", remoteFilePath)
	// 本地运行scp命令
	command := exec.Command("/bin/bash", "-c", shellCmd)
	rep, err := command.CombinedOutput()

	if nil != err {
		resMsg = err.Error()
		global.GVA_LOG.Error("运行scp:" + shellCmd + ", 返回 " + string(rep))
		return err.Error(), err
	}

	global.GVA_LOG.Info("scp 结果：", zap.Any("scp", string(rep)))

	global.GVA_LOG.Info("init 当前路径：", zap.String("init-path", pwd))
	err = utils.DoPathscp(username, password, ip, pwd+"/resource/shell/init.sh", remoteFilePath, port)
	if err != nil {
		resMsg = "请检查账号、密码、端口是否正确"
		global.GVA_LOG.Error("请检查账号、密码、端口是否正确")
		return err.Error(), err
	}

	session, err = utils.DoSshConnect(username, password, ip, int(port))
	if err != nil {
		resMsg = "请检查账号、密码、端口是否正确"
		global.GVA_LOG.Error("请检查账号、密码、端口是否正确")
		return err.Error(), err
	}

	defer session.Close()

	cmd := remoteFilePath + "/init.sh " + remoteFilePath + " >" + remoteFilePath + "/init.log"
	global.GVA_LOG.Info("开始运行 init.sh 脚本", zap.Any("cmd", cmd))
	out, err := session.CombinedOutput(cmd)
	if nil != err {
		global.GVA_LOG.Error("运行 init.sh 脚本失败", zap.Any("err", err))
		resMsg = err.Error()
	} else {
		global.GVA_LOG.Info("运行 init.sh 脚本成功，结果：", zap.Any("out", out))
		resMsg = ""
	}
	return resMsg, err

}

// 复制配置文件到本机器
func CopyConfig(requestData request.ConfigWatchDogEnvStruct, filename string) error {
	pwd, _ := os.Getwd()
	shellCmd := fmt.Sprintf("%s%s%s%s%s%s%s%s%d%s%s%s", pwd, "/resource/shell/copyconfig.sh ", requestData.CUsername,
		" ", requestData.CIp, " ", requestData.CPassword, " ", requestData.CPort, " ", requestData.CRemoteFilePath, " "+filename)

	// 本地运行scp命令
	command := exec.Command("/bin/bash", "-c", shellCmd)
	global.GVA_LOG.Info("复制配置命令：", zap.Any("cmd", shellCmd))
	rep, err := command.CombinedOutput()

	if nil != err {
		global.GVA_LOG.Error("运行copy:" + shellCmd + ", 返回 " + string(rep))
		return err
	}

	return nil
}
