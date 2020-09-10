package initialize

import (
	"gin-vue-admin/global"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"os"
	"time"
)

// 初始化elasticSearch客户端
func ElasticSearch() {
	esConfig := global.GVA_CONFIG.Es

	urls := esConfig.Urls
	if nil == urls || len(urls) == 0 {
		global.GVA_LOG.Info("es 地址为空")
	}

	username := esConfig.Username
	password := esConfig.Password

	if "" != username && "" != password {
		// 未开启密码
		global.GVA_ES, _ = elastic.NewClient(
			elastic.SetURL(urls...),
			elastic.SetBasicAuth(username, password),
			elastic.SetSniff(true),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetGzip(true),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC-ERR-", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC-INFO-", log.LstdFlags)),
			elastic.SetHeaders(http.Header{
				"X-Caller-Id": []string{"..."},
			}),
		)
	} else {
		// 开启密码
		elastic.NewClient(
			elastic.SetURL(urls...),
			elastic.SetSniff(true),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetGzip(true),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC-ERR-", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC-INFO-", log.LstdFlags)),
			elastic.SetHeaders(http.Header{
				"X-Caller-Id": []string{"..."},
			}),
		)
	}

}
