package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
)

func NewPs() *Ps {
	p := &Ps{}
	p.Data = make([]map[string]interface{}, 0)
	return p
}

type Ps struct {
	Base
	Data []map[string]interface{}
}

func (p *Ps) Run() error {
	pids := sigar.ProcList{}
	err := pids.Get()
	if err != nil {
		return err
	}

	for _, pid := range pids.List {
		state := sigar.ProcState{}
		mem := sigar.ProcMem{}
		time := sigar.ProcTime{}

		if err := state.Get(pid); err != nil {
			continue
		}
		if err := mem.Get(pid); err != nil {
			continue
		}
		if err := time.Get(pid); err != nil {
			continue
		}

		procData := make(map[string]interface{})
		procData["Pid"] = pid
		procData["Ppid"] = state.Ppid
		procData["StartTime"] = time.FormatStartTime()
		procData["FormatTotal"] = time.FormatTotal()
		procData["MemoryResident"] = mem.Resident / 1024
		procData["State"] = state.State
		procData["Name"] = state.Name

		if len(procData) > 0 {
			p.Data = append(p.Data, procData)
		}
	}

	return nil
}

func (p *Ps) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
