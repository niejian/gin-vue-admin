// router doc

package router

import (
	v1 "gin-vue-admin/api/v1"
	"github.com/gin-gonic/gin"
)

//InitK8sApi doc
//@Description: 获取ns，deploy等信息
//@Author niejian
//@Date 2021-05-06 11:41:43
//@param Router
func InitK8sApi(Router *gin.RouterGroup) {
	//ApiRouter := Router.Group("api").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	ApiRouter := Router.Group("k8sapi")
	{
		ApiRouter.GET("nsList", v1.ListNs)
		// 获取某个命名空间的deploy信息
		ApiRouter.GET("deploys/:ns", v1.ListDeploy)
		// 获取pod名称
		ApiRouter.GET("pods/:ns/:labels", v1.ListPodsByLabels)
		// 添加配置信息
		ApiRouter.POST("watchConf/addOrUpdate", v1.AddWatchdogConf)
		ApiRouter.GET("getConfByNsAndAppName/:ns/:appName", v1.GetConfByNsAndAppName)
		// 分页获取配置列表信息
		ApiRouter.POST("getConfList", v1.GetConfigList)
	}
}
