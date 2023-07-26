package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/libis/urlchecker-extended/pkg/checker"
	"github.com/libis/urlchecker-extended/pkg/config"
	"github.com/libis/urlchecker-extended/pkg/slack"
)

func main() {
	// A Slack Webhook must be specified as an environment variable.
	webhook := os.Getenv(config.EnvSlackWebhook)
	if webhook == "" {
		log.Fatalf("Missing '%s' Environment Variable", config.EnvSlackWebhook)
	}
	slack := slack.SlackClient{
		Webhook: webhook,
	}

	// A filename must be specified for the program to read.
	var filename string
	flag.StringVar(&filename, "filename", "", "JSON File with paths")

	// A hostname must be specified since this may be used in different
	// environments.
	var hostname string
	flag.StringVar(&hostname, "hostname", "", "Hostname of website")

	var protocol string
	flag.StringVar(&protocol, "protocol", "https", "Protocol to use")

	var workers int
	flag.IntVar(&workers, "workers", 5, "Number of concurrent workers")

	var sleepFlag int
	flag.IntVar(&sleepFlag, "sleep", 0, "Number of seconds to sleep between requests")

	flag.Parse()
	if filename == "" {
		log.Fatal("Missing filename flag")
	}
	if hostname == "" {
		log.Fatal("Missing hostname flag")
	}

	// Attempt to read the file specified.
	log.Printf("Reading %s...\n", filename)

	sleep := time.Duration(sleepFlag)

	// Attempt to parse the file content as JSON.
	checker.Check(filename, protocol, hostname, slack, workers, sleep)
}
