// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcDiskStats", NewProcDiskStats)
}

// NewProcDiskStats is ProcDiskStats constructor.
func NewProcDiskStats() readers.IReader {
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

// ToJson serialize Data field to JSON.
func (p *ProcDiskStats) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
