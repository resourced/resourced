package readers

import (
	"strings"
	"testing"
)

func TestNewPs(t *testing.T) {
	n := NewPs()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewPsRun(t *testing.T) {
	n := NewPs()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}
}

func TestNewPsToJson(t *testing.T) {
	n := NewPs()
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
	} else if !strings.Contains(string(jsonData), `Name`) {
		t.Errorf("jsonData does not contain 'Name' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Pid`) {
		t.Errorf("jsonData does not contain 'Pid' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `ParentPid`) {
		t.Errorf("jsonData does not contain 'ParentPid' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `StartTime`) {
		t.Errorf("jsonData does not contain 'StartTime' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `RunTime`) {
		t.Errorf("jsonData does not contain 'RunTime' key. jsonData: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `MemoryResident`) {
		t.Errorf("jsonData does not contain 'MemoryResident' key. jsonData: %s", jsonData)
	}
}
