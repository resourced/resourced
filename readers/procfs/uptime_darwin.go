// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcUptime", NewProcUptime)
}

// NewProcUptime is ProcUptime constructor.
func NewProcUptime() readers.IReader {
	p := &ProcUptime{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcUptime is a reader that scrapes /proc/uptime data.
type ProcUptime struct {
	Data map[string][]interface{}
}

func (p *ProcUptime) Run() error {
	return errors.New("/proc/uptime is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcUptime) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
