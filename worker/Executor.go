package worker

import (
	"go-quartz/common"
	"sync"
	"time"
)

type Executor struct {
}

var (
	G_executor   *Executor
	executorLock sync.Mutex
)

func initExecutor() {

	executorLock.Lock()
	defer executorLock.Unlock()

	if G_executor == nil {
		G_executor = &Executor{}
	}
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

		if info.Job.JobExeTarget != nil {
			info.Job.JobExeTarget.Execute() //执行任务对象
		} else if info.Job.ExeHandler != nil {
			info.Job.ExeHandler() //执行任务处理函数
		}

		exeResult.EndTime = time.Now()
		//推送执行结果
		G_scheduler.pushJobExecuteResult(exeResult)

	}(info)
}
