package core

import (
	"fmt"
	"gin-vue-admin/global"
	_ "gin-vue-admin/packfile"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const defaultConfigFile = "config.yaml"

// package init 函数在main方法之前调用
func init() {
	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	global.GVA_VP = v
	// 判断是否是k8s环境，如果是，那么数据库，es读取k8s部署文件中的配置信息
	enableK8s := os.Getenv("ENABLE_K8S")
	if "" != enableK8s && "true" == enableK8s {
		fmt.Println("读取k8s配置信息 ")
		esUrls := os.Getenv("ES_URLS")
		esUsername := os.Getenv("ES_USERNAME")
		esPassword := os.Getenv("ES_PASSWORD")

		if "" != esUrls && "" != esUsername && "" != esPassword {
			// 替换
			global.GVA_CONFIG.Es.Urls = strings.Split(esUrls, ",")
			global.GVA_CONFIG.Es.Username = esUsername
			global.GVA_CONFIG.Es.Password = esPassword
		}

		mysqlUrl := os.Getenv("MYSQL_URLS")
		mysqlUsername := os.Getenv("MYSQL_USERNAME")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")

		if "" != mysqlUrl && "" != mysqlUsername && "" != mysqlPassword {
			global.GVA_CONFIG.Mysql.Path = mysqlUrl
			global.GVA_CONFIG.Mysql.Username = mysqlUsername
			global.GVA_CONFIG.Mysql.Password = mysqlPassword
		}
	}
}
