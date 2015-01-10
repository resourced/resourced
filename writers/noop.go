package writers

import (
	"encoding/json"
)

// NewNoop is Noop constructor.
func NewNoop() *Noop {
	n := &Noop{}
	return n
}

// Noop is a writer that does not do anything.
type Noop struct {
	ReadersData map[string]interface{}
	Data        map[string]interface{}
}

func (n *Noop) Run() error {
	return nil
}

func (n *Noop) SetReadersData(readersJsonBytes map[string][]byte) {
	if n.ReadersData == nil {
		n.ReadersData = make(map[string]interface{})
	}

	for key, jsonBytes := range readersJsonBytes {
		var data interface{}
		err := json.Unmarshal(jsonBytes, &data)
		if err == nil {
			n.ReadersData[key] = data
		}
	}
}

func (n *Noop) GetReadersData() map[string]interface{} {
	return n.ReadersData
}

func (n *Noop) ToJson() ([]byte, error) {
	return json.Marshal(n.Data)
}
