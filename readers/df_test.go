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
