package agent

import (
	"os"
	"testing"
)

func TestConstructor(t *testing.T) {
	agent := createAgentForTest(t)

	if len(agent.Tags) == 0 {
		t.Error("agent.Tags should not be empty.")
	}

	for key, value := range agent.Tags {
		if key == "redis" && value != "3.0.1" {
			t.Errorf("agent.Tags is incorrect: %v", agent.Tags)
		}
		if key == "mysql" && value != "5.6.24" {
			t.Errorf("agent.Tags is incorrect: %v", agent.Tags)
		}
	}
}
