package executors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/resourced/resourced/libmap"
	"github.com/resourced/resourced/libstring"
)

func TestDiskCleanerRun(t *testing.T) {
	config := newConfigExecutorForTest(t)
	config.GoStruct = "DiskCleaner"
	config.Kind = "executor"
	config.Path = "/x/disk-cleaner/home"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Globs"] = make([]interface{}, 0)
	config.GoStructFields["Globs"] = append(config.GoStructFields["Globs"].([]interface{}), "~/*.log")

	// Create a file to be deleted
	seed, _ := libstring.GeneratePassword(32)

	err := ioutil.WriteFile(libstring.ExpandTildeAndEnv(fmt.Sprintf("~/testing-%v.log", seed)), []byte(""), 0644)
	if err != nil {
		t.Fatalf("Creating a file should always work. Error: %v", err)
	}

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	counterDB := libmap.NewTSafeMapCounter()
	executor.SetCounterDB(counterDB)

	executor.Run()

	dataJson, err := executor.ToJson()
	if err != nil {
		t.Fatalf("Failed to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(dataJson, &data)

	if data["Success"] == nil {
		t.Fatalf("There should always be output from uptime.")
	}
	if len(data["Success"].([]interface{})) == 0 {
		t.Fatalf("There should be success cleaned file data. Success: %v", data["Success"].([]interface{}))
	}
}
