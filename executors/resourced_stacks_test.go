package executors

import (
	"encoding/json"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func newResourcedStacksExecutorForTest(t *testing.T) resourced_config.Config {
	config := resourced_config.Config{}
	config.GoStruct = "ResourcedStacks"
	config.Kind = "executor"
	config.Path = "/x/resourced-stacks"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Name"] = "helloworld"
	config.GoStructFields["Root"] = "/Users/didip/tmp/resourced-stacks-testrepo"
	config.GoStructFields["GitRepo"] = "https://github.com/resourced/resourced-stacks-testrepo.git"
	config.GoStructFields["DryRun"] = true

	return config
}

func TestResourcedStacksRun(t *testing.T) {
	ResetConditionsMetByPath()

	config := newResourcedStacksExecutorForTest(t)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("ResourcedStacks constructor did not do its job")
	}
	if err != nil {
		t.Errorf("ResourcedStacks constructor did not do its job. Error: %v", err)
	}

	executor.Run()

	dataJson, err := executor.ToJson()
	if err != nil {
		t.Fatalf("Failed to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(dataJson, &data)

	if data["Error"] != nil {
		t.Fatalf("Pulling and running from ResourceD Stacks should work. Error: %v", err)
	}

	if data["Output"] == nil {
		t.Fatalf("There should always be output from ResourceD Stacks.")
	}
	if data["Output"].(string) == "" {
		t.Fatalf("There should always be output from ResourceD Stacks. Output: %v", data["Output"].(string))
	}
}
