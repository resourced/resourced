package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcUptimeRun(t *testing.T) {
	p := NewProcUptime()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/uptime data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/uptime data should fail on non-linux.")
		}
	}
}

func TestNewProcUptimeToJson(t *testing.T) {
	p := NewProcUptime()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/uptime data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"total", "idle"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
