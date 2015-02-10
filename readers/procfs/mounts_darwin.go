// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcMounts is ProcMounts constructor.
func NewProcMounts() *ProcMounts {
	p := &ProcMounts{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcMounts is a reader that scrapes /proc/mounts data.
type ProcMounts struct {
	Data map[string][]interface{}
}

func (p *ProcMounts) Run() error {
	return errors.New("/proc/mounts is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcMounts) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
