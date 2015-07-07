package executors

import (
	"testing"
)

func TestSimpleRun(t *testing.T) {
	s := &Shell{}
	s.Data = make(map[string]interface{})
	s.Command = "uptime"

	s.Run()
	if s.Data["Output"] == nil {
		t.Fatalf("There should always be output from uptime.")
	}
	if s.Data["Output"].(string) == "" {
		t.Fatalf("There should always be output from uptime. Output: %v", s.Data["Output"].(string))
	}
}

func TestCrossingTreshold(t *testing.T) {
	s := &Shell{}
	s.Data = make(map[string]interface{})
	s.Command = "uptime"
	s.LowThreshold = 1
	s.HighThreshold = 3

	if s.LowThresholdExceeded() {
		t.Errorf("In the beginning, low treshold should not be exceeded. Value: %v", s.LowThresholdExceeded())
	}

	s.ConditionMet()
	s.ConditionMet()

	if !s.LowThresholdExceeded() {
		t.Errorf("Low treshold should be exceeded. Value: %v", s.LowThresholdExceeded())
	}

	s.ConditionMet()
	s.ConditionMet()

	if !s.HighThresholdExceeded() {
		t.Errorf("High treshold should be exceeded. Value: %v", s.LowThresholdExceeded())
	}
}
