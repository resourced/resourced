// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcVmStat", NewProcVmStat)
}

// NewProcVmStat is ProcVmStat constructor.
func NewProcVmStat() readers.IReader {
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
