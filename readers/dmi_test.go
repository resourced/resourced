package readers

import (
	"runtime"
	"strings"
	"testing"
)

func TestNewDMI(t *testing.T) {
	d := NewDMI()
	if d.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestDMIRun(t *testing.T) {
	// 'dmidecode' is not available on Darwin, run test only on Linux
	if runtime.GOOS == "linux" {
		d := NewDMI()
		err := d.Run()
		if err != nil {
			t.Errorf("Fetching DMI data should always work. Error: %v", err)
		}
	}
}

func TestDMIToJson(t *testing.T) {
	// 'dmidecode' is not available on Darwin, run test only on Linux
	if runtime.GOOS == "linux" {
		d := NewDMI()
		err := d.Run()
		if err != nil {
			t.Errorf("Fetching DMI data should always work. Error: %v", err)
		}

		jsonData, err := d.ToJson()
		if err != nil {
			t.Errorf("Marshalling df data should always be successful. Error: %v", err)
		}

		jsonDataString := string(jsonData)

		keysToTest := []string{"0x0001", "0x0002", "DMIName", "UUID"}

		for _, key := range keysToTest {
			if !strings.Contains(jsonDataString, key) {
				t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
			}
		}
	}
}
