package executors

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/resourced/resourced/libstring"
)

func init() {
	Register("DiskCleaner", NewDiskCleaner)
}

func NewDiskCleaner() IExecutor {
	dc := &DiskCleaner{}
	dc.Data = make(map[string]interface{})

	return dc
}

type DiskCleaner struct {
	Base
	Data  map[string]interface{}
	Globs []interface{}
}

// Run shells out external program and store the output on c.Data.
func (dc *DiskCleaner) Run() error {
	dc.Data["Conditions"] = dc.Conditions

	if dc.IsConditionMet() && dc.LowThresholdExceeded() && !dc.HighThresholdExceeded() {
		successOutput := make([]string, 0)
		failOutput := make([]string, 0)

		for _, globInterface := range dc.Globs {
			glob := globInterface.(string)
			glob = libstring.ExpandTildeAndEnv(glob)

			matches, err := filepath.Glob(glob)
			if err != nil {
				dc.Data["Error"] = err.Error()
				dc.Data["ExitStatus"] = 1

				sendErr := dc.SendToMaster(dc.Data)
				if sendErr != nil {
					logrus.Error(sendErr)
				}

				return err
			}

			for _, fullpath := range matches {
				err := os.RemoveAll(fullpath)
				if err != nil {
					failOutput = append(failOutput, fullpath)
				} else {
					successOutput = append(successOutput, fullpath)
				}
			}
		}

		if len(failOutput) > 0 {
			dc.Data["ExitStatus"] = 1
		} else {
			dc.Data["ExitStatus"] = 0
		}

		dc.Data["Success"] = successOutput
		dc.Data["Failure"] = failOutput

		sendErr := dc.SendToMaster(dc.Data)
		if sendErr != nil {
			logrus.Error(sendErr)
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
// If there are no meaningful results, ToJson returns nil.
func (dc *DiskCleaner) ToJson() ([]byte, error) {
	successOutputInterface, successFound := dc.Data["Success"]
	failureOutputInterface, failureFound := dc.Data["Failure"]

	if !successFound && !failureFound {
		return nil, nil
	}

	successOutput := successOutputInterface.([]string)
	failureOutput := failureOutputInterface.([]string)

	if len(successOutput) == 0 && len(failureOutput) == 0 {
		return nil, nil
	}

	return json.Marshal(dc.Data)
}
