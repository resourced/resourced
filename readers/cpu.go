package readers

import (
	"encoding/json"

	gopsutil_cpu "github.com/shirou/gopsutil/cpu"
)

func init() {
	Register("CpuInfo", NewCpuInfo)
}

// NewCpuInfo is CpuInfo constructor.
func NewCpuInfo() IReader {
	c := &CpuInfo{}
	c.Data = make(map[string][]gopsutil_cpu.InfoStat)
	return c
}

// CpuInfo is a reader that scrapes cpu info data.
// Data source: https://github.com/shirou/gopsutil/tree/master/cpu
type CpuInfo struct {
	Data map[string][]gopsutil_cpu.InfoStat
}

func (c *CpuInfo) Run() error {
	data, err := gopsutil_cpu.Info()
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

// NewCpuStat is CpuStat constructor.
func NewCpuStat() *CpuStat {
	c := &CpuStat{}
	c.Data = make(map[string][]gopsutil_cpu.TimesStat)
	return c
}

// CpuStat is a reader that scrapes cpu stat data.
// Data source: https://github.com/shirou/gopsutil/tree/master/cpu
type CpuStat struct {
	Data map[string][]gopsutil_cpu.TimesStat
}

// Run gathers gopsutil_cpu.CPUTimes data.
func (c *CpuStat) Run() error {
	data, err := gopsutil_cpu.Times(true)
	if err != nil {
		return err
	}

	c.Data["CpuStat"] = data
	return nil
}

// ToJson serialize Data field to JSON.
func (c *CpuStat) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
