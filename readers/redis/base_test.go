package redis

import (
	"testing"
)

func TestBaseInitConnection(t *testing.T) {
	r := &Base{}
	err := r.initConnection()
	if err != nil {
		t.Errorf("Initializing connection should always be successful. Error: %v", err)
	}
	if len(connections) == 0 {
		t.Errorf("Initializing connection should always be successful.")
	}
}
