package executors

import (
	"encoding/json"
	"testing"

	"github.com/resourced/resourced/libmap"
)

func TestShellRun(t *testing.T) {
	config := newConfigExecutorForTest(t)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	counterDB := libmap.NewTSafeMapCounter(nil)
	executor.SetCounterDB(counterDB)

	executor.Run()

	dataJson, err := executor.ToJson()
	if err != nil {
		t.Fatalf("Failed to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(dataJson, &data)

	if data["Output"] == nil {
		t.Fatalf("There should always be output from uptime.")
	}
	if data["Output"].(string) == "" {
		t.Fatalf("There should always be output from uptime. Output: %v", data["Output"].(string))
	}
}
