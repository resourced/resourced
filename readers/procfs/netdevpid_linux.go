// +build linux

package procfs

import (
	"encoding/json"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/cloudfoundry/gosigar"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcNetDevPid", NewProcNetDevPid)
}

// NewProcNetDevPid is ProcNetDevPid constructor.
func NewProcNetDevPid() readers.IReader {
	p := &ProcNetDevPid{}
	p.Data = make(map[string]map[string]linuxproc.NetworkStat)
	return p
}

// ProcNetDevPid is a reader that scrapes /proc/$pid/net/dev data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/network_stat.go
type ProcNetDevPid struct {
	Data map[string]map[string]linuxproc.NetworkStat
}

func (p *ProcNetDevPid) Run() error {
	pids := sigar.ProcList{}
	err := pids.Get()
	if err != nil {
		return err
	}

	for _, pid := range pids.List {
		data, err := linuxproc.ReadNetworkStat(fmt.Sprintf("/proc/%v/net/dev", pid))
		if err != nil {
			return err
		}

		pidString := fmt.Sprintf("%v", pid)

		p.Data[pidString] = make(map[string]linuxproc.NetworkStat)

		for _, perIface := range data {
			p.Data[pidString][perIface.Iface] = perIface
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (p *ProcNetDevPid) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
