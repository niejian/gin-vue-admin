package v1

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/global/response"
	"gin-vue-admin/model/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const (
	DOWNLOAD_PATH = "/Users/a/myproject/go/src/gin-vue-admin/server/resource/downloads/"
	//DOWNLOAD_PATH = "/home/appadm/gva/resource/downloads/"
)

// scp 初始化数据
// 初始化看门狗
func WatchDogEnvInit(c *gin.Context) {
	var requestData request.InitWatchDogEnvStruct
	// 绑定参数
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"Ip":             {utils.NotEmpty()},
		"Username":       {utils.NotEmpty()},
		"Password":       {utils.NotEmpty()},
		"Port":           {utils.NotEmpty()},
		"RemoteFilePath": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	resp, err := service.InitWatchdogEnv(&requestData)
	if err != nil {
		global.GVA_LOG.Info("初始化环境失败，请重试：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("初始化环境失败，%v", err), c)
	} else {
		// 请求成功
		response.OkWithData(resp, c)
	}

}

// 下载配置文件
// scp target local
// http local
func DownloadConfig(c *gin.Context) {
	var requestData request.ConfigWatchDogEnvStruct
	// 绑定参数
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"CIp":             {utils.NotEmpty()},
		"CUsername":       {utils.NotEmpty()},
		"CPassword":       {utils.NotEmpty()},
		"CPort":           {utils.NotEmpty()},
		"CRemoteFilePath": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	session, err := connectTarget(requestData.CUsername, requestData.CPassword,
		requestData.CIp, requestData.CPort)
	defer session.Close()
	if err != nil {
		global.GVA_LOG.Info("连接用户名或密码错误，请重试：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("连接用户名或密码错误，请重试，%v", err), c)
	} else {
		filename := fmt.Sprintf("%s%s%s", requestData.CIp, ":", "watchDog.yaml")
		// 请求成功，执行复制脚本
		err := service.CopyConfig(requestData, filename)
		if nil != err {
			global.GVA_LOG.Info("复制配置文件失败，请重试：", zap.Any("err", err))
			response.FailWithMessage(fmt.Sprintf("复制配置文件失败，请重试，%v", err), c)
		} else {
			filePath := fmt.Sprintf("%s%s", DOWNLOAD_PATH, filename)

			// download file
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", "watchDog.yaml")) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
			//c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filePath)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
			c.Writer.Header().Set("Content-Type", "application/octet-stream")
			c.File(filePath)

		}

	}

}

func connectTarget(username, password, ip string, port int) (*ssh.Session, error) {
	session, err := utils.DoSshConnect(username, password,
		ip, port)

	return session, err
}

// 尝试连接
func Try2Connect(c *gin.Context) {
	var requestData request.InitWatchDogEnvStruct
	// 绑定参数
	c.ShouldBindJSON(&requestData)

	// 校验请求参数是否为空
	rules := utils.Rules{
		"Ip":       {utils.NotEmpty()},
		"Username": {utils.NotEmpty()},
		"Password": {utils.NotEmpty()},
		"Port":     {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	session, err := connectTarget(requestData.Username, requestData.Password, requestData.Ip, requestData.Port)

	defer session.Close()
	if err != nil {
		global.GVA_LOG.Info("连接用户名或密码错误，请重试：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("初始化环境失败，%v", err), c)
	} else {
		// 请求成功
		response.OkWithData("success", c)
	}

}

// 文件上传
func Upload(c *gin.Context) {
	var requestData request.ConfigWatchDogEnvStruct
	// 绑定参数
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"CIp":             {utils.NotEmpty()},
		"CUsername":       {utils.NotEmpty()},
		"CPassword":       {utils.NotEmpty()},
		"CPort":           {utils.NotEmpty()},
		"CRemoteFilePath": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	session, err := connectTarget(requestData.CUsername, requestData.CPassword,
		requestData.CIp, requestData.CPort)
	defer session.Close()
	if err != nil {
		global.GVA_LOG.Info("连接用户名或密码错误，请重试：", zap.Any("err", err))
		response.FailWithMessage(fmt.Sprintf("连接用户名或密码错误，请重试，%v", err), c)
	} else {
		localFile, _ = c.FormFile("")
	}
}
