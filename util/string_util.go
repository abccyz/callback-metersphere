package util

import (
	"strings"
	"time"
)

func GetServerName(url string) string {
	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]             // 获取最后一个部分
	name = strings.TrimSuffix(name, ".git") // 移除.git后缀
	return name
}

func ConvertTime(timeString string) string {
	timeLayout := "Mon Jan 2 15:04:05 2006 -0700"
	parsedTime, err := time.Parse(timeLayout, timeString)
	if err != nil {
		return ""
	}
	standardTime := parsedTime.Format("2006-01-02 15:04:05")
	return standardTime
}
