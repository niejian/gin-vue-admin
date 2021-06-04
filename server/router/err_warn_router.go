// router doc

package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

//InitWatchdogErrorWarnRouter doc
//@Description: 初始化异常统计告警
//@Author niejian
//@Date 2021-06-03 13:49:33
//@param Router
//@return R
func InitWatchdogErrorWarnRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("errorWarn").
		Use(middleware.JWTAuth()).
		Use(middleware.CasbinHandler())
	{
		router.POST("/addOrUpdate", v1.AddOrUpdateErrorWarn)
		router.GET("/getConfInfoByIndexName/:indexName", v1.GetConfInfoByIndexName)
		router.GET("/getUserInfo/:userIds", v1.GetUserInfo)

	}

	return router
}
