// +build linux

package procfs

import (
	"encoding/json"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

// NewProcMounts is ProcMounts constructor.
func NewProcMounts() *ProcMounts {
	p := &ProcMounts{}
	p.Data = make(map[string][]linuxproc.Mount)
	return p
}

// ProcMounts is a reader that scrapes /proc/mounts data.
// Data source: https://github.com/c9s/goprocinfo/blob/master/linux/mounts.go
type ProcMounts struct {
	Data map[string][]linuxproc.Mount
}

func (p *ProcMounts) Run() error {
	mounts, err := linuxproc.ReadMounts("/proc/mounts")
	if err != nil {
		return err
	}

	p.Data["Mounts"] = mounts.Mounts
	return nil
}

func (p *ProcMounts) ToJson() ([]byte, error) {
	return json.Marshal(p.Data)
}
