package readers

import (
	"encoding/json"
	"errors"
)

func NewGoStruct(name string) (IReaderWriter, error) {
	var structInstance IReaderWriter

	if name == "NetworkInterfaces" {
		structInstance = NewNetworkInterfaces()
	}
	if name == "Df" {
		structInstance = NewDf()
	}
	if name == "Du" {
		structInstance = NewDu()
	}
	if name == "Memory" {
		structInstance = NewMemory()
	}
	if name == "Ps" {
		structInstance = NewPs()
	}
	if name == "LoadAvg" {
		structInstance = NewLoadAvg()
	}
	if name == "Uptime" {
		structInstance = NewUptime()
	}
	if name == "Meminfo" {
		structInstance = NewMeminfo()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

type IReaderWriter interface {
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
