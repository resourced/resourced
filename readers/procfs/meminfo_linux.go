// +build linux

package procfs

import (
	"encoding/json"
	"github.com/guillermo/go.procmeminfo"
)

// NewProcMemInfo is ProcMemInfo constructor.
func NewProcMemInfo() *ProcMemInfo {
	p := &ProcMemInfo{}
	p.Data = make(map[string]uint64)
	return p
}

// ProcMemInfo is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/loadavg.go
type ProcMemInfo struct {
	Data map[string]uint64
}

func (p *ProcMemInfo) Run() error {
	meminfo := &procmeminfo.MemInfo{}
	err := meminfo.Update()

	p.Data = *meminfo

	return err
}

func (p *ProcMemInfo) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
