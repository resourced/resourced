package libmap

import (
	"encoding/json"
	"sync"
)

func NewTSafeMapBytes() *TSafeMapBytes {
	s := &TSafeMapBytes{}
	s.Data = make(map[string][]byte)
	return s
}

type TSafeMapBytes struct {
	Data map[string][]byte
	sync.RWMutex
}

func (s *TSafeMapBytes) Set(key string, value []byte) {
	s.Lock()
	s.Data[key] = value
	s.Unlock()
}

func (s *TSafeMapBytes) Get(key string) []byte {
	var data []byte

	s.Lock()
	data = s.Data[key]
	s.Unlock()

	return data
}

func (s *TSafeMapBytes) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}

func NewTSafeNestedMapInterface() *TSafeNestedMapInterface {
	s := &TSafeNestedMapInterface{}
	s.Data = make(map[string]interface{})
	return s
}

type TSafeNestedMapInterface struct {
	Data map[string]interface{}
	sync.RWMutex
}

func (s *TSafeNestedMapInterface) initNestedMap(key string) {
	// Split key by dot and loop deeper into the nesting & create the maps
	s.Lock()
	s.Unlock()
}

func (s *TSafeNestedMapInterface) Set(key string, value interface{}) {
	s.initNestedMap(key)

	s.Lock()
	s.Data[key] = value
	s.Unlock()
}

func (s *TSafeNestedMapInterface) Get(key string) interface{} {
	var data interface{}

	// Split key by dot and loop deeper into the nesting

	s.RLock()
	data = s.Data[key]
	s.RUnlock()

	return data
}

func (s *TSafeNestedMapInterface) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
