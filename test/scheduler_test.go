package test

import (
	"fmt"
	"go-quartz/worker"
	"testing"
	"time"
)

type JobA struct {
	JobName string
}

type JobB struct {
	JobName string
}

func (joba *JobA) Execute() {

	fmt.Println("exe task jobName ", joba.JobName)
}

func (jobb *JobB) Execute() {

	time.Sleep(7 * time.Second)
	fmt.Println("exe task jobName ", jobb.JobName)
}

func TestAddJobObject(t *testing.T) {

	var (
		jobA *JobA
		jobB *JobB
	)

	worker.Init()

	jobA = &JobA{
		JobName: "JobA",
	}

	jobB = &JobB{
		JobName: "JobB",
	}

	worker.G_scheduler.RegisterJobObject(jobA.JobName, "*/6 * * * * * *", jobA)
	worker.G_scheduler.RegisterJobObject(jobB.JobName, "*/6 * * * * * *", jobB)

	//for {

		//time.Sleep(1 * time.Second)
	//}
}

func TestAddJobHandler(t *testing.T) {

	worker.Init()

	worker.G_scheduler.RegisterJobHandler("jobE", "*/8 * * * * * *", func() {
		fmt.Println("exe task jobName", "jobE")
	})

	worker.G_scheduler.RegisterJobHandler("jobF", "*/10 * * * * * *", func() {
		time.Sleep(11 * time.Second)
		fmt.Println("exe task jobName", "jobF")
	})

	for {

		time.Sleep(1 * time.Second)
	}
}
