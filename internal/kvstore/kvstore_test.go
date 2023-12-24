package kvstore

import (
	"testing"
)

func TestAddItem(t *testing.T) {
	//key := "Test"
	//value := "Value"

	_, err := NewKVStore()

	if err != nil {
		t.Fatalf("Error initializing KV store: %s", err)
	}

}
