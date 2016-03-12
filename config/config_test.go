package config

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewConfigs(t *testing.T) {
	config, err := NewConfigs("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs")
	if err != nil {
		t.Fatalf("Initializing Configs should work. Error: %v", err)
	}

	if len(config.Readers) <= 0 {
		t.Errorf("Length of reader config should > 0. config.Readers: %v", config.Readers)
	}
	if len(config.Writers) <= 0 {
		t.Errorf("Length of reader config should > 0. len(config.Writers): %v", len(config.Writers))
	}
}

func TestNewReaderConfig(t *testing.T) {
	config, err := NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs/readers/docker-containers-memory.toml"), "reader")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	if config.GoStruct != "DockerContainersMemory" {
		t.Fatalf("Config is initialized incorrectly. config.GoStruct: %v", config.GoStruct)
	}
	if config.Path != "/docker/containers/memory" {
		t.Fatalf("Config is initialized incorrectly. config.Path: %v", config.Path)
	}
	if config.Interval != "3s" {
		t.Fatalf("Config is initialized incorrectly. config.Interval: %v", config.Interval)
	}
	if config.GoStructFields["CgroupBasePath"] != "/sys/fs/cgroup/memory/docker" {
		inJson, _ := json.Marshal(config.GoStructFields)
		t.Fatalf("Config is initialized incorrectly. config.GoStructFields: %v", string(inJson))
	}
}

func TestNewWriterConfigWithJsonProcessor(t *testing.T) {
	config, err := NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs/writers/stdout.toml"), "writer")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}
	if config.GoStructFields["JsonProcessor"] == "" {
		inJson, _ := json.Marshal(config.GoStructFields)
		t.Fatalf("Config is initialized incorrectly. config.GoStructFields: %v", string(inJson))
	}
}

func TestCommonJsonData(t *testing.T) {
	config, err := NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs/readers/docker-containers-memory.toml"), "reader")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	record := config.CommonJsonData()
	if len(record) == 0 {
		t.Error("common data should never be empty")
	}
	for _, key := range []string{"UnixNano", "Path", "Interval"} {
		if _, ok := record[key]; !ok {
			t.Errorf("%v data should never be empty.", key)
		}
	}
}

func TestPathWithPrefix(t *testing.T) {
	config, err := NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs/readers/docker-containers-memory.toml"), "reader")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	path := config.PathWithPrefix()
	if !strings.HasPrefix(path, "/r") {
		t.Errorf("Path should have been prefixed with /r. Path: %v", path)
	}
	if strings.HasPrefix(path, "/w") {
		t.Errorf("Path is prefixed incorrectly. Path: %v", path)
	}
}

func TestpathWithKindPrefix(t *testing.T) {
	config, err := NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs/readers/docker-containers-memory.toml"), "reader")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	toBeTested := config.PathWithKindPrefix("r", "/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = config.PathWithKindPrefix("r", "/r/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = config.PathWithKindPrefix("w", "/stuff")
	if toBeTested != "/w/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = config.PathWithKindPrefix("w", "/w/stuff")
	if toBeTested != "/w/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = config.PathWithKindPrefix("x", "/stuff")
	if toBeTested != "/w/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = config.PathWithKindPrefix("x", "/x/stuff")
	if toBeTested != "/x/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}
}
