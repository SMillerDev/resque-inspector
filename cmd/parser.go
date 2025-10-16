package cmd

import (
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"resque-inspector/resque"
	"resque-inspector/server"
	"slices"
	"strconv"
)

func ParseCommandLine(version string, date string) {
	subcommands := setupSubCommands()
	if len(flag.Args()) < 1 {
		log.Default().Fatal("expected a subcommands, got none!")
	}

	subcommand := subcommands[os.Args[1]]
	if subcommand == nil {
		flag.Usage()
		os.Exit(1)
	}

	err := subcommand.Parse(os.Args[2:])
	if err != nil {
		log.Default().Printf("Failed to parse command %s\n", err)
		os.Exit(1)
	}

	setupDsn()
	setupJson()
	setupDebug()

	switch os.Args[1] {
	case "queues":
		getAllQueuesCmd()
	case "workers":
		getAllWorkersCmd()
	case "jobs":
		if len(subcommand.Args()) == 0 {
			log.Default().Fatal("expected queue name to fetch jobs for, got none")
		}

		getAllJobsCmd(subcommand.Arg(0))
	case "clear":
		if len(subcommand.Args()) == 0 {
			log.Default().Fatal("expected queue to clear, got none")
		}

		clearQueueCmd(subcommand.Arg(0))
	case "retry":
		log.Default().Println("subcommand 'retry' is not yet implemented!")
	case "serve":
		server.Serve()
	case "version":
		fmt.Printf("resque-inspector %s, built at %s", version, date)
	default:
		log.Default().Printf("expected one of %v subcommands\n", maps.Keys(subcommands))
	}

	if resque.Client != nil {
		defer resque.Client.Close()
	}
}

func setupSubCommands() map[string]*flag.FlagSet {
	readSubcommands := []string{"queues", "jobs", "workers"}
	deleteSubcommands := []string{"clear"}
	writeSubcommands := []string{"retry"}
	utilSubcommands := []string{"serve", "version"}
	allSubCommands := slices.Concat(readSubcommands, deleteSubcommands, writeSubcommands, utilSubcommands)

	subcommands := map[string]*flag.FlagSet{}
	for _, subcommand := range allSubCommands {
		subcommands[subcommand] = flag.NewFlagSet(subcommand, flag.ExitOnError)
	}
	for _, subcommand := range readSubcommands {
		subcommands[subcommand].StringVar(&Filter, "filter", defaultFilter, "Regex to filter "+subcommand)
	}

	flag.StringVar(&baseHost, "host", defaultRedisHost, "The hostname to connect to")
	flag.IntVar(&basePort, "port", defaultRedisPort, "The port to connect to")

	flag.StringVar(&baseDsnFlag, "dsn", "", "DSN to connect to, will override the hostname and port")
	flag.BoolVar(&baseJsonOut, "json", false, "Output content in JSON format")
	flag.BoolVar(&baseDebug, "debug", false, "Print debug output")

	flag.Parse()

	for _, set := range subcommands {
		set.StringVar(&subHost, "host", defaultRedisHost, "The hostname to connect to")
		set.IntVar(&subPort, "port", defaultRedisPort, "The port to connect to")

		set.StringVar(&subDsnFlag, "dsn", "", "DSN to connect to, will override the hostname and port")
		set.BoolVar(&subJsonOut, "json", false, "Output content in JSON format")
		set.BoolVar(&subDebug, "debug", false, "Print debug output")
	}

	return subcommands
}

func parseEnvironmentForDSN() string {
	val, isset := os.LookupEnv("REDIS_DSN")
	if isset {
		return val
	}
	host, issetHost := os.LookupEnv("REDIS_HOST")
	port, issetPort := os.LookupEnv("REDIS_PORT")
	if issetHost || issetPort {
		if host == "" {
			host = defaultRedisHost
		}
		if port == "" {
			port = strconv.Itoa(defaultRedisPort)
		}
		iPort, _ := strconv.Atoi(port)
		return dsnFromHostPort(host, iPort)
	}

	if debug {
		log.Default().Println("No environment variable set, using default")
	}
	return ""
}
