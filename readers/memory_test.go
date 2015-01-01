package readers

import (
	"strings"
	"testing"
)

func TestNewMemory(t *testing.T) {
	n := NewMemory()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewMemoryRun(t *testing.T) {
	n := NewMemory()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewMemoryToJson(t *testing.T) {
	n := NewMemory()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling memory data should always be successful. Error: %v", err)
	}

	if strings.Contains(string(jsonData), "Error") {
		t.Errorf("jsonData shouldn't return error: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Memory`) {
		t.Errorf("jsonData does not contain 'Memory' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Swap`) {
		t.Errorf("jsonData does not contain 'Swap' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `ActualFree`) {
		t.Errorf("jsonData does not contain 'ActualFree' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `ActualUsed`) {
		t.Errorf("jsonData does not contain 'ActualUsed' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Used`) {
		t.Errorf("jsonData does not contain 'Used' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Free`) {
		t.Errorf("jsonData does not contain 'Free' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Total`) {
		t.Errorf("jsonData does not contain 'Total' key. jsonData: %s", jsonData)
	}
}
