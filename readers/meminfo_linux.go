// +build linux

package readers

import (
	"encoding/json"
	"github.com/guillermo/go.procmeminfo"
)

func NewMeminfo() *Meminfo {
	m := &Meminfo{}
	m.Data = make(map[string]uint64)
	return m
}

type Meminfo struct {
	Data map[string]uint64
}

func (m *Meminfo) Run() error {
	meminfo := &procmeminfo.MemInfo{}
	err := meminfo.Update()

	m.Data = *meminfo

	return err
}

func (m *Meminfo) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
