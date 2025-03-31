package cmd

import (
	"encoding/json"
	"log"
	"os"
	"resque-inspector/models"
	"resque-inspector/result"
)

func PrintJsonResult(data interface{}) {
	if !jsonOut {
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		log.Default().Fatalf("Could not marshal json: %s\n", err)
	}

	print(string(jsonData))
	os.Exit(0)
}

func PrintJobCmdResult(data result.Result[models.JobInterface]) {
	for _, val := range data.Items {
		println(val.Stringify())
	}
}

func filterFromCmdline() result.Filter {
	return result.Filter{}
}
