package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"resque-inspector/models"
	"resque-inspector/resque"
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

func PrintJobCmdResult(data resque.Result[models.JobInterface]) {
	fmt.Printf("Printing %d out of %d jobs\n", data.Selected, data.Total)
	for _, val := range data.Items {
		println(val.Stringify())
	}
}

func filterFromCmdline() resque.Filter {
	return resque.Filter{}
}
