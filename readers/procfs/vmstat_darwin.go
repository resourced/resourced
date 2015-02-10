// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcVmStat is ProcVmStat constructor.
func NewProcVmStat() *ProcVmStat {
	p := &ProcVmStat{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcVmStat is a reader that scrapes /proc/vmstat data.
type ProcVmStat struct {
	Data map[string][]interface{}
}

func (p *ProcVmStat) Run() error {
	return errors.New("/proc/vmstat is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcVmStat) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
