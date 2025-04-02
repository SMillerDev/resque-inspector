package cmd

import (
	"resque-inspector/models"
)

func getAllJobsCmd(name string) {
	queue := models.GetQueue(name)

	var jobResult = queue.GetJobList(filterFromCmdline(), 0, 100)
	PrintJsonResult(jobResult)
	PrintJobCmdResult(jobResult)
}
