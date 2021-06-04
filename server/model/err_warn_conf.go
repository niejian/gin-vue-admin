// model doc

package model

import "gorm.io/gorm"

type ErrWarnConf struct {
	gorm.Model
	IndexName string `json:"indexName" gorm:"comment:'索引名称, 应用名称'"`
	ToUserIds string `json:"toUserIds" gorm:"comment:'群聊成员'"`
	ChatId    string `json:"ChatId" gorm:"comment:'群聊Id'"`
	SendTime  string `json:"sendTime" gorm:"comment:'定时发送时间,cron表达式'"`
	GroupName string `json:"groupName" gorm:"comment:'群聊名称'"`
}
