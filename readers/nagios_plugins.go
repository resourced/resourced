package readers

import (
	"encoding/json"
	"strings"

	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
)

func init() {
	Register("NagiosPlugin", NewNagiosPlugin)
}

func NewNagiosPlugin() IReader {
	s := &NagiosPlugin{}
	s.Data = make(map[string]interface{})

	return s
}

type NagiosPlugin struct {
	Shell
}

// Run NagiosPlugins out external program and store the output on c.Data.
func (s *NagiosPlugin) Run() error {
	if s.Command != "" {
		s.Command = libstring.ExpandTildeAndEnv(s.Command)

		nagiosPluginOutputBytes, err := libprocess.NewCmd(s.Command).CombinedOutput()
		if err != nil {
			s.Data["ExitStatus"] = 1
		} else {
			s.Data["ExitStatus"] = 0
		}

		nagiosPluginOutput := strings.TrimSpace(string(nagiosPluginOutputBytes))

		if strings.Contains(nagiosPluginOutput, "OK") {
			s.Data["ExitStatus"] = 0
			s.Data["Message"] = nagiosPluginOutput
		} else if strings.Contains(nagiosPluginOutput, "WARNING") {
			s.Data["ExitStatus"] = 1
			s.Data["Message"] = nagiosPluginOutput
		} else if strings.Contains(nagiosPluginOutput, "CRITICAL") {
			s.Data["ExitStatus"] = 2
			s.Data["Message"] = nagiosPluginOutput
		} else if strings.Contains(nagiosPluginOutput, "UNKNOWN") {
			s.Data["ExitStatus"] = 3
			s.Data["Message"] = nagiosPluginOutput
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (s *NagiosPlugin) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
