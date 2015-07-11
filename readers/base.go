// Package readers provides objects that gathers resource data from a host.
package readers

import (
	"errors"
	resourced_config "github.com/resourced/resourced/config"
	"reflect"
)

var readerConstructors = make(map[string]func() IReader)

// Register makes any reader constructor available by name.
func Register(name string, constructor func() IReader) {
	if constructor == nil {
		panic("reader: Register reader constructor is nil")
	}
	if _, dup := readerConstructors[name]; dup {
		panic("reader: Register called twice for reader constructor " + name)
	}
	readerConstructors[name] = constructor
}

// NewGoStruct instantiates IReader
func NewGoStruct(name string) (IReader, error) {
	constructor, ok := readerConstructors[name]
	if !ok {
		return nil, errors.New("GoStruct is undefined.")
	}

	return constructor(), nil
}

// NewGoStructByConfig instantiates IReader given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IReader, error) {
	reader, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	// Populate IReader fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(reader).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return reader, err
}

// IReader is generic interface for all readers.
type IReader interface {
	Run() error
	ToJson() ([]byte, error)
}
