// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcUptime is ProcUptime constructor.
func NewProcUptime() *ProcUptime {
	p := &ProcUptime{}
	return p
}

// ProcUptime is a reader that scrapes /proc/uptime data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/uptime.go
type ProcUptime struct {
	Data *linuxproc.Uptime
}

func (p *ProcUptime) Run() error {
	uptime, err := linuxproc.ReadUptime("/proc/uptime")
	if err != nil {
		return err
	}

	p.Data = uptime
	return nil
}

func (p *ProcUptime) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
