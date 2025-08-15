package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"resque-inspector/models"
	"resque-inspector/resque"
	"resque-inspector/server"
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

func filterFromCmdline() models.Filter {
	return models.Filter{}
}

func setupJson() {
	jsonOut = baseJsonOut || subJsonOut
}
func setupDebug() {
	debug = baseDebug || subDebug
	resque.Debug = debug
	models.Debug = debug
}
func setupDsn() {
	var dsn string

	if baseDsnFlag != "" {
		dsn = baseDsnFlag
	} else if subDsnFlag != "" {
		dsn = subDsnFlag
	} else if baseHost != defaultRedisHost && basePort != defaultRedisPort {
		dsn = dsnFromHostPort(baseHost, basePort)
	} else {
		dsn = dsnFromHostPort(subHost, subPort)
	}

	envDsn := parseEnvironmentForDSN()
	if envDsn != "" {
		dsn = envDsn
	}

	resque.Dsn = dsn
	server.Dsn = dsn
}

func dsnFromHostPort(host string, port int) string {
	return fmt.Sprintf("redis://%s:%d", host, port)
}
