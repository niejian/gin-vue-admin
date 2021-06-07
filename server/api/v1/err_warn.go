// v1 doc

package v1

import (
	"fmt"
	"gin-vue-admin/global/response"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

//AddOrUpdateErrorWarn doc
//@Description: 创建或更新
//@Author niejian
//@Date 2021-06-03 13:51:21
//@param c
func AddOrUpdateErrorWarn(c *gin.Context) {
	var requestData request.ErrorWarnStruct
	// 绑定参数
	c.ShouldBindJSON(&requestData)
	// 校验请求参数是否为空
	rules := utils.Rules{
		"indexName": {utils.NotEmpty()},
		"toUserIds": {utils.NotEmpty()},
		"sendTime":  {utils.NotEmpty()},
		"groupName": {utils.NotEmpty()},
	}

	verifyErr := utils.Verify(requestData, rules)
	if nil != verifyErr {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	errorWarnConf := &model.ErrWarnConf{
		IndexName: requestData.IndexName,
		ToUserIds: strings.Join(requestData.ToUserIds, "|"),
		ChatId:    requestData.ChatId,
		SendTime:  requestData.SendTime,
		GroupName: requestData.GroupName,
		IsEnable:  requestData.IsEnable,
	}
	// 判断员工号是否合法
	userIds := requestData.ToUserIds
	//userIds := strings.Split(toUserIds, "|")
	_, err := service.GetWxUser(userIds)
	if nil != err {
		response.FailWithDetailed(response.ERROR, err.Error(), fmt.Sprintf("%v", err), c)
		return
	}
	// 判断新增还是更新
	id := requestData.ID
	if id <= 0 {
		// 新增
		errorWarnConf.ChatId = ""
		_, err := service.AddErrWarnConf(errorWarnConf)
		if err != nil {
			response.FailWithDetailed(response.ERROR, "添加失败", fmt.Sprintf("%v", err), c)
		} else {
			response.OkDetailed("", "添加成功", c)
		}
	} else {
		errorWarnConf.ID = uint(requestData.ID)
		err := service.UpdateErrWarnConf(errorWarnConf)
		if err == nil {
			response.OkDetailed("", "更新成功", c)
		} else {
			response.FailWithDetailed(response.ERROR, err.Error(), fmt.Sprintf("%v", err), c)
		}
	}
}

func GetConfInfoByIndexName(c *gin.Context) {
	indexName := c.Param("indexName")
	conf, err := service.GetWarnConfByIndexName(indexName)
	if err != nil {
		response.FailWithDetailed(response.ERROR, "查询失败", fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(conf, "查询成功", c)
	}
}

func GetUserInfo(c *gin.Context) {
	userIds := c.Param("userIds")
	userIdArr := strings.Split(userIds, ",")
	query, err := service.GetWxUser(userIdArr)

	if err != nil {
		response.FailWithDetailed(response.ERROR, "", fmt.Sprintf("%v", err), c)
	} else {

		response.OkDetailed(&query, "查询成功", c)

	}
}
