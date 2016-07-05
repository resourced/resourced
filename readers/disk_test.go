package readers

import (
	"runtime"
	"strings"
	"testing"
)

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

// ------------------------------------------------------------

func TestNewDiskIORun(t *testing.T) {
	if runtime.GOOS == "linux" {
		n := NewDiskIO()
		err := n.Run()
		if err != nil {
			t.Errorf("Parsing data should always be successful. Error: %v", err)
		}
	}
}

func TestNewDiskIOToJson(t *testing.T) {
	if runtime.GOOS == "linux" {
		n := NewDiskIO()
		err := n.Run()
		if err != nil {
			t.Errorf("Parsing data should always be successful. Error: %v", err)
		}

		jsonData, err := n.ToJson()
		if err != nil {
			t.Errorf("Marshalling data should always be successful. Error: %v", err)
		}

		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"readCount", "writeCount", "readBytes", "writeBytes", "readTime", "writeTime", "name", "ioTime", "serialNumber"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
