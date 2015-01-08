// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcDiskStats is ProcDiskStats constructor.
func NewProcDiskStats() *ProcDiskStats {
	c := &ProcDiskStats{}
	c.Data = make(map[string]linuxproc.DiskStat)
	return c
}

// ProcDiskStats is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/diskstat.go
type ProcDiskStats struct {
	Data map[string]linuxproc.DiskStat
}

func (c *ProcDiskStats) Run() error {
	diskstats, err := linuxproc.ReadDiskStats("/proc/diskstats")
	if err != nil {
		return err
	}

	for _, perDevice := range diskstats {
		c.Data[perDevice.Name] = perDevice
	}
	return nil
}

func (c *ProcDiskStats) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
