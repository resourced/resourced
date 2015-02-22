package mysql

import (
	"testing"
)

func TestBaseInitConnection(t *testing.T) {
	m := &Base{}
	err := m.initConnection()
	if err != nil {
		t.Errorf("Initializing connection should always be successful. Error: %v", err)
	}
	if len(connections) == 0 {
		t.Errorf("Initializing connection should always be successful.")
	}
}
