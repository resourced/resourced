package writers

import (
	"encoding/json"
	"strings"
	"testing"
)

func readersDataForNewrelicInsightsTest() map[string]interface{} {
	jsonData := `{
    "Data": {
        "LoadAvg15m": 1.59375,
        "LoadAvg1m": 1.5537109375,
        "LoadAvg5m": 1.68798828125
    },
    "GoStruct": "LoadAvg",
    "Host": {
        "Name":"MacBook-Pro.local",
        "Tags":[]
    },
    "Hostname": "example.com",
    "Interval": "1s",
    "Path": "/load-avg",
    "UnixNano": 1420607791403576000
}`
	data := make(map[string]interface{})
	json.Unmarshal([]byte(jsonData), &data)
	return data
}

func TestNewNewrelicInsightsToJson(t *testing.T) {
	n := NewNewrelicInsights()

	readersData := make(map[string]interface{})
	readersData["/load-avg"] = readersDataForNewrelicInsightsTest()
	n.Data = readersData

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("ToJson() should never fail. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"LoadAvg15m", "LoadAvg1m", "LoadAvg5m", "Hostname", "eventType"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}
