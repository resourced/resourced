// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcCpuInfo", NewProcCpuInfo)
}

// NewProcCpuInfo is ProcCpuInfo constructor.
func NewProcCpuInfo() readers.IReader {
	p := &ProcCpuInfo{}
	p.Data = make(map[string][]linuxproc.Processor)
	return p
}

// ProcCpuInfo is a reader that scrapes /proc/cpuinfo data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/cpuinfo.go
type ProcCpuInfo struct {
	Data map[string][]linuxproc.Processor
}

func (p *ProcCpuInfo) Run() error {
	cpuinfo, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		return err
	}

	p.Data["CpuInfo"] = cpuinfo.Processors
	return nil
}

// ToJson serialize Data field to JSON.
func (p *ProcCpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
