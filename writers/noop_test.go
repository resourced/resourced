package writers

import (
	"strings"
	"testing"
)

func TestNewNoopRun(t *testing.T) {
	n := NewNoop()
	err := n.Run()
	if err != nil {
		t.Errorf("Run() should never fail. Error: %v", err)
	}
}

func TestNewNoopSetJsonData(t *testing.T) {
	n := NewNoop()

	jsonData := `{
    "Data": {
        "LoadAvg15m": 1.59375,
        "LoadAvg1m": 1.5537109375,
        "LoadAvg5m": 1.68798828125
    },
    "GoStruct": "LoadAvg",
    "Hostname": "example.com",
    "Interval": "1s",
    "Path": "/load-avg",
    "Tags": [ ],
    "UnixNano": 1420607791403576000
}`

	err := n.SetData([]byte(jsonData))
	if err != nil {
		t.Errorf("Marshalling data should be successful. Error: %v", err)
	}

	jsonDataFromStruct, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling data should be successful. Error: %v", err)
	}

	jsonDataFromStructString := string(jsonDataFromStruct)

	if strings.Contains(jsonDataFromStructString, "Error") {
		t.Errorf("jsonDataFromStructString shouldn't return error: %v", jsonDataFromStructString)
	}

	keysToTest := []string{"LoadAvg15m", "LoadAvg1m", "LoadAvg5m"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataFromStructString, key) {
			t.Errorf("jsonDataFromStructString does not contain '%v' key. jsonDataFromStructString: %v", key, jsonDataFromStructString)
		}
	}
}
