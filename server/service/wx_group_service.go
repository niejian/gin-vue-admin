// service doc

package service

import (
	"encoding/json"
	"errors"
	"gin-vue-admin/config"
	"gin-vue-admin/global"
	"gin-vue-admin/utils"
	"go.uber.org/zap"
)

//CreateWxGroup doc
//@Description: 创建群聊
//@Author niejian
//@Date 2021-06-03 14:38:12
//@param userIds
//@return string
//@return error
func CreateWxGroup(userIds []string, groupName string) (string, error) {
	url := global.GVA_CONFIG.Wx.BaseUrl + global.GVA_CONFIG.Wx.Group.Create
	data := &config.WxGroupCreateStruct{
		CorpId:  global.GVA_CONFIG.Wx.CorpId,
		AgentId: global.GVA_CONFIG.Wx.AgentId,
		Name:    groupName,
		Users:   userIds,
	}
	var respData *config.CreateChatGroupResponse
	// 获取创建群聊后生成的chatid
	resp := utils.Post(url, data, "")
	err := json.Unmarshal([]byte(resp), &respData)
	if err != nil {
		global.GVA_LOG.Error("创建群聊结果转换失败", zap.String("err", err.Error()))
		return "", err
	}

	return respData.Data.ChatId, nil
}

//UpdateWxGroup doc
//@Description: 更新群聊信息（只能修改群聊名称和人员信息），人员变化在前一步就会判断
//@Author niejian
//@Date 2021-06-03 15:31:19
//@param chatId
//@param name
//@param delUserIdArr
//@param addUserIdArr
//@return string
//@return error
func UpdateWxGroup(chatId, name string, delUserIdArr, addUserIdArr []string) (string, error) {
	url := global.GVA_CONFIG.Wx.BaseUrl + global.GVA_CONFIG.Wx.Group.Update
	var data *config.WxGroupUpdateStruct

	/*
		if "" != name && len(delUserIdArr) > 0 && len(addUserIdArr) > 0 {
			data = &config.WxGroupUpdateStruct{
				CorpId:      global.GVA_CONFIG.Wx.CorpId,
				AgentId:     global.GVA_CONFIG.Wx.AgentId,
				Name:        name,
				ChatId:      chatId,
				AddUserList: addUserIdArr,
				DelUserList: delUserIdArr,
			}
		}

		if "" != name && len(delUserIdArr) == 0 && len(addUserIdArr) == 0 {
			data = &config.WxGroupUpdateStruct{
				CorpId:  global.GVA_CONFIG.Wx.CorpId,
				AgentId: global.GVA_CONFIG.Wx.AgentId,
				Name:    name,
				ChatId:  chatId,
			}
		}

		if "" == name && len(delUserIdArr) > 0 && len(addUserIdArr) > 0 {
			data = &config.WxGroupUpdateStruct{
				CorpId:      global.GVA_CONFIG.Wx.CorpId,
				AgentId:     global.GVA_CONFIG.Wx.AgentId,
				ChatId:      chatId,
				AddUserList: addUserIdArr,
				DelUserList: delUserIdArr,
			}
		}
	*/
	data = &config.WxGroupUpdateStruct{
		CorpId:      global.GVA_CONFIG.Wx.CorpId,
		AgentId:     global.GVA_CONFIG.Wx.AgentId,
		Name:        name,
		ChatId:      chatId,
		AddUserList: addUserIdArr,
		DelUserList: delUserIdArr,
	}
	respStr := utils.Post(url, data, "")
	var respData *config.UpdateChatGroupResponse
	err := json.Unmarshal([]byte(respStr), &respData)
	if err != nil {
		global.GVA_LOG.Error("更新群聊结果转换失败", zap.String("err", err.Error()))
		return "", err
	}

	if !respData.Success {
		return "", errors.New("更新群聊信息失败，请重试")
	}

	return "", nil
}

//QueryChatInfo doc
//@Description: 获取群聊详细信息
//@Author niejian
//@Date 2021-06-04 11:58:03
//@param chatId
//@return *config.QueryChatGroupResponse
//@return error
func QueryChatInfo(chatId string) (*config.QueryChatGroupResponse, error) {
	url := global.GVA_CONFIG.Wx.BaseUrl + global.GVA_CONFIG.Wx.Group.Query +
		"?corpId=" + global.GVA_CONFIG.Wx.CorpId + "&agentId=" + global.GVA_CONFIG.Wx.AgentId +
		"&chatId=" + chatId
	resp := utils.Get(url, "")
	var respData *config.QueryChatGroupResponse
	// 获取创建群聊后生成的chatid
	err := json.Unmarshal([]byte(resp), &respData)
	if err != nil {
		global.GVA_LOG.Error("查询群聊结果转换失败", zap.String("err", err.Error()))
		return nil, err
	}
	return respData, nil
}

//GetWxUser doc
//@Description: 获取企业微信成员信息
//@Author niejian
//@Date 2021-06-04 13:55:56
//@param userId
//@return *config.QueryUserInfoResponse
//@return error
func GetWxUser(userIds []string) ([]*config.QueryUserInfoResponse, error) {
	var datas []*config.QueryUserInfoResponse
	for _, userId := range userIds {
		url := global.GVA_CONFIG.Wx.BaseUrl + global.GVA_CONFIG.Wx.Group.User +
			"?corpId=" + global.GVA_CONFIG.Wx.CorpId + "&userId=" + userId
		resp := utils.Get(url, "")
		var respData *config.QueryUserInfoResponse
		// 获取创建群聊后生成的chatid
		err := json.Unmarshal([]byte(resp), &respData)
		if err != nil {
			global.GVA_LOG.Error("获取企业成员信息失败", zap.String("err", err.Error()))
			err = errors.New("员工号：" + userId + " 不存在，请检查")
			return nil, err
		}
		if !respData.Success {
			err = errors.New("员工号：" + userId + " 不存在，请检查")
			return nil, err
		}
		datas = append(datas, respData)
	}

	return datas, nil
}
