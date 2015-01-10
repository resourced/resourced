package writers

import (
	"errors"
)

// NewGoStruct instantiates IWriter
func NewGoStruct(name string) (IWriter, error) {
	var structInstance IWriter

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
