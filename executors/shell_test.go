package executors

import (
	"encoding/json"
	"testing"
)

func TestSimpleRun(t *testing.T) {
	ResetConditionsMetByPath()

	config := newConfigExecutorForTest(t)

	s := NewGoStructByConfig(config)
	if s == nil {
		t.Fatalf("Shell constructor did not do its job")
	}

	s.Run()

	dataJson, err := s.ToJson()
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
