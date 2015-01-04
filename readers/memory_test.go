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

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"Memory", "Swap", "ActualFree", "ActualUsed", "Used", "Free", "Total"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}
