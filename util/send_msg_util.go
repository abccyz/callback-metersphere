package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMessageToDingTalk(accessToken, message string) {
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", accessToken)
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("发送消息失败:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("消息发送成功")
	} else {
		fmt.Println("发送消息失败")
	}
}
