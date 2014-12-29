package config

import (
	"testing"
)

func TestConstructor(t *testing.T) {
	config, err := NewConfigStorage("~/go/src/github.com/resourced/resourced/tests/data/config-reader", "~/go/src/github.com/resourced/resourced/tests/data/config-writer")

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	if len(config.Readers) != 1 {
		t.Errorf("Length of reader config should == 1. config.Readers: %v", config.Readers)
	}
	if len(config.Writers) != 1 {
		t.Errorf("Length of reader config should == 1. config.Writers: %v", config.Writers)
	}
}

func TestRun(t *testing.T) {
	config, err := NewConfigStorage("~/go/src/github.com/resourced/resourced/tests/data/config-reader", "~/go/src/github.com/resourced/resourced/tests/data/config-writer")

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	_, err = config.Readers[0].Run()
	if err != nil {
		t.Fatalf("Reader command should work. Error: %v", err)
	}
}
