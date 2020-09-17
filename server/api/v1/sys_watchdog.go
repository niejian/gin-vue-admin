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
