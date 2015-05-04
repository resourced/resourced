// Package writers provides objects that can send colected resource data to external place.
package writers

import (
	"bytes"
	"encoding/json"
	"errors"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libprocess"
	"github.com/resourced/resourced/libstring"
	"reflect"
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
	if name == "NewrelicInsights" {
		structInstance = NewNewrelicInsights()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

// NewGoStructByConfig instantiates IWriter given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IWriter, error) {
	writer, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	// Populate IWriter fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(writer).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return writer, err
}

// IWriter is general interface for writer.
type IWriter interface {
	Run() error
	SetReadersDataInBytes(map[string][]byte)
	SetReadersData(map[string]interface{})
	GetReadersData() map[string]interface{}
	SetData(interface{})
	GetData() interface{}
	GetJsonProcessor() string
	GenerateData() error
	ToJson() ([]byte, error)
}

type Base struct {
	ReadersData   map[string]interface{}
	Data          interface{}
	JsonProcessor string
}

// Run executes the writer.
func (b *Base) Run() error {
	return nil
}

// SetReadersDataInBytes pulls readers data and store them on ReadersData field.
func (b *Base) SetReadersDataInBytes(readersJsonBytes map[string][]byte) {
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

// SetReadersData assigns ReadersData field.
func (b *Base) SetReadersData(readersData map[string]interface{}) {
	b.ReadersData = readersData
}

// GetReadersData returns ReadersData field.
func (b *Base) GetReadersData() map[string]interface{} {
	return b.ReadersData
}

// SetData assigns Data field.
func (b *Base) SetData(data interface{}) {
	b.Data = data
}

// GetData returns Data field.
func (b *Base) GetData() interface{} {
	return b.Data
}

// GetJsonProcessor returns json processor path.
func (b *Base) GetJsonProcessor() string {
	path := ""
	if b.JsonProcessor != "" {
		path = libstring.ExpandTildeAndEnv(b.JsonProcessor)
	}
	return path
}

// GenerateData pulls ReadersData field and set it to Data field.
// If JsonProcessor is defined, use it to mangle JSON and save the new JSON on Data field.
func (b *Base) GenerateData() error {
	var err error

	processorPath := b.GetJsonProcessor()
	if processorPath == "" {
		// If there's no JsonProcessor
		b.SetData(b.GetReadersData())

	} else {
		// If there is a JsonProcessor
		cmd := libprocess.NewCmd(processorPath)

		readersData := b.GetReadersData()

		readersDataJsonBytes, err := json.Marshal(readersData)
		if err != nil {
			return err
		}

		cmd.Stdin = bytes.NewReader(readersDataJsonBytes)

		postProcessingDataBytes, err := cmd.Output()
		if err != nil {
			return err
		}

		var postProcessingData interface{}
		err = json.Unmarshal(postProcessingDataBytes, &postProcessingData)
		if err != nil {
			return err
		}

		b.SetData(postProcessingData)
	}

	return err
}

// ToJson serialize Data field to JSON.
func (b *Base) ToJson() ([]byte, error) {
	return json.Marshal(b.Data)
}
