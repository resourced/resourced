package writers

import (
	"encoding/json"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func TestShellRun(t *testing.T) {
	config := &resourced_config.Config{}
	config.GoStruct = "Shell"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Command"] = "$GOPATH/src/github.com/resourced/resourced/tests/script-writer/stdin-stdout.py"
	config.ReaderPaths = []string{"/load-avg", "/free"}
	config.Path = "/python/loadavg-uptime-free"
	config.Interval = "1m"

	s, err := NewGoStructByConfig(*config)
	if err != nil {
		t.Errorf("Creating IWriter should work. Error: %v", err)
	}

	readersData := make(map[string]map[string]map[string]interface{})
	readersData["/load-avg"] = make(map[string]map[string]interface{})
	readersData["/load-avg"]["Data"] = make(map[string]interface{})
	readersData["/load-avg"]["Data"]["LoadAvg1m"] = 1.97119140625
	readersData["/free"] = make(map[string]map[string]interface{})
	readersData["/free"]["Data"] = make(map[string]interface{})
	readersData["/free"]["Data"]["Used"] = 8563695616

	readersDataJson, _ := json.Marshal(readersData)
	var readersDataMapInterface map[string]interface{}
	json.Unmarshal(readersDataJson, &readersDataMapInterface)

	s.SetReadersData(readersDataMapInterface)

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
