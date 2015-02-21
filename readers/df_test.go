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

func TestDfRun(t *testing.T) {
	n := NewDf()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}
}

func TestDfToJson(t *testing.T) {
	n := NewDf()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling df data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"Available", "DeviceName", "Total", "UsePercent", "Used"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}

func TestDfFSPaths(t *testing.T) {
	n := NewDf()
	n.FSPaths = "/tmp,/bin"

	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}

	if len(n.Data) == 0 {
		inJson, _ := n.ToJson()
		t.Errorf("df data should not be empty. Data: %v", string(inJson))
	}
	if len(n.Data["/tmp"]) == 0 {
		inJson, _ := n.ToJson()
		t.Errorf("df /tmp data should not be empty. Data: %v", string(inJson))
	}
	if len(n.Data["/bin"]) == 0 {
		inJson, _ := n.ToJson()
		t.Errorf("df /bin data should not be empty. Data: %v", string(inJson))
	}
}
