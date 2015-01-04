package readers

import (
	"strings"
	"testing"
)

func TestNewDu(t *testing.T) {
	n := NewDu()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewDuRun(t *testing.T) {
	n := NewDu()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}
}

func TestNewDuToJson(t *testing.T) {
	n := NewDu()
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
	} else if !strings.Contains(string(jsonData), `Path`) {
		t.Errorf("jsonData does not contain 'Path' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Total`) {
		t.Errorf("jsonData does not contain 'Total' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Free`) {
		t.Errorf("jsonData does not contain 'Free' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `InodesTotal`) {
		t.Errorf("jsonData does not contain 'InodesTotal' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `InodesFree`) {
		t.Errorf("jsonData does not contain 'InodesFree' key. jsonData: %s", jsonData)
	}
}
