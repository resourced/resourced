package procfs

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewProcNetDev(t *testing.T) {
	p := NewProcNetDev()
	if p.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewProcNetDevRun(t *testing.T) {
	p := NewProcNetDev()

	if runtime.GOOS == "linux" {
		err := p.Run()
		if err != nil {
			t.Errorf("Reading /proc/net/dev data should work on linux. Error: %v", err)
		}
	} else {
		err := p.Run()
		if err == nil {
			t.Error("Reading /proc/net/dev data should fail on non-linux.")
		}
	}
}

func TestNewProcNetDevToJson(t *testing.T) {
	p := NewProcNetDev()
	p.Run()

	jsonData, err := p.ToJson()
	if err != nil {
		t.Errorf("Marshalling /proc/net/dev data should always work. Error: %v", err)
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
