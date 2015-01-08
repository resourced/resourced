// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcLoadAvg is ProcLoadAvg constructor.
func NewProcLoadAvg() *ProcLoadAvg {
	c := &ProcLoadAvg{}
	c.Data = make(map[string]float64)
	return c
}

// ProcLoadAvg is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/loadavg.go
type ProcLoadAvg struct {
	Data map[string]float64
}

func (c *ProcLoadAvg) Run() error {
	loadavg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		return err
	}

	c.Data["last1min"] = loadavg.Last1Min
	c.Data["last5min"] = loadavg.Last5Min
	c.Data["last15min"] = loadavg.Last15Min

	return nil
}

func (c *ProcLoadAvg) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
