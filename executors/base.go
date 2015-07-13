// Package executors provides objects that gathers resource data from a host.
package executors

import (
	"errors"
	"reflect"
	"sync"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/queryparser"
)

var executorConstructors = make(map[string]func() IExecutor)

var ConditionMetByPathCounter = make(map[string]int)

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

func ResetConditionsMetByPath() {
	ConditionMetByPathCounter = make(map[string]int)
}

// NewGoStruct instantiates IExecutor
func NewGoStruct(name string) (IExecutor, error) {
	constructor, ok := executorConstructors[name]
	if !ok {
		return nil, errors.New("GoStruct is undefined.")
	}

	return constructor(), nil
}

// NewGoStructByConfig instantiates IExecutor given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IExecutor, error) {
	executor, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	executor.SetPath(config.Path)
	executor.SetInterval(config.Interval)

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

	executor.SetQueryParser()

	return executor, nil
}

// IExecutor is generic interface for all executors.
type IExecutor interface {
	SetPath(string)
	GetPath() string
	SetInterval(string)
	Run() error
	ToJson() ([]byte, error)
	SetQueryParser()
	SetReadersDataInBytes(map[string][]byte)
	IsConditionMet() bool
	LowThresholdExceeded() bool
	HighThresholdExceeded() bool
	ConditionMetCount() int
	OnLowThresholdExceeded() error
	OnHighThresholdExceeded() error
}

type Base struct {
	Command          string
	Path             string
	Interval         string
	LowThreshold     int
	HighThreshold    int
	Conditions       string
	ReadersDataBytes map[string][]byte
	qp               *queryparser.QueryParser
	sync.RWMutex
}

func (b *Base) SetPath(path string) {
	b.Path = path
}

func (b *Base) GetPath() string {
	return b.Path
}

func (b *Base) SetInterval(interval string) {
	b.Interval = interval
}

func (b *Base) SetQueryParser() {
	if b.Conditions == "" {
		b.Conditions = `[true]`
	}
	b.qp = queryparser.New([]byte(b.Conditions))
}

// SetReadersDataInBytes pulls readers data and store them on ReadersData field.
func (b *Base) SetReadersDataInBytes(readersJsonBytes map[string][]byte) {
	b.ReadersDataBytes = readersJsonBytes
}

func (b *Base) IsConditionMet() bool {
	result, err := b.qp.EvalJSONExpressions(b.ReadersDataBytes, nil)
	if err != nil {
		return false
	}
	if result == true {
		b.conditionMet()
	} else {
		b.conditionNoLongerMet()
	}
	return result
}

func (b *Base) conditionMet() {
	b.Lock()
	ConditionMetByPathCounter[b.Path] = ConditionMetByPathCounter[b.Path] + 1
	b.Unlock()
}

func (b *Base) conditionNoLongerMet() {
	b.Lock()
	ConditionMetByPathCounter[b.Path] = 0
	b.Unlock()
}

func (b *Base) ConditionMetCount() int {
	return ConditionMetByPathCounter[b.Path]
}

func (b *Base) LowThresholdExceeded() bool {
	if b.LowThreshold == 0 {
		return false
	}
	return ConditionMetByPathCounter[b.Path] > b.LowThreshold
}

func (b *Base) HighThresholdExceeded() bool {
	if b.HighThreshold == 0 {
		return false
	}
	return ConditionMetByPathCounter[b.Path] > b.HighThreshold
}

func (b *Base) OnLowThresholdExceeded() error {
	return nil
}

func (b *Base) OnHighThresholdExceeded() error {
	return nil
}
