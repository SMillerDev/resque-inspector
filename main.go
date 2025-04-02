package main

import (
	"resque-inspector/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.ParseCommandLine(version, date)
}
