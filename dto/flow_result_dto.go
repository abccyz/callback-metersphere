package dto

type FlowResult struct {
	Event        string    `json:"event"`
	Action       string    `json:"action"`
	Task         Task      `json:"task"`
	Pipeline     *Pipeline `json:"pipeline"`
	Artifacts    []string  `json:"artifacts"`
	Sources      []Source  `json:"sources"`
	GlobalParams []string  `json:"globalParams"`
}

type Task struct {
	PipelineId          string `json:"pipelineId"`
	PipelineName        string `json:"pipelineName"`
	StageName           string `json:"stageName"`
	TaskName            string `json:"taskName"`
	BuildNumber         string `json:"buildNumber"`
	StatusCode          string `json:"statusCode"`
	StatusName          string `json:"statusName"`
	PipelineUrl         string `json:"pipelineUrl"`
	Message             string `json:"message"`
	ExecutorName        string `json:"executorName"`
	PipelineTags        string `json:"pipelineTags"`
	PipelineEnvironment string `json:"pipelineEnvironment"`
	FlowInstId          string `json:"flowInstId"`
	PipelineInstId      string `json:"pipelineInstId"`
	PipelineMark        string `json:"pipelineMark"`
}

type Pipeline struct {
}

type Source struct {
	Name string `json:"name"`
	Sign string `json:"sign"`
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	Repo             string   `json:"repo"`
	Branch           string   `json:"branch"`
	CommitID         string   `json:"commitId"`
	PreviousCommitID *string  `json:"privousCommitId"`
	CommitMsg        string   `json:"commitMsg"`
	Args             []string `json:"args"`
}
