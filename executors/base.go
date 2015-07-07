// Package executors provides objects that gathers resource data from a host.
package executors

import (
	"reflect"
	"sync"

	resourced_config "github.com/resourced/resourced/config"
)

var ConditionMetByPathCounter map[string]int64

func init() {
	ConditionMetByPathCounter = make(map[string]int64)
}

// NewGoStruct instantiates IExecutor
func NewGoStruct(name string) (IExecutor, error) {
	return nil, nil
}

// NewGoStructByConfig instantiates IExecutor given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IExecutor, error) {
	executor, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	// Populate IExecutor fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(executor).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return executor, err
}

// IExecutor is generic interface for all executors.
type IExecutor interface {
	Run() error
	ToJson() ([]byte, error)
	ConditionMet()
	LowTresholdExceeded() bool
	HighTresholdExceeded() bool
}

type Base struct {
	Command      string
	Path         string
	Interval     string
	LowTreshold  int64
	HighTreshold int64
	Conditions   []interface{}
	sync.RWMutex
}

func (b *Base) ConditionMet() {
	b.Lock()
	ConditionMetByPathCounter[b.Path] = ConditionMetByPathCounter[b.Path] + 1
	b.Unlock()
}

func (b *Base) LowTresholdExceeded() bool {
	return ConditionMetByPathCounter[b.Path] > b.LowTreshold
}

func (b *Base) HighTresholdExceeded() bool {
	return ConditionMetByPathCounter[b.Path] > b.HighTreshold
}
