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
	Data map[string]interface{}
}

func (n *Noop) Run() error {
	return nil
}

func (n *Noop) SetData(jsonBytes []byte) error {
	var fullJsonMap map[string]interface{}

	err := json.Unmarshal(jsonBytes, &fullJsonMap)
	if err != nil {
		return err
	}

	n.Data = fullJsonMap["Data"].(map[string]interface{})
	return nil
}
