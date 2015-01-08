// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcCpuInfo is ProcCpuInfo constructor.
func NewProcCpuInfo() *ProcCpuInfo {
	c := &ProcCpuInfo{}
	c.Data = make(map[string][]linuxproc.Processor)
	return c
}

// ProcCpuInfo is a reader that scrapes /proc/cpuinfo data.
// Data source: https://github.com/shirou/gopsutil/tree/master/cpu
type ProcCpuInfo struct {
	Data map[string][]linuxproc.Processor
}

func (c *ProcCpuInfo) Run() error {
	cpuinfo, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		return err
	}

	c.Data["ProcCpuInfo"] = cpuinfo.Processors
	return nil
}

func (c *ProcCpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
