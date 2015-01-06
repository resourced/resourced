package readers

import (
	"runtime"
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

// ------------------------------------------------------------

func TestNewDiskIO(t *testing.T) {
	n := NewDiskIO()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

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

		keysToTest := []string{"read_count", "write_count", "read_bytes", "write_bytes", "read_time", "write_time", "name", "io_time", "serial_number"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
