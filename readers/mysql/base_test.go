package mysql

import (
	"strings"
	"testing"
)

func TestBaseInitConnection(t *testing.T) {
	m := &Base{}
	m.Retries = 1

	err := m.initConnection()
	if strings.Contains(err.Error(), "connection refused") {
		t.Infof("Local MySQL is not running. Stop testing.")
		return
	}

	if err != nil {
		t.Errorf("Initializing connection should always be successful. Error: %v", err)
	}
	if len(connections) == 0 {
		t.Errorf("Initializing connection should always be successful.")
	}
}
