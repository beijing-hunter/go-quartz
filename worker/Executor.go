package worker

import (
	"go-quartz/common"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func initExecutor() {
	G_executor = &Executor{}
}

//执行任务
func (executor *Executor) executeJob(info *common.JobExecuteInfo) {
	go func(info *common.JobExecuteInfo) {

		var (
			exeResult *common.JobExecuteResult
		)

		exeResult = &common.JobExecuteResult{
			JobExecuteInfo: info,
			OutPut:         make([]byte, 0),
		}

		exeResult.StartTime = time.Now()
		info.Job.JobExeTarget.Execute()
		exeResult.EndTime = time.Now()
		//推送执行结果
		G_scheduler.pushJobExecuteResult(exeResult)

	}(info)
}
