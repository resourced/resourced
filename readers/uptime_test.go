package readers

import (
	"strings"
	"testing"
)

func TestNewUptimeRun(t *testing.T) {
	n := NewUptime()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing load average data should always be successful. Error: %v", err)
	}
}

func TestNewUptimeToJson(t *testing.T) {
	n := NewUptime()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing load average data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling load average data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"CurrentTimeUnixNano", "CurrentTime", "Uptime", "LoadAvg1m", "LoadAvg5m", "LoadAvg15m", "TimeZone"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}
