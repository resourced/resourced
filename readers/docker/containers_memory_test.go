// +build docker
package docker

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewDockerContainersMemoryRun(t *testing.T) {
	n := NewDockerContainersMemory()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewDockerContainersMemoryToJson(t *testing.T) {
	n := NewDockerContainersMemory()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling memory data should always be successful. Error: %v", err)
	}

	if runtime.GOOS == "darwin" {
		if !strings.Contains(string(jsonData), "Error") {
			t.Errorf("jsonData should return error on darwin: %s", jsonData)
		}
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"container_id", "cache", "rss", "rss_huge", "mapped_file", "pgpgin", "pgpgout",
			"pgfault", "pgmajfault", "inactive_anon", "active_anon", "inactive_file", "active_file", "unevictable", "hierarchical_memory_limit",
			"total_cache", "total_rss", "total_rss_huge", "total_mapped_file", "total_pgpgin", "total_pgpgout", "total_pgfault", "total_pgmajfault",
			"total_inactive_anon", "total_active_anon", "total_inactive_file", "total_active_file", "total_unevictable"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
