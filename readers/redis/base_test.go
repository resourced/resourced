package redis

import (
	"strings"
	"testing"
)

func TestBaseInitConnection(t *testing.T) {
	r := &Base{}
	err := r.initConnection()

	if err != nil && strings.Contains(err.Error(), "connection refused") {
		t.Logf("Local Redis is not running. Stop testing.")
		return
	}

	if err != nil {
		t.Logf("Failed to initialize connection. Error: %v", err)
	}
}
