package readers

import (
	"strings"
	"testing"
)

func TestNewHostInfo(t *testing.T) {
	n := NewHostInfo()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewHostInfoRun(t *testing.T) {
	n := NewHostInfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewHostInfoToJson(t *testing.T) {
	n := NewHostInfo()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling memory data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"Hostname", "Uptime", "Procs", "OS", "Platform", "PlatformFamily", "PlatformVersion", "VirtualizationSystem", "VirtualizationRole"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}

func TestNewHostUsers(t *testing.T) {
	n := NewHostUsers()
	if n.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestNewHostUsersRun(t *testing.T) {
	n := NewHostUsers()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}
}

func TestNewHostUsersToJson(t *testing.T) {
	n := NewHostUsers()
	err := n.Run()
	if err != nil {
		t.Errorf("Parsing memory data should always be successful. Error: %v", err)
	}

	jsonData, err := n.ToJson()
	if err != nil {
		t.Errorf("Marshalling memory data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}

	keysToTest := []string{"user", "terminal", "host", "started"}

	for _, key := range keysToTest {
		if !strings.Contains(jsonDataString, key) {
			t.Errorf("jsonDataString does not contain '%v' key. jsonDataString: %v", key, jsonDataString)
		}
	}
}
