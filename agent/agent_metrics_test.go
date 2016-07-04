package agent

import (
	"encoding/json"
	"testing"
)

func TestBuildResultDBPayloadFromGraphiteMetric(t *testing.T) {
	key := "test.simple"
	value := 3.4

	agent := createAgentForTest(t)

	payload := agent.buildResultDBPayloadFromKeyValueMetric(key, value)

	if payload["Data"] == nil {
		t.Errorf("Failed to parse Graphite metric to ResourceD payload")
	}
	if payload["Data"].(map[string]interface{})["simple"] != 3.4 {
		t.Errorf("Failed to parse Graphite metric to ResourceD payload. Expected: %v, Result: %v", value, payload["Data"].(map[string]interface{})["simple"])
	}
}

func TestBuildResultDBPayloadFromGraphiteMetricWithHostname(t *testing.T) {
	key := "test.localhost.simple"
	value := 3.4

	agent := createAgentForTest(t)

	payload := agent.buildResultDBPayloadFromKeyValueMetric(key, value)

	if payload["Data"] == nil {
		t.Errorf("Failed to parse Graphite metric to ResourceD payload")
	}
	if payload["Data"].(map[string]interface{})["localhost"].(map[string]interface{})["simple"] != 3.4 {
		t.Errorf("Failed to parse Graphite metric to ResourceD payload. Expected: %v, Result: %v", value, payload["Data"].(map[string]interface{})["localhost"].(map[string]interface{})["simple"])
	}
}
