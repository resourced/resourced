// Package loggers provides objects that gathers resource data from a host.
package loggers

import (
	"errors"
	"os"
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
	SetLoglines(string, []string)
	GetLoglines(string) []string
	ResetLoglines(string)
	GetData() *libmap.TSafeMapStrings
	GetFiles() []string
	GetMaxLengthWipeTrigger() int64
}

type ILoggerFile interface {
	ILogger
	RunBlocking(string)
}

func NewBase() ILogger {
	b := &Base{}
	b.Data = libmap.NewTSafeMapStrings(nil)
	b.MaxLengthWipeTrigger = 1000000

	return b
}

type Base struct {
	Files                []interface{}
	Data                 *libmap.TSafeMapStrings
	MaxLengthWipeTrigger int64
	Target               string
	DenyList             []interface{}
}

// Run tails the file continuously.
func (b *Base) RunBlocking(file string) {
	t, err := tail.TailFile(file, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
	})
	if err == nil {
		if !b.Data.Exists(file) {
			b.Data.Set(file, make([]string, 0))
		}

		for line := range t.Lines {
			b.Data.Append(file, line.Text)
		}
	}
}

// SetLoglines sets loglines.
func (b *Base) SetLoglines(file string, loglines []string) {
	b.Data.Set(file, loglines)
}

// GetLoglines returns loglines.
func (b *Base) GetLoglines(file string) []string {
	return b.Data.Get(file)
}

// ResetLoglines wipes it clean.
func (b *Base) ResetLoglines(file string) {
	b.Data.Reset(file)
}

// GetData returns data.
func (b *Base) GetData() *libmap.TSafeMapStrings {
	return b.Data
}

// GetFiles returns data.
func (b *Base) GetFiles() []string {
	result := make([]string, len(b.Files))

	for i, fileInterface := range b.Files {
		result[i] = fileInterface.(string)
	}

	return result
}

// GetMaxLengthWipeTrigger returns MaxLengthWipeTrigger
func (b *Base) GetMaxLengthWipeTrigger() int64 {
	return b.MaxLengthWipeTrigger
}
