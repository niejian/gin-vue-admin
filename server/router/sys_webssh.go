package router

import (
	v1 "gin-vue-admin/api/v1"
	"github.com/gin-gonic/gin"
)

// web终端
func InitSysWebssh(Router *gin.RouterGroup) (R gin.IRoutes) {
	WebsshRouter := Router.Group("")
	//Group("").
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler())
	{
		WebsshRouter.GET("/ws/:id", v1.WsSsh)
	}

	return WebsshRouter
}
