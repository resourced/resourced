package writers

import (
	"encoding/json"
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
	h := NewHttp()

	readersData := make(map[string][]byte)
	readersData["/load-avg"] = jsonReadersDataForHttpTest()

	h.SetReadersData(readersData)

	return h
}

func TestNewHttpSetReadersData(t *testing.T) {
	h := newWriterForHttpTest()

	key := "/load-avg"
	_, ok := h.GetReadersData()[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, h.GetReadersData())
	}
}

func TestNewHttpRequest(t *testing.T) {
	h := newWriterForHttpTest()
	h.Url = "http://example.com/"
	h.Method = "POST"

	readersData := h.GetReadersData()
	dataJson, err := json.Marshal(readersData)

	if err != nil {
		t.Errorf("Failed to generate readers data. Error: %v", err)
	}

	_, err = h.NewHttpRequest(dataJson)
	if err != nil {
		t.Errorf("Failed to create Request struct. Error: %v", err)
	}
}

func TestNewHttpRun(t *testing.T) {
	h := newWriterForHttpTest()
	h.Url = "http://example.com/"
	h.Method = "POST"

	err := h.Run()
	if err != nil {
		t.Errorf("Run() should never fail. Error: %v", err)
	}
}
