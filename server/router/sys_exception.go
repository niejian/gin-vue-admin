package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitExceptionViewRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	ExceptionViewRouter := Router.
		Group("exception").
		Use(middleware.JWTAuth()).
		Use(middleware.CasbinHandler())

	{
		// 添加路由信息
		// 获取过去x天的异常报表信息
		ExceptionViewRouter.POST("exceptionOverview", v1.GetExceptionOverview)
		// 1. 获取当天的详细预览信息
		ExceptionViewRouter.POST("viewException", v1.GetExceptionView)
		// 2. 通过索引名、字段名获取异常详细信息
		ExceptionViewRouter.POST("exceptionDetails", v1.GetExceptionDetails)
		// 3. 通过索引名获取信息
		ExceptionViewRouter.POST("indexException", v1.GetExceptionDailyViewByIndexName)
		ExceptionViewRouter.GET("getExceptionById/:indexName/:id", v1.GetExceptionById)
	}

	return ExceptionViewRouter

}
