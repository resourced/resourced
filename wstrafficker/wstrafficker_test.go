package wstrafficker

import (
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, _, err := NewClient("http://localhost:55555", "http://localhost:55655/api/ws", nil)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("Creating new client with basic settings should not fail. Error: %v", err)
	}
}

func TestBasic(t *testing.T) {
	_, err := NewWSTrafficker("http://localhost:55555", "http://localhost:55655/api/ws", nil)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("Creating new trafficker with basic settings should not fail. Error: %v", err)
	}
}
