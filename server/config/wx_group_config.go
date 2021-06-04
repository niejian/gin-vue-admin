// config doc

package config

//WxGroupCreateStruct doc
//@Description: 创建群聊结构体
//@Author niejian
type WxGroupCreateStruct struct {
	CorpId  string   `json:"corpId"`
	AgentId string   `json:"agentId"`
	Name    string   `json:"name"`
	Users   []string `json:"users"`
}

type WxGroupChatMsgStruct struct {
	CorpId  string                   `json:"corpId"`
	AgentId string                   `json:"agentId"`
	Data    WxGroupChatMsgDataStruct `json:"data"`
}

type WxGroupChatMsgDataStruct struct {
	ChatId  string `json:"chatid"`
	MsgType string `json:"msgtype" `
	Safe    int    `json:"safe"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type WxGroupUpdateStruct struct {
	CorpId      string   `json:"corpId"`
	AgentId     string   `json:"agentId"`
	Name        string   `json:"name"`
	ChatId      string   `json:"chatId"`
	AddUserList []string `json:"addUserList"`
	DelUserList []string `json:"delUserList"`
}

//CreateChatGroupResponse doc
//@Description: 创建群聊返回结构体
//@Author niejian
type CreateChatGroupResponse struct {
	ResponseCode int                         `json:"responseCode"`
	ResponseMsg  string                      `json:"responseMsg"`
	Success      bool                        `json:"success"`
	Data         CreateChatGroupResponseData `json:"data"`
}

type CreateChatGroupResponseData struct {
	ErrCode int    `json:"errCode"`
	ChatId  string `json:"chatid"`
	ErrMsg  string `json:"errmsg"`
}

type UpdateChatGroupResponse struct {
	ResponseCode int         `json:"responseCode"`
	ResponseMsg  string      `json:"responseMsg"`
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
}

type QueryChatGroupResponse struct {
	ResponseCode int                  `json:"responseCode"`
	ResponseMsg  string               `json:"responseMsg"`
	Success      bool                 `json:"success"`
	Data         QueryChatGroupStruct `json:"data"`
}

type QueryChatGroupStruct struct {
	ErrCode  int      `json:"errCode"`
	ErrMsg   string   `json:"errmsg"`
	ChatInfo ChatInfo `json:"chatinfo"`
}

type ChatInfo struct {
	ChatId   string   `json:"chatid"`
	Name     string   `json:"name"`
	Owner    string   `json:"owner"`
	UserList []string `json:"userList"`
}

//QueryUserInfoResponse doc
//@Description: 获取成员信息结果
//@Author niejian
type QueryUserInfoResponse struct {
	ResponseCode int         `json:"responseCode"`
	ResponseMsg  string      `json:"responseMsg"`
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
}
