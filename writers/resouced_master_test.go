package writers

import (
	"testing"
)

func jsonReadersDataForResourcedMasterTest() []byte {
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
	return []byte(jsonData)
}

func newWriterForResourcedMasterTest() *ResourcedMaster {
	n := NewResourcedMaster()

	readersData := make(map[string][]byte)
	readersData["/load-avg"] = jsonReadersDataForResourcedMasterTest()

	n.SetReadersDataInBytes(readersData)

	return n
}

func TestNewResourcedMasterSetReadersDataInBytes(t *testing.T) {
	n := newWriterForResourcedMasterTest()

	key := "/load-avg"
	_, ok := n.GetReadersData()[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, n.GetReadersData())
	}
}

func TestNewResourcedMasterRun(t *testing.T) {
	n := newWriterForResourcedMasterTest()
	n.Url = "http://localhost:55655/"
	n.Method = "POST"

	err := n.Run()
	if err != nil {
		t.Errorf("Run() should never fail. Error: %v", err)
	}
}
