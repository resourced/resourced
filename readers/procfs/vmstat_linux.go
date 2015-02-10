// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcVmStat is ProcVmStat constructor.
func NewProcVmStat() *ProcVmStat {
	p := &ProcVmStat{}
	return p
}

// ProcVmStat is a reader that scrapes /proc/vmstat data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/vmstat.go
type ProcVmStat struct {
	Data *linuxproc.VMStat
}

func (p *ProcVmStat) Run() error {
	data, err := linuxproc.ReadVMStat("/proc/vmstat")
	if err != nil {
		return err
	}

	p.Data = data
	return nil
}

// ToJson serialize Data field to JSON.
func (p *ProcVmStat) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
