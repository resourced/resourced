package writers

import (
	"bytes"
	"encoding/json"

	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
)

func init() {
	Register("Shell", NewShell)
}

func NewShell() IWriter {
	s := &Shell{}
	s.Data = make(map[string]interface{})

	return s
}

type Shell struct {
	Base
	Command string
	Data    map[string]interface{}
}

// Run shells out external program and store the output on c.Data.
func (s *Shell) Run() error {
	if s.Command != "" {
		s.Command = libstring.ExpandTildeAndEnv(s.Command)

		readersDataJsonBytes, err := json.Marshal(s.GetReadersData())
		if err != nil {
			return err
		}

		cmd := libprocess.NewCmd(s.Command)
		cmd.Stdin = bytes.NewReader(readersDataJsonBytes)

		outputJson, err := cmd.CombinedOutput()

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
