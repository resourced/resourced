package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
)

func NewPs() *Ps {
	p := &Ps{}
	p.Data = make(map[string][]map[string]interface{})
	return p
}

type Ps struct {
	Data map[string][]map[string]interface{}
}

func (p *Ps) Run() error {
	pids := sigar.ProcList{}
	err := pids.Get()
	if err != nil {
		return err
	}

	p.Data["Processes"] = make([]map[string]interface{}, 0)

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
		procData["Name"] = state.Name
		procData["Pid"] = pid
		procData["ParentPid"] = state.Ppid
		procData["StartTime"] = time.FormatStartTime()
		procData["RunTime"] = time.FormatTotal()
		procData["MemoryResident"] = mem.Resident / 1024
		procData["State"] = string(state.State)

		if len(procData) > 0 {
			p.Data["Processes"] = append(p.Data["Processes"], procData)
		}
	}

	return nil
}

func (p *Ps) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
