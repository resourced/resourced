package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcDiskStatsRun(t *testing.T) {
	p := NewProcDiskStats()

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

func TestNewProcDiskStatsToJson(t *testing.T) {
	p := NewProcDiskStats()
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

		keysToTest := []string{"major", "minor", "name", "read_ios", "read_merges", "read_sectors", "read_ticks",
			"write_ios", "write_merges", "write_sectors", "write_ticks", "in_flight", "io_ticks", "time_in_queue"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
