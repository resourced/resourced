package executors

import (
	"encoding/json"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func TestDynamicConstructor(t *testing.T) {
	config := resourced_config.Config{}
	config.GoStruct = "Shell"
	config.Kind = "executor"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Command"] = "uptime"

	s := NewGoStructByConfig(config)
	if s == nil {
		t.Fatalf("Shell constructor did not do its job")
	}

	// Test simple run, see if it works
	s.Run()

	inJson, err := s.ToJson()
	if err != nil {
		t.Errorf("Failed to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(inJson, &data)

	if data["Output"] == nil {
		t.Fatalf("There should always be output from uptime.")
	}
	if data["Output"].(string) == "" {
		t.Fatalf("There should always be output from uptime. Output: %v", data["Output"].(string))
	}

}
