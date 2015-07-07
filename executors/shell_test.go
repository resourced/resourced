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
	s.LowTreshold = 1
	s.HighTreshold = 3

	if s.LowTresholdExceeded() {
		t.Errorf("In the beginning, low treshold should not be exceeded. Value: %v", s.LowTresholdExceeded())
	}

	s.ConditionMet()
	s.ConditionMet()

	if !s.LowTresholdExceeded() {
		t.Errorf("Low treshold should be exceeded. Value: %v", s.LowTresholdExceeded())
	}

	s.ConditionMet()
	s.ConditionMet()

	if !s.HighTresholdExceeded() {
		t.Errorf("High treshold should be exceeded. Value: %v", s.LowTresholdExceeded())
	}
}
