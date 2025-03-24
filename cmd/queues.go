package cmd

import (
	"fmt"
	"resque-inspector/models"
)

func getAllQueuesCmd() {
	result := models.GetQueueList(Filter)
	PrintJsonResult(result)

	fmt.Printf("Queues: %d\n\n", result.Total)
	for _, queue := range result.Items {
		fmt.Printf("\t%s: %d items\n", queue.Name, queue.JobCount)
	}
}
