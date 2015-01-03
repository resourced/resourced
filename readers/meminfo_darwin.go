// +build darwin

package readers

import (
	"encoding/json"
)

func NewMeminfo() *Meminfo {
	m := &Meminfo{}
	m.Data = make(map[string]string)
	return m
}

type Meminfo struct {
	Base
	Data map[string]string
}

func (m *Meminfo) Run() error {
	m.Data["Error"] = "/proc/meminfo is only available on Linux."
	return nil
}

func (m *Meminfo) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
