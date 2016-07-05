// +build darwin
package readers

import (
	"encoding/json"
	"errors"
)

func init() {
	Register("IOStat", NewIOStat)
}

func NewIOStat() IReader {
	ios := &IOStat{}
	ios.Data = make(map[string]map[string]float64)
	return ios
}

type IOStat struct {
	Data map[string]map[string]float64
}

// Run gathers load average information from gosigar.
func (ios *IOStat) Run() error {
	return errors.New("iostat -x is only available on Linux.")
}

// ToJson serialize Data field to JSON.
func (ios *IOStat) ToJson() ([]byte, error) {
	return json.Marshal(ios.Data)
}
