package linkedlistnode

import (
	"chord/internal/kvstore"
	"context"
	"encoding/json"
	"fmt"
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
	kvStore    *kvstore.KVStore
}

func NewNode(address string) (*LinkedlistNode, error) {

	kvStore, err := kvstore.NewKVStore()

	if err != nil {
		return nil, err
	}
	return &LinkedlistNode{
		httpServer: nil,
		address:    address,
		waitGroup:  nil,
		kvStore:    kvStore,
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
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/put", node.putHandler)

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

// Put data into the key value store
// JSON format
// {"key" : "the_key", "value" : "the_value"}
func (node *LinkedlistNode) putHandler(w http.ResponseWriter, req *http.Request) {
	listNodeLogger.Infof("Received a put request")

	rawRequestBody, err := io.ReadAll(req.Body)
	if err != nil {
		listNodeLogger.Errorf("Error reading request body\n")
		listNodeLogger.Errorf("Body: %s\n", rawRequestBody)
		listNodeLogger.Errorf("%s", err)
		Error(w, "Internal server error", 500)
		return
	}

	requestBody := make(map[string]string)

	err = json.Unmarshal(rawRequestBody, &requestBody)

	if err != nil {
		listNodeLogger.Errorf("Error unmarshaling JSON\n")
		listNodeLogger.Errorf("Body: %s\n", rawRequestBody)
		listNodeLogger.Errorf("%s\n", err)
		Error(w, "Malformed request body. Is it valid JSON?", 400)
		return
	}

	requestBodyFormatted, _ := json.Marshal(requestBody)
	listNodeLogger.Infof("Body: %s\n", requestBodyFormatted)

	if requestBody["key"] == "" || requestBody["value"] == "" {
		listNodeLogger.Errorf("Missing key or value in input JSON\n")
		Error(w, "Either key or value not included", 400)
	}

	err = node.kvStore.Put(requestBody["key"], requestBody["value"])

	if err != nil {
		listNodeLogger.Errorf("Error writing to key value store\n")
		listNodeLogger.Errorf("%s\n", err)
		Error(w, "Internal server error", 500)
	}

	allData := node.kvStore.ToString()
	listNodeLogger.Infof("All Data: %s\n", allData)

}

// Ping (pong)
func pingHandler(w http.ResponseWriter, req *http.Request) {
	listNodeLogger.Infof("Received a ping\n")
	listNodeLogger.Infof("%s\t%s\n", req.Method, req.URL.EscapedPath())
	requestBody, err := io.ReadAll(req.Body)

	// Write out body if included
	if err == nil {
		listNodeLogger.Infof("Body: %s\n", requestBody)
	}

	io.WriteString(w, "pong\n")
}

func Error(w http.ResponseWriter, error string, code int) {
	responseBodyMap := make(map[string]string)

	responseBodyMap["statusCode"] = fmt.Sprintf("%d", code)
	responseBodyMap["statusMessage"] = error

	responseBody, err := json.Marshal(responseBodyMap)

	if err != nil {
		listNodeLogger.Errorf("Error sending HTTP error response.\n")
		http.Error(w, error, 500)
		return
	}

	w.WriteHeader(code)
	io.WriteString(w, fmt.Sprintf("%s", responseBody))

}
