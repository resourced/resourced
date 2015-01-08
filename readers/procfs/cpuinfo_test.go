package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcCpuInfo(t *testing.T) {
	n := NewProcCpuInfo()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewProcCpuInfoRun(t *testing.T) {
	n := NewProcCpuInfo()

	if runtime.GOOS == "linux" {
		err := n.Run()
		if err != nil {
			t.Errorf("Reading /proc/cpuinfo data should work on linux. Error: %v", err)
		}
	} else {
		err := n.Run()
		if err == nil {
			t.Error("Reading /proc/cpuinfo data should fail on non-linux.")
		}
	}
}

func TestNewProcCpuInfoToJson(t *testing.T) {
	n := NewProcCpuInfo()
	n.Run()

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/cpuinfo data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"vendor_id", "model", "model_name", "flags", "cores", "mhz"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
