package readers

import (
	"encoding/json"
	gopsutil_cpu "github.com/shirou/gopsutil/cpu"
)

func NewCpuInfo() *CpuInfo {
	c := &CpuInfo{}
	c.Data = make(map[string][]gopsutil_cpu.CPUInfoStat)
	return c
}

type CpuInfo struct {
	Data map[string][]gopsutil_cpu.CPUInfoStat
}

func (c *CpuInfo) Run() error {
	data, err := gopsutil_cpu.CPUInfo()
	if err != nil {
		return err
	}

	c.Data["CpuInfo"] = data
	return nil
}

func (c *CpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}

// ----------------------------------------------------------------

func NewCpuStat() *CpuStat {
	c := &CpuStat{}
	c.Data = make(map[string][]gopsutil_cpu.CPUTimesStat)
	return c
}

type CpuStat struct {
	Data map[string][]gopsutil_cpu.CPUTimesStat
}

func (c *CpuStat) Run() error {
	data, err := gopsutil_cpu.CPUTimes(true)
	if err != nil {
		return err
	}

	c.Data["CpuStat"] = data
	return nil
}

func (c *CpuStat) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
