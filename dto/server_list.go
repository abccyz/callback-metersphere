package dto

type Server struct {
	ServerName  string `json:"serverName"`
	ProjectId   string `json:"projectId"`
	TestPlanId  string `json:"testPlanId"`
	DelayedCall int    `json:"delayedCall"` //延迟调用时间为秒
}
