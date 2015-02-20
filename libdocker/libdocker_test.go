package libdocker

import (
	"testing"
)

func TestAllContainers(t *testing.T) {
	_, err := AllContainers()
	if err != nil {
		t.Errorf("Gettting Docker containers info should not fail. Error: %v", err)
	}
}

func TestAllInspectedContainers(t *testing.T) {
	_, err := AllInspectedContainers()
	if err != nil {
		t.Errorf("Gettting Docker containers info should not fail. Error: %v", err)
	}
}

func TestAllImages(t *testing.T) {
	_, err := AllImages()
	if err != nil {
		t.Errorf("Gettting Docker images info should not fail. Error: %v", err)
	}
}

func TestAllInspectedImages(t *testing.T) {
	_, err := AllInspectedImages()
	if err != nil {
		t.Errorf("Gettting Docker images info should not fail. Error: %v", err)
	}
}
