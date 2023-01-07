package linkedlistnode

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

var listNodeLogger = log.WithFields(log.Fields{"process": "listNodeLogger"})

type LinkedlistNode struct {
	httpServer *http.Server
	address    string
	waitGroup  *sync.WaitGroup
}

func NewNode(address string) (*LinkedlistNode, error) {
	return &LinkedlistNode{
		httpServer: nil,
		address:    address,
		waitGroup:  nil,
	}, nil
}

// Start a node
func (node *LinkedlistNode) Start() (bool, error) {

	listNodeLogger.Infof("Starting linked list node")

	// Set wait group
	if node.waitGroup != nil {
		node.waitGroup.Add(1)
	}

	// Set address
	node.httpServer = &http.Server{Addr: node.address}

	// Initialize handlers
	http.HandleFunc("/hello", pingHandler)

	// Serve
	go func() {
		listNodeLogger.Infof("Serving node on %s\n", node.address)

		err := node.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			listNodeLogger.Errorf("Unexpected stop on Http Server: %v\n", err)
			node.Shutdown()
		} else {
			listNodeLogger.Infof("Http Server closed %v\n", err)
		}
	}()
	return true, nil
}

func (node *LinkedlistNode) AddWaitGroup(wg *sync.WaitGroup) {
	node.waitGroup = wg
}

// Shutdown terminates the HTTP server listening for logs
func (node *LinkedlistNode) Shutdown() {

	// Shutdown server
	if node.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		err := node.httpServer.Shutdown(ctx)
		if err != nil {
			listNodeLogger.Errorf("Failed to shutdown http server gracefully %s\n", err)
		} else {
			node.httpServer = nil
		}
	}

	// Stop waitgroup
	if node.waitGroup != nil {
		node.waitGroup.Done()
	}

}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	listNodeLogger.Infof("Received a ping\n")
	io.WriteString(w, "Hello, world!\n")
}
