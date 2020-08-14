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

func TestAddJob(t *testing.T) {

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

	worker.G_scheduler.RegisterJob(jobA.JobName, "*/6 * * * * * *", jobA)
	worker.G_scheduler.RegisterJob(jobB.JobName, "*/6 * * * * * *", jobB)

	for {

		time.Sleep(1 * time.Second)
	}
}
