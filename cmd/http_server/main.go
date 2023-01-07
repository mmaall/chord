package main

import (
	"chord/internal/linkedlistnode"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {

	logger := log.WithFields(log.Fields{"process": "main"})

	const address string = "localhost:8080"

	logger.Infof("Lets get chord going\n")

	var processWaitGroup sync.WaitGroup

	nodeServer, err := linkedlistnode.NewNode(address)
	nodeServer.AddWaitGroup(&processWaitGroup)

	if err != nil {
		logger.Errorf("Error starting chord node\n")
		logger.Errorf("%d\n", err)
	}

	nodeServer.Start()

	processWaitGroup.Wait()
}
