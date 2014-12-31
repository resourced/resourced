package readers

import (
	"strings"
	"testing"
)

func TestNewNetworkInterfaces(t *testing.T) {
	n := NewNetworkInterfaces()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewNetworkInterfacesRun(t *testing.T) {
	n := NewNetworkInterfaces()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing net.Interfaces() data should always be successful. Error: %v", err)
	}
}

func TestNewNetworkInterfacesToJson(t *testing.T) {
	n := NewNetworkInterfaces()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing net.Interfaces() data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling net.Interfaces() data should always be successful. Error: %v", err)
	}

	if strings.Contains(string(jsonData), "Error") {
		t.Errorf("jsonData shouldn't return error: %s", jsonData)
	} else if !strings.Contains(string(jsonData), `Addresses`) {
		t.Errorf("jsonData does not contain 'Addresses' key. jsonData: %s", jsonData)
	}
}
