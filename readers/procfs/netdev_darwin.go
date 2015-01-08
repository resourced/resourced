// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcNetDev is ProcNetDev constructor.
func NewProcNetDev() *ProcNetDev {
	p := &ProcNetDev{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcNetDev is a reader that scrapes /proc/net/dev data.
type ProcNetDev struct {
	Data map[string][]interface{}
}

func (p *ProcNetDev) Run() error {
	return errors.New("/proc/net/dev is only available on Linux.")
}

func (p *ProcNetDev) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
