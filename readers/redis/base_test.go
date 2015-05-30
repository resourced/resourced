package redis

import (
	"testing"
)

func TestBaseInitConnection(t *testing.T) {
	r := &Base{}
	err := r.initConnection()
	if err != nil {
		t.Logf("Failed to initialize connection. Error: %v", err)
	}
}
