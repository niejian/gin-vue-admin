// initialize doc

package initialize

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/service"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

//InitializeCron doc
//@Description: 初始化定时任务
//@Author niejian
//@Date 2021-06-07 10:41:59
func InitializeCron() {
	var confs []model.ErrWarnConf
	// 查询所有定时任务信息
	err := global.GVA_DB.Table("err_warn_conf").Where("is_enable = ?", 1).Find(&confs).Error
	if err != nil {
		global.GVA_LOG.Error("查询定时任务失败，请重试", zap.String("err", err.Error()))
		panic(err.Error())
	}
	// 开启定时任务
	for _, data := range confs {
		c := cron.New()
		expression := data.SendTime
		global.GVA_LOG.Info("定时任务", zap.String("indexName", data.IndexName), zap.String("表达式", expression))
		err := c.AddFunc(data.SendTime, func() {
			global.GVA_LOG.Info("定时任务开始启动。。。。")
			service.SendErrWarn(&data)
		})
		if nil != err {
			global.GVA_LOG.Error("err", zap.String("err", err.Error()))
		}
		global.Cron_Map[data.IndexName] = c
		c.Start()
	}

}
