package request

type SshRequestStruct struct {
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

// 初始化看门狗结构体
type InitWatchDogEnvStruct struct {
	Ip             string `json:"ip"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Port           int    `json:"port"`
	RemoteFilePath string `json:"remoteFilePath"`
}

// 初始化看门狗结构体
type ConfigWatchDogEnvStruct struct {
	CIp             string `json:"cip"`
	CUsername       string `json:"cusername"`
	CPassword       string `json:"cpassword"`
	CPort           int    `json:"cport"`
	CRemoteFilePath string `json:"cremoteFilePath"`
}
