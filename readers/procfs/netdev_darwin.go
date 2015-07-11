// +build darwin

package procfs

import (
	"encoding/json"
	"errors"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("ProcNetDev", NewProcNetDev)
}

// NewProcNetDev is ProcNetDev constructor.
func NewProcNetDev() readers.IReader {
	p := &ProcNetDev{}
	p.Data = make(map[string][]interface{})
	return p
}

// ProcNetDev is a reader that scrapes /proc/net/dev data.
type ProcNetDev struct {
	Data map[string][]interface{}
}

func (p *ProcNetDev) Run() error {
	return errors.New("/proc/net/dev is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (p *ProcNetDev) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
