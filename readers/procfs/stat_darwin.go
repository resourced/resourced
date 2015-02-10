// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcStat is ProcStat constructor.
func NewProcStat() *ProcStat {
	p := &ProcStat{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcStat is a reader that scrapes /proc/stat data.
type ProcStat struct {
	Data map[string][]interface{}
}

func (p *ProcStat) Run() error {
	return errors.New("/proc/stat is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcStat) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
