package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
)

func NewMemory() *Memory {
	m := &Memory{}
	m.Data = make(map[string]map[string]interface{})
	return m
}

type Memory struct {
	Data map[string]map[string]interface{}
}

func (m *Memory) Run() error {
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
	m.Data["Memory"]["Used"] = mem.Used
	m.Data["Memory"]["Free"] = mem.Free
	m.Data["Memory"]["ActualUsed"] = mem.ActualUsed
	m.Data["Memory"]["ActualFree"] = mem.ActualFree

	m.Data["Swap"]["Total"] = swap.Total
	m.Data["Swap"]["Used"] = swap.Used
	m.Data["Swap"]["Free"] = swap.Free

	return nil
}

func (m *Memory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
