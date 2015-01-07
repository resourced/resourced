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
	InputData map[string]interface{}
	Data      map[string]interface{}
}

func (n *Noop) Run() error {
	return nil
}

func (n *Noop) SetData(jsonBytes []byte) error {
	err := json.Unmarshal(jsonBytes, &n.InputData)
	if err != nil {
		return err
	}
	return err
}

func (n *Noop) ToJson() ([]byte, error) {
	return json.Marshal(n.Data)
}
