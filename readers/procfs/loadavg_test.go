package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcLoadAvgRun(t *testing.T) {
	p := NewProcLoadAvg()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/loadavg data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/loadavg data should fail on non-linux.")
		}
	}
}

func TestNewProcLoadAvgToJson(t *testing.T) {
	p := NewProcLoadAvg()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/loadavg data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"last1min", "last5min", "last15min"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
