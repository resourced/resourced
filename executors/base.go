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

var ConditionMetByPathCounterMutex = &sync.RWMutex{}

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

	return executor, nil
}

// IExecutor is generic interface for all executors.
type IExecutor interface {
	SetPath(string)
	GetPath() string
	SetInterval(string)
	Run() error
	ToJson() ([]byte, error)
	SetQueryParser(map[string][]byte)
	SetReadersDataInBytes(map[string][]byte)
	SetTags(map[string]string)
	IsConditionMet() bool
	LowThresholdExceeded() bool
	HighThresholdExceeded() bool
	ConditionMetCount() int
}

type Base struct {
	// Command: Shell command to execute.
	Command string

	// Path: ResourceD URL path. Example:
	// /uptime -> http://localhost:55555/x/uptime
	Path string

	Interval string

	// LowThreshold: minimum count of valid conditions
	LowThreshold int

	// HighThreshold: maximum count of valid conditions
	HighThreshold int

	// Conditions for when executor should run. It uses javascript.
	Conditions string

	ReadersDataBytes map[string][]byte

	qp *queryparser.QueryParser
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

func (b *Base) SetQueryParser(readersJsonBytes map[string][]byte) {
	b.qp = queryparser.New(readersJsonBytes, nil)
}

// SetReadersDataInBytes pulls readers data and store them on ReadersData field.
func (b *Base) SetReadersDataInBytes(readersJsonBytes map[string][]byte) {
	b.ReadersDataBytes = readersJsonBytes

	b.SetQueryParser(readersJsonBytes)
}

// SetTags assigns all host tags to qp (QueryParser).
func (b *Base) SetTags(tags map[string]string) {
	b.qp.SetTags(tags)
}

func (b *Base) IsConditionMet() bool {
	if b.Conditions == "" {
		b.Conditions = "true"
	}

	result, err := b.qp.Parse(b.Conditions)
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
	ConditionMetByPathCounterMutex.Lock()
	defer ConditionMetByPathCounterMutex.Unlock()

	ConditionMetByPathCounter[b.Path] = b.ConditionMetCount() + 1
}

func (b *Base) conditionNoLongerMet() {
	ConditionMetByPathCounterMutex.Lock()
	defer ConditionMetByPathCounterMutex.Unlock()

	ConditionMetByPathCounter[b.Path] = 0
}

func (b *Base) ConditionMetCount() int {
	ConditionMetByPathCounterMutex.RLock()
	defer ConditionMetByPathCounterMutex.RUnlock()

	return ConditionMetByPathCounter[b.Path]
}

func (b *Base) LowThresholdExceeded() bool {
	ConditionMetByPathCounterMutex.RLock()
	defer ConditionMetByPathCounterMutex.RUnlock()

	return b.ConditionMetCount() > b.LowThreshold
}

func (b *Base) HighThresholdExceeded() bool {
	ConditionMetByPathCounterMutex.RLock()
	defer ConditionMetByPathCounterMutex.RUnlock()

	if b.HighThreshold == 0 {
		return false
	}
	return b.ConditionMetCount() > b.HighThreshold
}
