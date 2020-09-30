package utils

import "time"

var (
	FORMATTIMESTR = "2006-01-02 15:04:05"
)

// 时间戳转换字符串
func FormatTimeByTimestamp(timestampLong int64) string {
	//time.Unix()

	return time.Unix(timestampLong/1000, 0).Format(FORMATTIMESTR)
}

func FormatByNow() string {
	return time.Now().Format(FORMATTIMESTR)
}
