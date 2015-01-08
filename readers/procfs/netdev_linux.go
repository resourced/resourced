// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcNetDev is ProcNetDev constructor.
func NewProcNetDev() *ProcNetDev {
	p := &ProcNetDev{}
	p.Data = make(map[string]linuxproc.NetworkStat)
	return p
}

// ProcNetDev is a reader that scrapes /proc/net/dev data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/network_stat.go
type ProcNetDev struct {
	Data map[string]linuxproc.NetworkStat
}

func (p *ProcNetDev) Run() error {
	data, err := linuxproc.ReadNetworkStat("/proc/net/dev")
	if err != nil {
		return err
	}

	for _, perIface := range data {
		p.Data[perIface.Iface] = perIface
	}

	return nil
}

func (p *ProcNetDev) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
