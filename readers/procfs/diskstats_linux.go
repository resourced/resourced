// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcDiskStats is ProcDiskStats constructor.
func NewProcDiskStats() *ProcDiskStats {
	p := &ProcDiskStats{}
	p.Data = make(map[string]linuxproc.DiskStat)
	return p
}

// ProcDiskStats is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/diskstat.go
type ProcDiskStats struct {
	Data map[string]linuxproc.DiskStat
}

func (p *ProcDiskStats) Run() error {
	diskstats, err := linuxproc.ReadDiskStats("/proc/diskstats")
	if err != nil {
		return err
	}

	for _, perDevice := range diskstats {
		p.Data[perDevice.Name] = perDevice
	}
	return nil
}

func (p *ProcDiskStats) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
