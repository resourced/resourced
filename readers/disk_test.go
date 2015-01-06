package readers

import (
	"strings"
	"testing"
)

func TestNewDiskPartitions(t *testing.T) {
	n := NewDiskPartitions()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewDiskPartitionsRun(t *testing.T) {
	n := NewDiskPartitions()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing df data should always be successful. Error: %v", err)
	}
}

func TestNewDiskPartitionsToJson(t *testing.T) {
	n := NewDiskPartitions()
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

	keysToTest := []string{"device", "mountpoint", "fstype", "opts"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}
