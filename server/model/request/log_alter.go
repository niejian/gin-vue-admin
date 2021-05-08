// k8s doc

package request

type ErrorLogAlterConfig struct {
	// 发送工号
	ToUserIds string `json:"toUserIds"`
	// 忽略异常
	Ignores string `json:"ignores"`
	// errs 告警异常
	Errs string `json:"errs"`
	// 应用名称
	AppName string `json:"appName"`
	// 命名空间
	Namespace string `json:"ns"`
	Id        int    `json:"ID"`
}
