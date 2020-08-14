package worker

import (
	"fmt"
	"go-quartz/common"
	"time"
)

type Scheduler struct {
	jobEventChan      chan *common.JobEvent              //etcd任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo  //任务执行表
	jobExeResultChan  chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

func initSchduler() {

	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExeResultChan:  make(chan *common.JobExecuteResult, 1000),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
	}

	go G_scheduler.scheduleLoop()
}

//注册任务
func (scheduler *Scheduler) RegisterJob(jobName string, cronExpr string, execute common.IJobExecute) {
	var (
		job      *common.Job
		jobEvent *common.JobEvent
	)

	job = &common.Job{
		Name:         jobName,
		JobExeTarget: execute,
		CronExpr:     cronExpr,
	}

	jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
	scheduler.jobEventChan <- jobEvent
}

//删除任务
func (scheduler *Scheduler) DelJob(jobName string) {

	var (
		job      *common.Job
		jobEvent *common.JobEvent
	)

	job = &common.Job{
		Name: jobName,
	}

	jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
	scheduler.jobEventChan <- jobEvent
}

//处理任务事件
func (scheduler *Scheduler) handlerJobEvent(event *common.JobEvent) {

	var (
		jobSchedulerPlan *common.JobSchedulePlan
		err              error
		jobIsExisted     bool
	)
	switch event.EventType {
	case common.JOB_EVENT_SAVE:

		if jobSchedulerPlan, err = common.BuildSchedulerPlan(event.Job); err != nil {
			return
		}

		scheduler.jobPlanTable[event.Job.Name] = jobSchedulerPlan
	case common.JOB_EVENT_DELETE:

		if jobSchedulerPlan, jobIsExisted = scheduler.jobPlanTable[event.Job.Name]; jobIsExisted {
			delete(scheduler.jobPlanTable, event.Job.Name)
		}
	}
}

//处理任务执行结果
func (scheduler *Scheduler) handlerJobResult(result *common.JobExecuteResult) {

	delete(scheduler.jobExecutingTable, result.JobExecuteInfo.Job.Name)
	//fmt.Println("任务执行完成", result.JobExecuteInfo.Job.Name, string(result.OutPut), result.Err)
}

//执行任务
func (scheduler *Scheduler) trySchedule() (scheduleAfterTime time.Duration) {

	var (
		now      time.Time
		jobPlan  *common.JobSchedulePlan
		nearTime *time.Time
	)

	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfterTime = 1 * time.Second
		return
	}

	now = time.Now()

	for _, jobPlan = range scheduler.jobPlanTable {

		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO:执行任务
			scheduler.tryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //设置下次执行时间
		}

		//统计最近一个要过期的任务时间（用于计算调度休眠时间）
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	//下次调度间隔（最近要执行的任务调度-当前时间）
	scheduleAfterTime = (*nearTime).Sub(now)
	return scheduleAfterTime
}

//开始执行单个任务计划
func (scheduler *Scheduler) tryStartJob(plan *common.JobSchedulePlan) {

	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)

	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[plan.Job.Name]; jobExecuting {
		fmt.Println("任务正在执行，跳过执行jobName=", plan.Job.Name)
		return
	}

	jobExecuteInfo = common.BuildJobExecuteInfo(plan)

	//保存任务执行信息
	scheduler.jobExecutingTable[plan.Job.Name] = jobExecuteInfo

	//执行任务
	fmt.Println("执行任务jobName=", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.executeJob(jobExecuteInfo)
}

func (scheduler *Scheduler) pushJobExecuteResult(result *common.JobExecuteResult) {
	scheduler.jobExeResultChan <- result
}

//调度协程
func (scheduler *Scheduler) scheduleLoop() {

	var (
		jobEvent          *common.JobEvent
		scheduleAfterTime time.Duration //调度间隔时间
		scheduleTimer     *time.Timer   //调度定时器
		exeResult         *common.JobExecuteResult
	)

	scheduleAfterTime = scheduler.trySchedule()
	scheduleTimer = time.NewTimer(scheduleAfterTime)

	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //任务事件监听处理
			scheduler.handlerJobEvent(jobEvent)
		case <-scheduleTimer.C: //调度执行任务
			scheduleAfterTime = scheduler.trySchedule()
			scheduleTimer.Reset(scheduleAfterTime)
		case exeResult = <-scheduler.jobExeResultChan:
			scheduler.handlerJobResult(exeResult)
		}
	}
}
