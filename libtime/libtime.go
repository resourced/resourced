package libtime

import (
	"time"
)

func SleepString(definition string) error {
	delayTime, err := time.ParseDuration(definition)
	if err != nil {
		return err
	}

	time.Sleep(delayTime)
	return nil
}

func ParseIsoString(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999Z", dateString)
}
