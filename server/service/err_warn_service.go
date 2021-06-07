// service doc

package service

import (
	"bytes"
	"errors"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

//AddErrWarnConf doc
//@Description: 添加
//@Author niejian
//@Date 2021-06-03 14:07:09
//@param errWarnConf
//@return *model.ErrWarnConf
//@return error
func AddErrWarnConf(errWarnConf *model.ErrWarnConf) (*model.ErrWarnConf, error) {
	// 调用企业微信接口，创建群聊信息，获取到chatId
	chatId, err := CreateWxGroup(strings.Split(errWarnConf.ToUserIds, "|"), errWarnConf.GroupName)
	if err != nil {
		return nil, err
	}
	errWarnConf.ChatId = chatId
	err = global.GVA_DB.Table("err_warn_conf").Create(&errWarnConf).Error
	// 开始定时任务
	if 1 == errWarnConf.IsEnable {
		StartCron(errWarnConf)
	}
	return errWarnConf, err
}

//UpdateErrWarnConf doc
//@Description: 更新
//@Author niejian
//@Date 2021-06-03 14:08:20
//@param errWarnConf
func UpdateErrWarnConf(errWarnConf *model.ErrWarnConf) error {
	// 查询该记录的创建时间等，然后赋值
	var data model.ErrWarnConf
	// 按住键查询
	global.GVA_DB.Table("err_warn_conf").Find(&data, errWarnConf.ID)
	// 查到数据
	if data.ID > 0 {
		var err error
		errWarnConf.CreatedAt = data.CreatedAt
		errWarnConf.DeletedAt = data.DeletedAt
		// 调用查询群聊接口，获取到具体的群聊信息
		queryData, errQuery := QueryChatInfo(errWarnConf.ChatId)
		if nil != errQuery {
			return errors.New("查询此群聊信息失败，请重试")
		}
		// 查询失败
		if !queryData.Success {
			return errors.New("查询此群聊信息失败，请重试")
		}

		newUserIds := errWarnConf.ToUserIds
		// 获取到原有的群聊成员信息
		delUserIdArr := queryData.Data.ChatInfo.UserList
		// 前端传进的成员信息
		addUserIdArr := strings.Split(newUserIds, "|")
		//same := StringSliceEqual(userIdArr, newUserIdArr)
		//groupName := data.GroupName
		newGroupName := errWarnConf.GroupName
		//var emptyStrArr []string
		_, err = UpdateWxGroup(errWarnConf.ChatId, newGroupName, delUserIdArr, addUserIdArr)

		//if groupName != newGroupName && same {
		//	_, err = UpdateWxGroup(errWarnConf.ChatId, newGroupName, emptyStrArr, emptyStrArr)
		//}
		//if groupName != newGroupName && !same {
		//	_, err = UpdateWxGroup(errWarnConf.ChatId, newGroupName, userIdArr, newUserIdArr)
		//}
		//
		//if groupName == newGroupName && !same {
		//	_, err = UpdateWxGroup(errWarnConf.ChatId, "", userIdArr, newUserIdArr)
		//}
		if err == nil {
			global.GVA_DB.Table("err_warn_conf").Save(errWarnConf)
			// 关闭定时任务
			StopCron(errWarnConf)
			// 开始定时任务
			if 1 == errWarnConf.IsEnable {
				StartCron(errWarnConf)
			}
			return nil
		}
		return err
	}
	return errors.New("查无此数据")
}

func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

//GetWarnConfByIndexName doc
//@Description: 通过索引名获取告警提醒配置
//@Author niejian
//@Date 2021-06-04 10:35:08
//@param indexName
//@return *model.ErrWarnConf
//@return error
func GetWarnConfByIndexName(indexName string) (*model.ErrWarnConf, error) {
	var confs []model.ErrWarnConf
	err := global.GVA_DB.Table("err_warn_conf").Where("index_name = ?", indexName).Find(&confs).Error
	if nil != err {
		return nil, err
	}
	if len(confs) > 0 {
		return &confs[0], nil
	}
	return nil, nil
}

//SendErrWarn doc
//@Description: 发送统计告警信息
//@Author niejian
//@Date 2021-06-07 10:53:10
//@param data
func SendErrWarn(data *model.ErrWarnConf) {
	// 分类统计前一天改索引的异常信息
	//aggs := GetExceptionOverview(data.IndexName, 1)
	durationHour, _ := time.ParseDuration(fmt.Sprintf("%d%s", 0, "h"))
	warnTime := time.Now().Add(durationHour).Format(FormatStartTime)
	aggs := IndexOverview(data.IndexName, warnTime)
	// 告警信息
	var msg bytes.Buffer

	if len(aggs) > 0 {
		appName := strings.ReplaceAll(data.IndexName, indexPrefix, "")
		msg.WriteString("时间：" + warnTime + "\n")
		msg.WriteString("应用：" + appName + "\n")
		for _, aggIndex := range aggs {
			// 异常关键字
			e := aggIndex.Key
			// 告警次数
			count := aggIndex.DocCount
			msg.WriteString("异常：" + e + ", 数量：" + strconv.FormatInt(count, 10) + "\n")

		}
		// 发送群聊信息
		_, err := SendWxGroupMsg(data.ChatId, msg.String())
		if err != nil {
			SendWxGroupMsg(data.ChatId, msg.String())
		}
	} else {
		global.GVA_LOG.Error("无告警信息", zap.String("indexName", data.IndexName))
	}
}
