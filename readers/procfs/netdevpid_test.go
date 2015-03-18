package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcNetDevPid(t *testing.T) {
	p := NewProcNetDevPid()
	if p.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewProcNetDevPidRun(t *testing.T) {
	p := NewProcNetDevPid()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/$pid/net/dev data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/$pid/net/dev data should fail on non-linux.")
		}
	}
}

func TestNewProcNetDevPidToJson(t *testing.T) {
	p := NewProcNetDevPid()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/$pid/net/dev data should always work. Error: %v", err)
	}

	if runtime.GOOS == "linux" {
		jsonDataString := string(jsonData)

		if strings.Contains(jsonDataString, "Error") {
			t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
		}

		keysToTest := []string{"iface", "rxbytes", "rxpackets", "rxerrs", "rxdrop", "rxfifo", "rxframe",
			"rxcompressed", "rxmulticast", "txbytes", "txpackets", "txerrs", "txdrop", "txfifo", "txcolls", "txcarrier", "txcompressed"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
