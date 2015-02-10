package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	gopsutil_process "github.com/shirou/gopsutil/process"
	"strconv"
)

func NewPs() *Ps {
	p := &Ps{}
	p.Data = make(map[string]map[string]interface{})
	return p
}

type Ps struct {
	Data map[string]map[string]interface{}
}

// Run gathers ps information from gosigar.
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
		procData["Name"] = state.Name
		procData["Pid"] = pid
		procData["ParentPid"] = state.Ppid
		procData["StartTime"] = time.FormatStartTime()
		procData["RunTime"] = time.FormatTotal()
		procData["MemoryResident"] = mem.Resident / 1024
		procData["State"] = string(state.State)

		gopsutilProcess, err := gopsutil_process.NewProcess(int32(pid))
		if err != nil {
			continue
		}

		mmaps, err := gopsutilProcess.MemoryMaps(false)
		if err == nil {
			procData["MemoryMaps"] = mmaps
		}

		ios, err := gopsutilProcess.IOCounters()
		if err == nil {
			procData["IOCounters"] = ios
		}

		ctxSwitches, err := gopsutilProcess.NumCtxSwitches()
		if err == nil {
			procData["CtxSwitches"] = ctxSwitches
		}

		if len(procData) > 0 {
			p.Data[strconv.Itoa(pid)] = procData
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (p *Ps) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
