package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

// 看门狗环境初始化、配置更新路由
func InitWatchdogRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	WatchdogRouter := Router.Group("watchdog").
		Use(middleware.JWTAuth()).
		Use(middleware.CasbinHandler())
	{
		WatchdogRouter.POST("/init", v1.WatchDogEnvInit)

	}

	return WatchdogRouter
}
