package memcache

import (
	"encoding/json"

	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("MemcacheStats", NewMemcacheStats)
}

func NewMemcacheStats() readers.IReader {
	r := &MemcacheStats{}
	r.Data = make(map[string]interface{})

	return r
}

type MemcacheStats struct {
	Data map[string]interface{}
	Base
}

func (r *MemcacheStats) Run() error {
	data, err := r.Stats()
	if err != nil {
		return err
	}

	r.Data = data

	return nil
}

// ToJson serialize Data field to JSON.
func (r *MemcacheStats) ToJson() ([]byte, error) {
	return json.Marshal(r.Data)
}
