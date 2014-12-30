package libtime

import (
	"testing"
)

func TestSleepString(t *testing.T) {
	if SleepString("1ms") != nil {
		t.Errorf("Failed to sleep")
	}
}
