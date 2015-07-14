package executors

import (
	"encoding/json"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func newConfigExecutorForTest(t *testing.T) resourced_config.Config {
	config := resourced_config.Config{}
	config.GoStruct = "Shell"
	config.Kind = "executor"
	config.Path = "/x/uptime"
	config.GoStructFields = make(map[string]interface{})
	config.GoStructFields["Command"] = "uptime"

	return config
}

func TestDynamicConstructor(t *testing.T) {
	ResetConditionsMetByPath()

	config := newConfigExecutorForTest(t)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	// Test simple run, see if it works
	executor.Run()

	inJson, err := executor.ToJson()
	if err != nil {
		t.Errorf("Failed to serialize data to JSON. Error: %v", err)
	}

	var data map[string]interface{}
	json.Unmarshal(inJson, &data)

	if data["Output"] == nil {
		t.Fatalf("There should always be output from uptime.")
	}
	if data["Output"].(string) == "" {
		t.Fatalf("There should always be output from uptime. Output: %v", data["Output"].(string))
	}
}

func TestIsConditionMetDefaultQuery(t *testing.T) {
	ResetConditionsMetByPath()

	config := newConfigExecutorForTest(t)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	if executor.IsConditionMet() == false {
		t.Fatalf("The default query is [true], so IsConditionMet should always be true. Value: %v", executor.IsConditionMet())
	}
}

func TestIsConditionsMetCustomQuery(t *testing.T) {
	ResetConditionsMetByPath()

	config := newConfigExecutorForTest(t)
	config.GoStructFields["Conditions"] = `["<", {"/r/load-avg": "LoadAvg1m"}, 100]`

	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	executor.SetReadersDataInBytes(data)

	if executor.IsConditionMet() == false {
		t.Fatalf("Based on the data, conditions should have been met. Value: %v", executor.IsConditionMet())
	}
}

func TestRunAndCheckConditionsMet(t *testing.T) {
	ResetConditionsMetByPath()

	config := newConfigExecutorForTest(t)
	config.GoStructFields["Conditions"] = `["<", {"/r/load-avg": "LoadAvg1m"}, 100]`
	config.GoStructFields["LowThreshold"] = 1
	config.GoStructFields["HighThreshold"] = 2

	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	executor, err := NewGoStructByConfig(config)
	if executor == nil {
		t.Fatalf("Shell constructor did not do its job")
	}
	if err != nil {
		t.Errorf("Shell constructor did not do its job. Error: %v", err)
	}

	executor.SetReadersDataInBytes(data)

	err = executor.Run()
	if err != nil {
		t.Fatalf("Running uptime should always work. Error: %v", err)
	}

	if executor.LowThresholdExceeded() == true {
		t.Errorf("Ran only 1 time. LowThreshold should not have been exceeded yet. ConditionMetByPathCounter: %v", ConditionMetByPathCounter["/x/uptime"])
	}
	executor.Run()
	if executor.HighThresholdExceeded() == true {
		t.Errorf("Ran 2 times. LowThreshold should have been exceeded. ConditionMetByPathCounter: %v", ConditionMetByPathCounter["/x/uptime"])
	}
}
