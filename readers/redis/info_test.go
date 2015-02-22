package redis

import (
	"testing"
)

func TestRedisInfoRun(t *testing.T) {
	r := NewRedisInfo()
	err := r.initConnection()
	if err != nil {
		t.Errorf("Initializing connection should always be successful. Error: %v", err)
	}
	if len(connections) == 0 {
		t.Errorf("Initializing connection should always be successful.")
	}

	err = r.Run()
	if err != nil {
		t.Errorf("fetching INFO data should always be successful. Error: %v", err)
	}

	if len(r.Data) == 0 {
		inJson, _ := r.ToJson()
		t.Errorf("fetching INFO data should always be successful. Data: %v", string(inJson))
	}
}
