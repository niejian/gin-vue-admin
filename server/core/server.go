package core

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/initialize"
	"github.com/gin-contrib/static"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint {
		// 初始化redis服务
		initialize.Redis()
	}

	if global.GVA_CONFIG.System.UseEs {
		// 初始化es服务
		initialize.ElasticSearch()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")
	Router.Use(static.Serve("/fe", static.LocalFile("./fe", false)))
	Router.Static("/static", "./fe/static")
	//Router.POST("/watchdog/downloadConfig", v1.DownloadConfig)

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Debug("server run success on ", zap.String("address", address))

	fmt.Printf(`欢迎使用 Gin-Vue-Admin
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
