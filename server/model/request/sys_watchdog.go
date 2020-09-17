package request

// 初始化看门狗结构体
type InitWatchDogEnvStruct struct {
	Ip             string `json:"ip"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Port           int    `json:"port"`
	RemoteFilePath string `json:"remoteFilePath"`
}
