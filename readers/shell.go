package readers

import (
	"encoding/json"

	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
)

func init() {
	Register("Shell", NewShell)
}

func NewShell() IReader {
	s := &Shell{}
	s.Data = make(map[string]interface{})

	return s
}

type Shell struct {
	Command string
	Data    map[string]interface{}
}

// Run shells out external program and store the output on c.Data.
func (s *Shell) Run() error {
	if s.Command != "" {
		s.Command = libstring.ExpandTildeAndEnv(s.Command)

		outputJson, err := libprocess.NewCmd(s.Command).CombinedOutput()

		var output map[string]interface{}
		json.Unmarshal(outputJson, &output)

		s.Data["Output"] = output

		if err != nil {
			s.Data["ExitStatus"] = 1
		} else {
			s.Data["ExitStatus"] = 0
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (s *Shell) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
