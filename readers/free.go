package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
)

func init() {
	Register("Free", NewFree)
}

func NewFree() IReader {
	m := &Free{}
	m.Data = make(map[string]map[string]interface{})
	return m
}

// Free is a reader that scrapes swapory data and presents it in the form similar to `free`.
// Data source: https://github.com/cloudfoundry/gosigar/tree/master
type Free struct {
	Data map[string]map[string]interface{}
}

// Run gathers free data from gosigar.
func (m *Free) Run() error {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		return err
	}

	swap := sigar.Swap{}
	err = swap.Get()
	if err != nil {
		return err
	}

	m.Data["Memory"] = make(map[string]interface{})
	m.Data["Swap"] = make(map[string]interface{})

	m.Data["Memory"]["Total"] = mem.Total
	m.Data["Memory"]["TotalMB"] = mem.Total / 1000 / 1000
	m.Data["Memory"]["TotalGB"] = mem.Total / 1000 / 1000 / 1000

	m.Data["Memory"]["Used"] = mem.Used
	m.Data["Memory"]["UsedMB"] = mem.Used / 1000 / 1000
	m.Data["Memory"]["UsedGB"] = mem.Used / 1000 / 1000 / 1000

	m.Data["Memory"]["UsedPercent"] = int(float64(mem.Used) / float64(mem.Total) * 100)

	m.Data["Memory"]["Free"] = mem.Free
	m.Data["Memory"]["FreeMB"] = mem.Free / 1000 / 1000
	m.Data["Memory"]["FreeGB"] = mem.Free / 1000 / 1000 / 1000

	m.Data["Memory"]["FreePercent"] = int(float64(mem.Free) / float64(mem.Total) * 100)

	m.Data["Memory"]["ActualUsed"] = mem.ActualUsed
	m.Data["Memory"]["ActualUsedMB"] = mem.ActualUsed / 1000 / 1000
	m.Data["Memory"]["ActualUsedGB"] = mem.ActualUsed / 1000 / 1000 / 1000

	m.Data["Memory"]["ActualUsedPercent"] = int(float64(mem.ActualUsed) / float64(mem.Total) * 100)

	m.Data["Memory"]["ActualFree"] = mem.ActualFree
	m.Data["Memory"]["ActualFreeMB"] = mem.ActualFree / 1000 / 1000
	m.Data["Memory"]["ActualFreeGB"] = mem.ActualFree / 1000 / 1000 / 1000

	m.Data["Memory"]["ActualFreePercent"] = int(float64(mem.ActualFree) / float64(mem.Total) * 100)

	m.Data["Swap"]["Total"] = swap.Total
	m.Data["Memory"]["TotalMB"] = swap.Total / 1000 / 1000
	m.Data["Memory"]["TotalGB"] = swap.Total / 1000 / 1000 / 1000

	m.Data["Swap"]["Used"] = swap.Used
	m.Data["Memory"]["UsedMB"] = swap.Used / 1000 / 1000
	m.Data["Memory"]["UsedGB"] = swap.Used / 1000 / 1000 / 1000

	m.Data["Swap"]["Free"] = swap.Free
	m.Data["Swap"]["FreeMB"] = swap.Free / 1000 / 1000
	m.Data["Swap"]["FreeGB"] = swap.Free / 1000 / 1000 / 1000

	m.Data["Memory"]["SwapFreePercent"] = int(float64(swap.Free) / float64(swap.Total) * 100)

	return nil
}

// ToJson serialize Data field to JSON.
func (m *Free) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
