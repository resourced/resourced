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

func TestNewNoopSetReadersData(t *testing.T) {
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
	readersData := make(map[string][]byte)
	readersData["/load-avg"] = []byte(jsonData)

	n.SetReadersData(readersData)

	key := "/load-avg"
	_, ok := n.ReadersData[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, n.ReadersData)
	}
}
