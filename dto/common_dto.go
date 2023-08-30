package dto

type Environment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectId string `json:"projectId"`
}

type ShareInfoResponse struct {
	Data ShareInfoData `json:"data"`
}

type ShareInfoData struct {
	ShareUrl string `json:"shareUrl"`
}

// TestPlan --------------执行测试计划结构体--------------
type TestPlanDTO struct {
	Mode                  string            `json:"mode"`
	ReportType            string            `json:"reportType"`
	OnSampleError         bool              `json:"onSampleError"`
	RunWithinResourcePool bool              `json:"runWithinResourcePool"`
	ResourcePoolId        interface{}       `json:"resourcePoolId"`
	EnvMap                map[string]string `json:"envMap"`
	TestPlanId            string            `json:"testPlanId"`
	ProjectId             string            `json:"projectId"`
	UserId                string            `json:"userId"`
	TriggerMode           string            `json:"triggerMode"`
	EnvironmentType       string            `json:"environmentType"`
	EnvironmentGroupId    string            `json:"environmentGroupId"`
	RequestOriginator     string            `json:"requestOriginator"`
}

// RunPlanResultDTO 执行测试计划返回结构体
type RunPlanResultDTO struct {
	Success string `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// CommitDTO 代码提交信息
type CommitDTO struct {
	CommitMsg    string `json:"commitMsg"`
	CommitId     string `json:"commitId"`
	CommitAuthor string `json:"commitAuthor"`
	CommitTime   string `json:"commitTime"`
}

type InitServerDTO struct {
	ServerPort          string `json:"serverPort"`
	DingDingAccessToken string `json:"dingDingAccessToken"`
	MeterSphereServer   string `json:"meterSphereServer"`
}
