// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
)

// NewProcDiskStats is ProcDiskStats constructor.
func NewProcDiskStats() *ProcDiskStats {
	p := &ProcDiskStats{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcDiskStats is a reader that scrapes /proc/diskstats data.
type ProcDiskStats struct {
	Data map[string][]interface{}
}

func (p *ProcDiskStats) Run() error {
	return errors.New("/proc/diskstats is only available on Linux.")
}

func (p *ProcDiskStats) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
