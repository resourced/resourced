package writers

import (
	"testing"
)

func jsonReadersDataForHttpTest() []byte {
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

func newWriterForHttpTest() *Http {
	n := NewHttp()

	readersData := make(map[string][]byte)
	readersData["/load-avg"] = jsonReadersDataForHttpTest()

	n.SetReadersData(readersData)

	return n
}

func TestNewHttpSetReadersData(t *testing.T) {
	n := newWriterForHttpTest()

	key := "/load-avg"
	_, ok := n.GetReadersData()[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, n.GetReadersData())
	}
}

func TestNewHttpRun(t *testing.T) {
	n := newWriterForHttpTest()
	n.Url = "http://example.com/"
	n.Method = "POST"

	err := n.Run()
	if err != nil {
		t.Errorf("Run() should never fail. Error: %v", err)
	}
}
