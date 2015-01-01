package readers

import (
	"strings"
	"testing"
)

func TestNewDf(t *testing.T) {
	n := NewDf()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewDfRun(t *testing.T) {
	n := NewDf()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}
}

func TestNewDfToJson(t *testing.T) {
	n := NewDf()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling df data should always be successful. Error: %v", err)
	}

	if strings.Contains(string(jsonData), "Error") {
		t.Errorf("jsonData shouldn't return error: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Available`) {
		t.Errorf("jsonData does not contain 'Available' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `DeviceName`) {
		t.Errorf("jsonData does not contain 'DeviceName' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Total`) {
		t.Errorf("jsonData does not contain 'Total' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `UsePercent`) {
		t.Errorf("jsonData does not contain 'UsePercent' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Used`) {
		t.Errorf("jsonData does not contain 'Used' key. jsonData: %s", jsonData)
	}
}
