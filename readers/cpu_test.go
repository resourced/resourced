package readers

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewCpuInfoRun(t *testing.T) {
	n := NewCpuInfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Reading cpu data should always be successful. Error: %v", err)
	}
}

func TestNewCpuInfoToJson(t *testing.T) {
	n := NewCpuInfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Reading cpu data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling cpu data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"cpu", "vendorId", "family", "model", "stepping", "cores", "modelName", "cacheSize", "flags"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}

// ------------------------------------------------------------

func TestNewCpuStat(t *testing.T) {
	n := NewCpuStat()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewCpuStatRun(t *testing.T) {
	// Darwin version of gopsutil_cpu.CPUTimes is broken, so we are only testing this on Linux.
	if runtime.GOOS == "linux" {

		n := NewCpuStat()
		err := n.Run()
		if err != nil {
			t.Errorf("Reading cpu stat data should always be successful. Error: %v", err)
		}
	}
}

func TestNewCpuStatToJson(t *testing.T) {
	// Darwin version of gopsutil_cpu.CPUTimes is broken, so we are only testing this on Linux.
	if runtime.GOOS == "linux" {
		n := NewCpuStat()
		err := n.Run()
		if err != nil {
			t.Fatalf("Reading cpu stat data should always be successful. Error: %v", err)
		}

		jsonData, err := n.ToJson()
		if err != nil {
			t.Errorf("Marshalling cpu stat data should always be successful. Error: %v", err)
		}

		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"cpu", "user", "system", "idle", "nice", "iowait", "irq", "softirq", "steal", "guest", "guestNice", "stolen"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
