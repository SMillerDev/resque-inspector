package cmd

import (
	"resque-inspector/models"
)

func getAllJobsCmd(name string) {
	queue := models.GetQueue(name)

	var jobResult = queue.GetJobList(filterFromCmdline())
	PrintJsonResult(jobResult)
	PrintJobCmdResult(jobResult)
}
