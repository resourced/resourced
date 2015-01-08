package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcVmStatRun(t *testing.T) {
	p := NewProcVmStat()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/cpuinfo data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/cpuinfo data should fail on non-linux.")
		}
	}
}

func TestNewProcVmStatToJson(t *testing.T) {
	p := NewProcVmStat()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/cpuinfo data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"nr_free_pages", "nr_alloc_batch", "nr_inactive_anon", "nr_active_anon", "nr_inactive_file", "nr_active_file"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
