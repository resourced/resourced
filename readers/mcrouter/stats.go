package mcrouter

import (
	"encoding/json"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("McRouterStats", NewMcRouterStats)
}

func NewMcRouterStats() readers.IReader {
	r := &McRouterStats{}
	r.Data = make(map[string]interface{})

	return r
}

type McRouterStats struct {
	Data map[string]interface{}
	Base
}

func (r *McRouterStats) Run() error {
	data, err := r.Stats()
	if err != nil {
		return err
	}

	r.Data = data

	return nil
}

// ToJson serialize Data field to JSON.
func (r *McRouterStats) ToJson() ([]byte, error) {
	return json.Marshal(r.Data)
}
