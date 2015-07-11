package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcMemInfoRun(t *testing.T) {
	p := NewProcMemInfo()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/meminfo data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/meminfo data should fail on non-linux.")
		}
	}
}

func TestNewProcMemInfoToJson(t *testing.T) {
	p := NewProcMemInfo()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/meminfo data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"Active", "AnonHugePages", "AnonPages", "Bounce", "Buffers", "Cached", "CommitLimit",
			"DirectMap2M", "DirectMap4k", "Dirty", "Inactive", "MemFree", "MemTotal", "KernelStack", "PageTables"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
