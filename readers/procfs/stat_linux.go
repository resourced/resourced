// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcStat", NewProcStat)
}

// NewProcStat is ProcStat constructor.
func NewProcStat() readers.IReader {
	p := &ProcStat{}
	return p
}

// ProcStat is a reader that scrapes /proc/stat data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/stat.go
type ProcStat struct {
	Data *linuxproc.Stat
}

func (p *ProcStat) Run() error {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return err
	}

	p.Data = stat
	return nil
}

// ToJson serialize Data field to JSON.
func (p *ProcStat) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
