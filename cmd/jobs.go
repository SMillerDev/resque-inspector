package cmd

import (
	"os"
	"resque-inspector/models"
	"strconv"
)

func getAllJobsCmd(queue string) {
	result := models.GetQueueList(queue)
	if result.Total == 0 {
		println("no such queue found: " + queue)
		os.Exit(1)
	}
	if result.Total > 1 {
		println(strconv.Itoa(result.Total) + " queues found for: " + queue)
		os.Exit(1)
	}

	var jobResult = result.Items[0].GetJobList(Filter)
	PrintJsonResult(jobResult)
	PrintJobCmdResult(jobResult)
}
