package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"time"
)

//任务执行接口
type IJobExecute interface {
	Execute()
}

//任务执行处理函数
type ExecuteHandler func()

type Job struct {
	Name         string         `json:"name"`
	JobExeTarget IJobExecute    `json:"jobExeTarget"`
	CronExpr     string         `json:"cronExpr"`
	ExeHandler   ExecuteHandler `json:"exeHandler"`
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time
}

type JobExecuteInfo struct {
	Job      *Job
	PlanTime time.Time //任务计划执行时间
	RealTime time.Time //任务实际执行时间
}

type JobExecuteResult struct {
	JobExecuteInfo *JobExecuteInfo
	OutPut         []byte
	Err            error
	StartTime      time.Time
	EndTime        time.Time
}

type ApiResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type JobEvent struct {
	EventType int //save,delete
	Job       *Job
}

//构建Api响应结果
func BuildApiResult(code int, msg string, data interface{}) (jsonValue string) {

	apiResult := ApiResult{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	jsonValueByte, _ := json.Marshal(apiResult)
	jsonValue = string(jsonValueByte)
	return jsonValue
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	jobEvent = &JobEvent{
		EventType: eventType,
		Job:       job,
	}

	return jobEvent
}

//构建调度计划
func BuildSchedulerPlan(job *Job) (plan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	plan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return plan, err
}

func BuildJobExecuteInfo(jobSchedulerPlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulerPlan.Job,
		PlanTime: jobSchedulerPlan.NextTime,
		RealTime: time.Now(),
	}

	return jobExecuteInfo
}
