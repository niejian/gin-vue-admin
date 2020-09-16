package v1

import (
	"gin-vue-admin/global/response"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
)

// scp 初始化数据
// 初始化看门狗
func WatchDogEnvInit(c *gin.Context) {
	var requestData request.InitWatchDogEnvStruct

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

}
