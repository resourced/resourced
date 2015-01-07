package writers

import (
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

	keysToTest := []string{"LoadAvg15m", "LoadAvg1m", "LoadAvg5m"}
	realData := n.InputData["Data"].(map[string]interface{})

	for _, key := range keysToTest {
		_, ok := realData[key]
		if !ok {
			t.Errorf("Key does not exist. Key: %v", key)
		}
	}
}
