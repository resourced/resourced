// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcLoadAvg is ProcLoadAvg constructor.
func NewProcLoadAvg() *ProcLoadAvg {
	p := &ProcLoadAvg{}
	p.Data = make(map[string]float64)
	return p
}

// ProcLoadAvg is a reader that scrapes /proc/diskstats data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/loadavg.go
type ProcLoadAvg struct {
	Data map[string]float64
}

func (p *ProcLoadAvg) Run() error {
	loadavg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		return err
	}

	p.Data["last1min"] = loadavg.Last1Min
	p.Data["last5min"] = loadavg.Last5Min
	p.Data["last15min"] = loadavg.Last15Min

	return nil
}

// ToJson serialize Data field to JSON.
func (p *ProcLoadAvg) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
