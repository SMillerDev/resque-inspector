package cmd

import (
	"fmt"
	"resque-inspector/models"
)

func getAllWorkersCmd() {
	result := models.GetWorkerList(filterFromCmdline())

	PrintJsonResult(result)

	fmt.Printf("Workers: %d\n", result.Total)
	for workerName, slots := range result.Items {
		fmt.Printf("Worker: %s (%d slots)\n", workerName, len(slots))
		for _, slot := range slots {
			fmt.Printf("\tHost: %s, Socket: %s\n", slot.Host, slot.Socket)
		}
	}
}
