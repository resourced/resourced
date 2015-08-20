package varnish

import (
	"encoding/json"
	"os/exec"

	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("VarnishStats", NewVarnishStats)
}

func NewVarnishStats() readers.IReader {
	v := &VarnishStats{}
	v.Data = make(map[string]interface{})

	return v
}

type VarnishStats struct {
	Data map[string]interface{}
}

func (v *VarnishStats) Run() error {
	output, err := exec.Command("varnishstat", "-1", "-j").CombinedOutput()
	if err != nil {
		return err
	}

	return json.Unmarshal(output, &v.Data)
}

// ToJson serialize Data field to JSON.
func (v *VarnishStats) ToJson() ([]byte, error) {
	return json.Marshal(v.Data)
}
