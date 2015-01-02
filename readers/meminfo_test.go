// +build linux

package readers

import (
	"strings"
	"testing"
)

func TestNewMeminfo(t *testing.T) {
	n := NewMeminfo()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewMeminfoRun(t *testing.T) {
	n := NewMeminfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewMeminfoToJson(t *testing.T) {
	n := NewMeminfo()
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
	} else if !strings.Contains(string(jsonData), `MemTotal`) {
		t.Errorf("jsonData does not contain 'MemTotal' key. jsonData: %s", jsonData)
	}
}
