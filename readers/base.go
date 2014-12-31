package readers

import (
	"encoding/json"
)

type IReader interface {
	Run() error
	ToJson() ([]byte, error)
}

type Base struct {
	Data map[string]interface{}
}

func (b *Base) Run() error {
	return nil
}

func (b *Base) ToJson() ([]byte, error) {
	return json.Marshal(b.Data)
}
