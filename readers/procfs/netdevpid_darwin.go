// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcNetDevPid is ProcNetDevPid constructor.
func NewProcNetDevPid() *ProcNetDevPid {
	p := &ProcNetDevPid{}
	return p
}

// ProcNetDevPid is a reader that scrapes /proc/$pid/net/dev data.
type ProcNetDevPid struct {
	Data map[string]interface{}
}

func (p *ProcNetDevPid) Run() error {
	return errors.New("/proc/net/dev is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcNetDevPid) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
