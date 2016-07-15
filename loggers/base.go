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

// Register makes any logger constructor available by name.
func Register(name string, constructor func() ILogger) {
	if constructor == nil {
		panic("logger: Register logger constructor is nil")
	}
	if _, dup := loggerConstructors[name]; dup {
		panic("logger: Register called twice for logger constructor " + name)
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
	lgr, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	lgr.SetSource(config.Source)
	lgr.SetBufferSize(config.BufferSize)
	lgr.SetTargets(config.Targets)

	// Populate ILogger fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(lgr).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return lgr, err
}

// ILogger is generic interface for all loggers.
type ILogger interface {
	SetSource(string)
	GetSource() string
	SetBufferSize(int64)
	GetBufferSize() int64
	SetTargets([]resourced_config.LogTargetConfig)
	SetLoglines(string, []string)
	GetLoglines(string) []string
	GetLoglinesLength(string) int
	ResetLoglines(string)
	GetTargets() []resourced_config.LogTargetConfig
}

type ILoggerChannel interface {
	ILogger
	RunBlockingChannel(string, <-chan interface{})
}

type ILoggerFile interface {
	ILogger
	RunBlockingFile(string)
}

func NewBase() ILogger {
	b := &Base{}
	b.Data = libmap.NewTSafeMapStrings(nil)
	b.BufferSize = 1000000

	return b
}

type Base struct {
	Source     string
	BufferSize int64
	Targets    []resourced_config.LogTargetConfig

	Data *libmap.TSafeMapStrings
}

// RunBlockingChannel pulls log line from channel continuously.
func (b *Base) RunBlockingChannel(name string, ch <-chan interface{}) {
	for line := range ch {
		b.Data.Append(name, line.(string))
	}
}

// RunBlockingFile tails the file continuously.
func (b *Base) RunBlockingFile(file string) {
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

// SetSource
func (b *Base) SetSource(source string) {
	b.Source = source
}

// GetSource returns the source field.
func (b *Base) GetSource() string {
	return b.Source
}

// SetBufferSize sets BufferSize
func (b *Base) SetBufferSize(bufferSize int64) {
	b.BufferSize = bufferSize
}

// GetBufferSize returns BufferSize
func (b *Base) GetBufferSize() int64 {
	return b.BufferSize
}

// SetTargets sets []LogTargetConfig
func (b *Base) SetTargets(targets []resourced_config.LogTargetConfig) {
	b.Targets = targets
}

// SetLoglines sets loglines.
func (b *Base) SetLoglines(file string, loglines []string) {
	b.Data.Set(file, loglines)
}

// GetLoglines returns loglines.
func (b *Base) GetLoglines(file string) []string {
	return b.Data.Get(file)
}

// GetLoglinesLength returns the count of loglines.
func (b *Base) GetLoglinesLength(file string) int {
	return len(b.Data.Get(file))
}

// ResetLoglines wipes it clean.
func (b *Base) ResetLoglines(file string) {
	b.Data.Reset(file)
}

// GetTargets returns slice of LogTargetConfig.
func (b *Base) GetTargets() []resourced_config.LogTargetConfig {
	return b.Targets
}
