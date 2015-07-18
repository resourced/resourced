package agent

import (
	resourced_config "github.com/resourced/resourced/config"
	"os"
	"strings"
	"testing"
)

func createConfigForAgentWriterTest(t *testing.T) resourced_config.Config {
	config := resourced_config.Config{}
	config.Path = "/insights/du"
	config.Interval = "3s"
	config.Kind = "writer"
	config.ReaderPaths = []string{"/du"}
	config.GoStruct = "StdOut"
	config.GoStructFields = make(map[string]interface{})
	return config
}

func TestRunGoStructWriterWithJsonFlattener(t *testing.T) {
	agent := createAgentForTest(t)

	for _, readerConfig := range agent.Configs.Readers {
		if readerConfig.Path == "/du" {
			agent.Run(readerConfig)

			config := createConfigForAgentWriterTest(t)
			config.GoStructFields["JsonProcessor"] = os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/data/script-writer/json-flattener.py")

			writerData, err := agent.runGoStructWriter(config)
			if err != nil {
				t.Fatalf("runGoStructWriter should not fail. Error: %v", err)
			}

			writerDataString := string(writerData)

			keysToTest := []string{"/du.Data./.DeviceName"}

			for _, key := range keysToTest {
				if !strings.Contains(writerDataString, key) {
					t.Errorf("writerDataString does not contain '%v' key. writerDataString: %v", key, writerDataString)
				}
			}
		}
	}
}

func TestRunGoStructWriterWithInsightsDuFormatter(t *testing.T) {
	agent := createAgentForTest(t)

	for _, readerConfig := range agent.Configs.Readers {
		if readerConfig.Path == "/du" {
			agent.Run(readerConfig)

			config := createConfigForAgentWriterTest(t)
			config.GoStructFields["JsonProcessor"] = os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/data/script-writer/insights/du-formatter.py")

			writerData, err := agent.runGoStructWriter(config)
			if err != nil {
				t.Fatalf("runGoStructWriter should not fail. Error: %v", err)
			}

			writerDataString := string(writerData)

			keysToTest := []string{"Hostname", "DeviceName", "Free", "InodesFree", "InodesTotal", "InodesUsed", "Path", "Total", "Used", "eventType"}

			for _, key := range keysToTest {
				if !strings.Contains(writerDataString, key) {
					t.Errorf("writerDataString does not contain '%v' key. writerDataString: %v", key, writerDataString)
				}
			}
		}
	}
}
