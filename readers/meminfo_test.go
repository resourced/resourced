package readers

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewMeminfo(t *testing.T) {
	n := NewMeminfo()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewMeminfoRun(t *testing.T) {
	n := NewMeminfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewMeminfoToJson(t *testing.T) {
	n := NewMeminfo()
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
		if strings.Contains(string(jsonData), "Error") {
			t.Errorf("jsonData shouldn't return error on linux: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Active`) {
			t.Errorf("jsonData does not contain 'Active' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `AnonHugePages`) {
			t.Errorf("jsonData does not contain 'AnonHugePages' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `AnonPages`) {
			t.Errorf("jsonData does not contain 'AnonPages' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Bounce`) {
			t.Errorf("jsonData does not contain 'Bounce' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Buffers`) {
			t.Errorf("jsonData does not contain 'Buffers' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Cached`) {
			t.Errorf("jsonData does not contain 'Cached' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `CommitLimit`) {
			t.Errorf("jsonData does not contain 'CommitLimit' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `DirectMap2M`) {
			t.Errorf("jsonData does not contain 'DirectMap2M' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `DirectMap4k`) {
			t.Errorf("jsonData does not contain 'DirectMap4k' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Dirty`) {
			t.Errorf("jsonData does not contain 'Dirty' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Inactive`) {
			t.Errorf("jsonData does not contain 'Inactive' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `MemAvailable`) {
			t.Errorf("jsonData does not contain 'MemAvailable' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `MemFree`) {
			t.Errorf("jsonData does not contain 'MemFree' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `MemTotal`) {
			t.Errorf("jsonData does not contain 'MemTotal' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `KernelStack`) {
			t.Errorf("jsonData does not contain 'KernelStack' key. jsonData: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `PageTables`) {
			t.Errorf("jsonData does not contain 'PageTables' key. jsonData: %s", jsonData)
		}
	}
}
