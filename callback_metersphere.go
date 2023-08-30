package main

import (
	"CallbackMeterSphere/dto"
	"CallbackMeterSphere/net"
	"CallbackMeterSphere/util"
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	serversMap = make(map[string]dto.Server)
	pool       = util.NewWorkerPool(5)
	initServer dto.InitServerDTO
)

func init() {
	_, err := os.Stat("init.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("init.json文件不存在")
			log.Fatalf("请检查当前环境下是否存在init.json文件！")
		} else {
			fmt.Println("无法访问server_list.json文件:", err)
		}
		os.Exit(0)
	} else {
		initData, _ := os.ReadFile("init.json")
		json.Unmarshal(initData, &initServer)
	}
	var serverList []dto.Server
	//初始化读取服务列表
	_, err = os.Stat("server_list.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("server_list.json文件不存在")
			log.Fatalf("请检查当前环境下是否存在server_list.json文件！")
		} else {
			fmt.Println("无法访问server_list.json文件:", err)
		}
	} else {
		data, err := os.ReadFile("server_list.json")
		if err != nil {
			log.Fatalf("读取失败server_list.json: %v", err)
		}
		err = json.Unmarshal(data, &serverList)
		if err != nil {
			log.Fatalf("无法分析server_list.json: %v", err)
		}
		serverName := ""
		for _, server := range serverList {
			serversMap[server.ServerName] = server
			serverName += server.ServerName + ", "
		}
		log.Printf("数据初始化完成,可用服务数:[%d] 服务名[%s]...", len(serverList), serverName)
	}
}

// 获取测试计划的环境
func getTestPlanEnv(projectID, name string) (string, error) {
	url := fmt.Sprintf("/api/environment/list/%s", projectID)
	cli := net.NewCli(initServer.MeterSphereServer)
	resp, err := cli.Get(url)
	if err != nil {
		return "0", err
	}
	var response struct {
		Success bool              `json:"success"`
		Message interface{}       `json:"message"`
		Data    []dto.Environment `json:"data"`
	}
	_ = json.Unmarshal([]byte(resp), &response)
	for _, env := range response.Data {
		if env.Name == name {
			return env.ID, nil
		}
	}
	return "0", fmt.Errorf("具有名称的环境 '%s' 找不到", name)
}

// 执行测试计划
func runTestPlan(testPlanId, projectId, envId, userId string) (dto.RunPlanResultDTO, error) {
	path := "/track/test/plan/run"
	testPlan := dto.TestPlanDTO{
		Mode:                  "serial",
		ReportType:            "iddReport",
		OnSampleError:         true,
		RunWithinResourcePool: false,
		ResourcePoolId:        nil,
		EnvMap:                map[string]string{projectId: envId},
		TestPlanId:            testPlanId,
		ProjectId:             projectId,
		UserId:                userId,
		TriggerMode:           "MANUAL",
		EnvironmentType:       "JSON",
		EnvironmentGroupId:    "",
		RequestOriginator:     "TEST_PLAN",
	}
	var rtpDTO = dto.RunPlanResultDTO{}
	jsonData, err := json.Marshal(testPlan)
	if err != nil {
		return rtpDTO, err
	}
	cli := net.NewCli(initServer.MeterSphereServer)
	resp, err := cli.Post(path, jsonData)
	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return rtpDTO, err
	}
	json.Unmarshal([]byte(resp), &rtpDTO)
	return rtpDTO, nil
}

// 获取测试报告
func getReportURL(customData string) (string, error) {
	url := "/track/share/generate/expired"
	body := map[string]interface{}{
		"customData": customData,
		"shareType":  "PLAN_DB_REPORT",
		"lang":       nil,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	cli := net.NewCli(initServer.MeterSphereServer)
	resp, err := cli.Post(url, jsonBody)
	if err != nil {
		return "", err
	}
	var response dto.ShareInfoResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return "", err
	}
	reportURL := fmt.Sprintf("%s/track/share-plan-report%s", initServer.MeterSphereServer, response.Data.ShareUrl)
	return reportURL, nil
}

// buildCompleted 流水线构建完成后触发请求
func buildCompleted(ctx *fasthttp.RequestCtx) {
	// 这两行可以获取PostBody数据，文件上传也有用
	postBody := ctx.PostBody()
	flowResult := dto.FlowResult{}
	//log.Print(string(postBody))
	err := json.Unmarshal(postBody, &flowResult)
	if err != nil {
		return
	}
	sName := util.GetServerName(flowResult.Sources[0].Data.Repo)
	targetServer, ok := serversMap[sName]
	if !ok {
		log.Fatalf("server_list中不存在服务名%s ", sName)
	}
	//通知metersphere服务干活了
	noticeMetersphere(targetServer, flowResult)
	//获取返回的服务名称
	fmt.Fprint(ctx, "{\"code\":1,\"msg\":\"处理完成\"}")

}

// 调用meter-sphere服务
func noticeMetersphere(server dto.Server, flow dto.FlowResult) {
	//获取环境的id
	envId, _ := getTestPlanEnv(server.ProjectId, "测试环境")
	//执行测试计划 需要延迟的时间根据server_list中的配置文件来决定
	pool.Schedule(func() {
		log.Printf("正在执行%s服务的测试计划...", server.ServerName)
		data, err := runTestPlan(server.TestPlanId, server.ProjectId, envId, "zsj")
		if err != nil {
			log.Fatalf(err.Error())
		}
		//获取测试报告
		reportUrl, err := getReportURL(data.Data)
		if err != nil {
			log.Fatalf(err.Error())
		}
		var commits []dto.CommitDTO
		log.Printf("测试报告： %s", reportUrl)
		json.Unmarshal([]byte(flow.Sources[0].Data.CommitMsg), &commits)
		decodedString, err := url.QueryUnescape(commits[0].CommitMsg)
		if err != nil {
			fmt.Println("URL解码失败:", err)
			decodedString = commits[0].CommitMsg
		}

		msg := fmt.Sprintf("\n构建流水号: %s\n构建环境: %s\n构建备注: %s\n构建者: %s\n\n代码提交者: %s\n代码提交信息: %s\n代码提交分支: %s\n代码提交ID: %s\n代码提交时间: %s\n\n测试报告地址: %s\n请 @%s 关注测试报告",
			flow.Task.BuildNumber,
			flow.Task.PipelineEnvironment,
			flow.Task.PipelineMark,
			flow.Task.ExecutorName,
			commits[0].CommitAuthor,
			decodedString,
			flow.Sources[0].Data.Branch,
			commits[0].CommitId,
			util.ConvertTime(commits[0].CommitTime),
			reportUrl,
			commits[0].CommitAuthor)
		util.SendMessageToDingTalk(initServer.DingDingAccessToken, msg)

	}, time.Duration(server.DelayedCall)*time.Second)
}

func main() {
	port := initServer.ServerPort
	// 创建路由
	router := fasthttprouter.New()
	// post方法
	router.POST("/api/yunxiao/callback/meterspher", buildCompleted)
	log.Print("服务启动完成,监听端口:", port)
	fasthttp.ListenAndServe(":"+port, router.Handler)
}
