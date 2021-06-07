// service doc

package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

var indexPrefix = "watchdog_store_"

//StartCron doc
//@Description: 开启定时任务
//@Author niejian
//@Date 2021-06-07 09:33:44
//@param indexName
//@param cronExpression
func StartCron(data *model.ErrWarnConf) *cron.Cron {
	c := cron.New()
	c.AddFunc(data.SendTime, func() {
		SendErrWarn(data)
	})
	global.Cron_Map[data.IndexName] = c
	c.Start()
	global.GVA_LOG.Info("创建定时任务", zap.Any("data", &data))
	return c
}

//StopCron doc
//@Description: 停止定时任务
//@Author niejian
//@Date 2021-06-07 11:01:46
//@param data
func StopCron(data *model.ErrWarnConf) {
	c := global.Cron_Map[data.IndexName]
	if nil != c {
		global.GVA_LOG.Info("停止运行定时任务", zap.Any("data", &data))
		c.Stop()
		delete(global.Cron_Map, data.IndexName)
		return
	}
	global.GVA_LOG.Info("未发现定时任务", zap.Any("data", &data))
}
