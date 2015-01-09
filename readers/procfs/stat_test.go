package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcStat(t *testing.T) {
	p := NewProcStat()
	if p.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewProcStatRun(t *testing.T) {
	p := NewProcStat()

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

func TestNewProcStatToJson(t *testing.T) {
	p := NewProcStat()
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

		keysToTest := []string{"cpu_all", "user", "nice", "system", "idle", "iowait"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
