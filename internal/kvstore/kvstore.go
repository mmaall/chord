package kvstore

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

var kvLogger = log.WithFields(log.Fields{"process": "kvStoreLogger"})

type KVStore struct {
	kvMap map[string]string
}

func NewKVStore() (*KVStore, error) {
	return &KVStore{
		kvMap: make(map[string]string),
	}, nil
}

// Get a value from the key value store
func (kvStore *KVStore) Get(key string) string {
	kvLogger.Infof("Get %s", key)
	return kvStore.kvMap[key]
}

// Put a value to the key value store
func (kvStore *KVStore) Put(key string, value string) error {
	kvLogger.Infof("Putting %s to key %s", value, key)
	kvStore.kvMap[key] = value
	return nil
}

func (kvStore *KVStore) ToString() string {

	return fmt.Sprintf("%s", kvStore.kvMap)

}
