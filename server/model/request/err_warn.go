// request doc

package request

type ErrorWarnStruct struct {
	IndexName string   `json:"indexName"`
	ToUserIds []string `json:"toUserIds"`
	SendTime  string   `json:"sendTime"`
	ChatId    string   `json:"chatId"`
	ID        int      `json:"ID"`
	GroupName string   `json:"groupName"`
	IsEnable  int      `json:"isEnable"`
}
