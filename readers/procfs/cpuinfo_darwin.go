// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcCpuInfo is ProcCpuInfo constructor.
func NewProcCpuInfo() *ProcCpuInfo {
	p := &ProcCpuInfo{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcCpuInfo is a reader that scrapes /proc/cpuinfo data.
type ProcCpuInfo struct {
	Data map[string][]interface{}
}

func (p *ProcCpuInfo) Run() error {
	return errors.New("/proc/cpuinfo is only available on Linux.")
}

func (p *ProcCpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
