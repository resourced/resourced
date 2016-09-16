package executors

import (
	"encoding/json"

	"github.com/resourced/resourced/libprocess"
)

func init() {
	Register("Shell", NewShell)
}

func NewShell() IExecutor {
	s := &Shell{}
	s.Data = make(map[string]interface{})

	return s
}

type Shell struct {
	Base
	Data map[string]interface{}
}

// Run shells out external program and store the output on c.Data.
func (s *Shell) Run() error {
	s.Data["Conditions"] = s.Conditions

	if s.IsConditionMet() && s.LowThresholdExceeded() && !s.HighThresholdExceeded() {
		output, err := libprocess.NewCmd(s.Command).CombinedOutput()
		s.Data["Output"] = string(output)

		if err != nil {
			s.Data["Error"] = err.Error()
			s.Data["ExitStatus"] = 1
		} else {
			s.Data["ExitStatus"] = 0
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
// If there are no meaningful results, ToJson returns nil.
func (s *Shell) ToJson() ([]byte, error) {
	output, outputFound := s.Data["Output"]
	errorString, errorFound := s.Data["Error"]

	if !outputFound && !errorFound {
		return []byte("{}"), nil
	}

	if output.(string) == "" && errorString.(string) == "" {
		return []byte("{}"), nil
	}

	return json.Marshal(s.Data)
}
