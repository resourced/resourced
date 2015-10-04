// +build docker
package docker

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewDockerContainersCpuRun(t *testing.T) {
	n := NewDockerContainersCpu()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewDockerContainersCpuToJson(t *testing.T) {
	if runtime.GOOS == "linux" {
		n := NewDockerContainersCpu()
		err := n.Run()
		if err != nil {
			t.Errorf("Parsing memory data should always be successful. Error: %v", err)
		}

		jsonData, err := n.ToJson()
		if err != nil {
			t.Errorf("Marshalling memory data should always be successful. Error: %v", err)
		}

		if runtime.GOOS == "linux" {
			jsonDataString := string(jsonData)

			if strings.Contains(jsonDataString, "Error") {
				t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
			}

			keysToTest := []string{"user", "system", "idle", "nice", "iowait", "irq", "softirq",
				"steal", "guest", "guest_nice", "stolen"}

			for _, key := range keysToTest {
				if !strings.Contains(jsonDataString, key) {
					t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
				}
			}
		}
	}
}
