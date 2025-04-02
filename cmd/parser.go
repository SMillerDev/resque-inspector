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

var Filter string

var jsonOut bool
var subJsonOut bool
var baseJsonOut bool

var subDsnFlag string
var baseDsnFlag string
var subHost string
var baseHost string
var subPort int
var basePort int

const defaultRedisPort = 6379
const defaultRedisHost = "127.0.0.1"
const defaultFilter = ".*"

func ParseCommandLine(version string, date string) {
	subcommands := setupSubCommands()
	if len(flag.Args()) < 1 {
		log.Default().Fatal("expected a subcommands, got none!")
	}

	subcommand := subcommands[os.Args[1]]
	if subcommand == nil {
		for _, subcommand := range subcommands {
			fmt.Printf("%s\n", subcommand.Name())
		}
		flag.Usage()
		os.Exit(1)
	}

	err := subcommand.Parse(os.Args[2:])
	if err != nil {
		log.Default().Printf("Failed to parse command %w\n", err)
		os.Exit(1)
	}

	setupDsn()
	setupJson()

	switch os.Args[1] {
	case "queues":
		getAllQueuesCmd()
	case "workers":
		getAllWorkersCmd()
	case "jobs":
		if len(subcommand.Args()) == 0 {
			println("expected queue Name to fetch jobs for, got none")
			os.Exit(1)
		}

		getAllJobsCmd(subcommand.Arg(0))
	case "clear":
		fmt.Println("subcommand 'clear'")
	case "retry":
		log.Default().Println("subcommand 'retry'")
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

func setupJson() {
	jsonOut = baseJsonOut || subJsonOut
}
func setupDsn() {
	var dsn string

	if baseDsnFlag != "" {
		dsn = baseDsnFlag
	} else if subDsnFlag != "" {
		dsn = subDsnFlag
	} else if baseHost != defaultRedisHost && basePort != defaultRedisPort {
		dsn = fmt.Sprintf("%s:%d", baseHost, basePort)
	} else {
		dsn = fmt.Sprintf("%s:%d", subHost, subPort)
	}

	envDsn := parseEnvironmentForDSN()
	if envDsn != "" {
		dsn = envDsn
	}

	resque.Dsn = dsn
	server.Dsn = dsn
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

	flag.Parse()

	for _, set := range subcommands {
		set.StringVar(&subHost, "host", defaultRedisHost, "The hostname to connect to")
		set.IntVar(&subPort, "port", defaultRedisPort, "The port to connect to")

		set.StringVar(&subDsnFlag, "dsn", "", "DSN to connect to, will override the hostname and port")
		set.BoolVar(&subJsonOut, "json", false, "Output content in JSON format")
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
		return fmt.Sprintf("%s:%s", host, port)
	}

	log.Default().Println("No environment variable set, using default")
	return ""
}
