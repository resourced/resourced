package readers

import (
	"encoding/json"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func TestShellRun(t *testing.T) {
	config := &resourced_config.Config{}
	config.GoStruct = "Shell"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Command"] = "$GOPATH/src/github.com/resourced/resourced/tests/data/script-reader/darwin-memory.py"
	config.Path = "/memory.darwin"
	config.Interval = "3s"

	s, err := NewGoStructByConfig(*config)
	if err != nil {
		t.Errorf("Creating IReader should work. Error: %v", err)
	}

	err = s.Run()
	if err != nil {
		t.Errorf("Running shell command should work. Error: %v", err)
	}

	inJson, err := s.ToJson()
	if err != nil {
		t.Errorf("Unable to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(inJson, &data)

	if data == nil {
		t.Errorf("Failed to run shell command. Data: %v", data)
	}
	if data["ExitStatus"].(float64) != 0 {
		t.Errorf("Failed to run shell command. Data: %v", data)
	}
}
