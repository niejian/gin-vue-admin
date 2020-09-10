package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitExceptionViewRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	ExceptionViewRouter := Router.Group("exception").
		Use(middleware.JWTAuth()).
		Use(middleware.CasbinHandler())

	{
		// 添加路由信息
		// 1. 获取预览信息
		ExceptionViewRouter.POST("viewException", v1.GetExceptionView)
		// 2. 获取异常详细信息
		ExceptionViewRouter.POST("exceptionDetails")
	}

	return ExceptionViewRouter

}
