// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcCpuInfo", NewProcCpuInfo)
}

// NewProcCpuInfo is ProcCpuInfo constructor.
func NewProcCpuInfo() readers.IReader {
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

// ToJson serialize Data field to JSON.
func (p *ProcCpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
