package cmd

import (
	"fmt"
	"resque-inspector/models"
)

func getAllQueuesCmd() {
	result := models.GetQueueList(filterFromCmdline())
	PrintJsonResult(result)

	fmt.Printf("Classes: %d\n\n", result.Total)
	for _, queue := range result.Items {
		fmt.Printf("\t%s: %d items\n", queue.Name, queue.JobCount)
	}
}

func clearQueueCmd(queue string) {
	err := models.GetQueue(queue).Clear()
	if err != nil {
		fmt.Printf("Error clearing queue: %s\n", err)
	}
	fmt.Printf("Cleared queue: %s\n", queue)
}
