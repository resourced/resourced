// Package libtime provides time related library functions.
package libtime

import (
	"time"
)

// SleepString is a convenience function that performs `time.Sleep` given string duration.
func SleepString(definition string) error {
	delayTime, err := time.ParseDuration(definition)
	if err != nil {
		return err
	}

	time.Sleep(delayTime)
	return nil
}

func ParseDurationWithDefault(definition, defaultDef string) time.Duration {
	dur, err := time.ParseDuration(definition)
	if err != nil {
		dur, _ = time.ParseDuration(defaultDef) // Default
	}

	return dur
}
