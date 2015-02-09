package host

import (
	"testing"
)

func TestNewHost(t *testing.T) {
	h := NewHost("localhost")
	if h == nil {
		t.Error("Creating host should always be successful.")
	}
}

func TestNewHostByHostname(t *testing.T) {
	_, err := NewHostByHostname()
	if err != nil {
		t.Errorf("On common systems, creating Host should always be successful. Error: %v", err)
	}
}
