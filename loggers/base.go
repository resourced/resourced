// Package loggers provides objects that gathers resource data from a host.
package loggers

import (
	"errors"
	"reflect"

	"github.com/hpcloud/tail"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libmap"
)

var loggerConstructors = make(map[string]func() ILogger)

func init() {
	Register("Base", NewBase)
}

// Register makes any reader constructor available by name.
func Register(name string, constructor func() ILogger) {
	if constructor == nil {
		panic("reader: Register reader constructor is nil")
	}
	if _, dup := loggerConstructors[name]; dup {
		panic("reader: Register called twice for reader constructor " + name)
	}
	loggerConstructors[name] = constructor
}

// NewGoStruct instantiates ILogger
func NewGoStruct(name string) (ILogger, error) {
	constructor, ok := loggerConstructors[name]
	if !ok {
		return nil, errors.New("GoStruct is undefined.")
	}

	return constructor(), nil
}

// NewGoStructByConfig instantiates ILogger given Config struct
func NewGoStructByConfig(config resourced_config.Config) (ILogger, error) {
	reader, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	// Populate ILogger fields dynamically
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

// ILogger is generic interface for all readers.
type ILogger interface {
	RunBlocking()
	GetData() *libmap.TSafeMapStrings
	GetFile() string
	GetAutoPruneLength() int64
}

func NewBase() ILogger {
	b := &Base{}
	b.Data = libmap.NewTSafeMapStrings(map[string][]string{
		"Loglines": make([]string, 0),
	})
	b.AutoPruneLength = 1000000

	return b
}

type Base struct {
	File            string
	Data            *libmap.TSafeMapStrings
	AutoPruneLength int64
}

// Run tails the file continuously.
func (b *Base) RunBlocking() {
	t, err := tail.TailFile(b.File, tail.Config{Follow: true})
	if err == nil {
		for line := range t.Lines {
			b.Data.Append("Loglines", line.Text)
		}
	}
}

// GetData returns data.
func (b *Base) GetData() *libmap.TSafeMapStrings {
	return b.Data
}

// GetFile returns data.
func (b *Base) GetFile() string {
	return b.File
}

// GetAutoPruneLength returns AutoPruneLength
func (b *Base) GetAutoPruneLength() int64 {
	return b.AutoPruneLength
}
