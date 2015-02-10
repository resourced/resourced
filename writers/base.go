// Package writers provides objects that can send colected resource data to external place.
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
	if name == "Http" {
		structInstance = NewHttp()
	}
	if name == "ResourcedMaster" {
		structInstance = NewResourcedMaster()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

// IWriter is general interface for writer.
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

// Run executes the writer.
func (b *Base) Run() error {
	return nil
}

// SetReadersData pulls readers data and store them on ReadersData field.
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

// GetReadersData returns ReadersData field.
func (b *Base) GetReadersData() map[string]interface{} {
	return b.ReadersData
}

// ToJson serialize Data field to JSON.
func (b *Base) ToJson() ([]byte, error) {
	return json.Marshal(b.Data)
}
