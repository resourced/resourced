// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcCpuInfo is ProcCpuInfo constructor.
func NewProcCpuInfo() *ProcCpuInfo {
	c := &ProcCpuInfo{}
	c.Data = make(map[string][]interface{})
	return c
}

// ProcCpuInfo is a reader that scrapes /proc/cpuinfo data.
type ProcCpuInfo struct {
	Data map[string][]interface{}
}

func (c *ProcCpuInfo) Run() error {
	return errors.New("/proc/cpuinfo is only available on Linux.")
}

func (c *ProcCpuInfo) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
