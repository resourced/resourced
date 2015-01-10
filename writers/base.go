package writers

import (
	"encoding/json"
	"errors"
)

// NewGoStruct instantiates IWriter
func NewGoStruct(name string) (IWriter, error) {
	var structInstance IWriter

	if name == "StdOut" {
		structInstance = NewStdOut()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

type IWriter interface {
	Run() error
	SetReadersData(map[string][]byte)
	GetReadersData() map[string]interface{}
	ToJson() ([]byte, error)
}

type Base struct {
	ReadersData map[string]interface{}
	Data        map[string]interface{}
}

func (b *Base) Run() error {
	return nil
}

func (b *Base) SetReadersData(readersJsonBytes map[string][]byte) {
	if b.ReadersData == nil {
		b.ReadersData = make(map[string]interface{})
	}

	for key, jsonBytes := range readersJsonBytes {
		var data interface{}
		err := json.Unmarshal(jsonBytes, &data)
		if err == nil {
			b.ReadersData[key] = data
		}
	}
}

func (b *Base) GetReadersData() map[string]interface{} {
	return b.ReadersData
}

func (b *Base) ToJson() ([]byte, error) {
	return json.Marshal(b.Data)
}
