package es

// es存储的数据实体信息
type Exception struct {
	Id           string `json:"id"`
	Ip           string `json:"ip"`
	AppName      string `json:"appName"`
	CreateTime   int64  `json:"createTime"`
	ExceptionTag string `json:"exceptionTag"`
	From         string `json:"from"`
	Msg          string `json:"msg"`
	Year         string `json:"year"`
	Month        string `json:"month"`
	Day          string `json:"day"`
	CreateDate   string `json:"createDate"` //20200916
}
