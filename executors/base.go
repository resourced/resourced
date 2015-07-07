// Package executors provides objects that gathers resource data from a host.
package executors

import (
	"reflect"
	"sync"

	resourced_config "github.com/resourced/resourced/config"
)

var executorConstructors = make(map[string]func() IExecutor)

var ConditionMetByPathCounter = make(map[string]int64)

// Register makes any executor constructor available by name.
func Register(name string, constructor func() IExecutor) {
	if constructor == nil {
		panic("executor: Register executor constructor is nil")
	}
	if _, dup := executorConstructors[name]; dup {
		panic("executor: Register called twice for executor constructor " + name)
	}
	executorConstructors[name] = constructor
}

// NewGoStruct instantiates IExecutor
func NewGoStruct(name string) IExecutor {
	return executorConstructors[name]()
}

// NewGoStructByConfig instantiates IExecutor given Config struct
func NewGoStructByConfig(config resourced_config.Config) IExecutor {
	executor := NewGoStruct(config.GoStruct)

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

	return executor
}

// IExecutor is generic interface for all executors.
type IExecutor interface {
	Run() error
	ToJson() ([]byte, error)
	ConditionMet()
	LowThresholdExceeded() bool
	HighThresholdExceeded() bool
}

type Base struct {
	Command       string
	Path          string
	Interval      string
	LowThreshold  int64
	HighThreshold int64
	Conditions    []interface{}
	sync.RWMutex
}

func (b *Base) ConditionMet() {
	b.Lock()
	ConditionMetByPathCounter[b.Path] = ConditionMetByPathCounter[b.Path] + 1
	b.Unlock()
}

func (b *Base) LowThresholdExceeded() bool {
	return ConditionMetByPathCounter[b.Path] > b.LowThreshold
}

func (b *Base) HighThresholdExceeded() bool {
	return ConditionMetByPathCounter[b.Path] > b.HighThreshold
}
