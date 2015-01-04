package config

import (
	"os"
	"testing"
)

func TestConstructor(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	config, err := NewConfigStorage(gopath+"/src/github.com/resourced/resourced/tests/data/config-reader", gopath+"/src/github.com/resourced/resourced/tests/data/config-writer")

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	if len(config.Readers) <= 0 {
		t.Errorf("Length of reader config should > 0. config.Readers: %v", config.Readers)
	}
	if len(config.Writers) != 1 {
		t.Errorf("Length of reader config should == 1. config.Writers: %v", config.Writers)
	}
}
