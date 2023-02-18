package main

import (
	"chord/internal/linkedlistnode"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

func main() {

	logger := log.WithFields(log.Fields{"process": "main"})

	// Parse command line arguments
	var (
		help    bool
		address = "localhost:8080"
	)

	getopt.FlagLong(&help, "help", 'h', "Help")
	getopt.FlagLong(&address, "address", 'f', "Where the server is listening. Ex. localhost:8080")

	getopt.Parse()

	if help {
		getopt.Usage()
		os.Exit(0)
	}

	logger.Infof("Lets get chord going\n")

	var processWaitGroup sync.WaitGroup

	// Initialize server and register wait group
	nodeServer, err := linkedlistnode.NewNode(address)
	nodeServer.AddWaitGroup(&processWaitGroup)

	if err != nil {
		logger.Errorf("Error starting chord node\n")
		logger.Errorf("%d\n", err)
	}

	// Start server
	nodeServer.Start()

	// Wait for server to exit

	processWaitGroup.Wait()
}
