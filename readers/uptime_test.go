package readers

import (
	"strings"
	"testing"
)

func TestNewUptime(t *testing.T) {
	n := NewUptime()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

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

	if strings.Contains(string(jsonData), "Error") {
		t.Errorf("jsonData shouldn't return error: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `CurrentTimeUnixNano`) {
		t.Errorf("jsonData does not contain 'CurrentTimeUnixNano' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `CurrentTime`) {
		t.Errorf("jsonData does not contain 'CurrentTime' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Uptime`) {
		t.Errorf("jsonData does not contain 'Uptime' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `LoadAvg1m`) {
		t.Errorf("jsonData does not contain 'LoadAvg1m' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `LoadAvg5m`) {
		t.Errorf("jsonData does not contain 'LoadAvg5m' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `LoadAvg15m`) {
		t.Errorf("jsonData does not contain 'LoadAvg15m' key. jsonData: %s", jsonData)
	}
}
