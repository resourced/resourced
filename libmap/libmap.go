package libmap

import (
	"encoding/json"
	"strings"
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

func NewTSafeMapCounter() *TSafeMapCounter {
	s := &TSafeMapCounter{}
	s.Data = make(map[string]int)
	return s
}

type TSafeMapCounter struct {
	Data map[string]int
	sync.RWMutex
}

func (s *TSafeMapCounter) Incr(key string, value int) {
	s.Lock()
	s.Data[key] = s.Data[key] + value
	s.Unlock()
}

func (s *TSafeMapCounter) Get(key string) int {
	s.RLock()
	defer s.RUnlock()

	data, ok := s.Data[key]
	if !ok {
		data = 0
	}

	return data
}

func (s *TSafeMapCounter) Reset(key string) {
	s.Lock()
	s.Data[key] = 0
	s.Unlock()
}

func (s *TSafeMapCounter) ToJson() ([]byte, error) {
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
	// Split key by dot, loop deeper into the nesting & create the maps
	keyParts := strings.Split(key, ".")

	s.Lock()
	m := s.Data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		_, ok := m[keyPart]
		if !ok {
			m[keyPart] = make(map[string]interface{})
		}

		m = m[keyPart].(map[string]interface{})
	}
	s.Unlock()
}

func (s *TSafeNestedMapInterface) Set(key string, value interface{}) {
	s.initNestedMap(key)

	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	s.Lock()
	m := s.Data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	m[lastPart] = value
	s.Unlock()
}

func (s *TSafeNestedMapInterface) Get(key string) interface{} {
	var data interface{}

	// Split key by dot and loop deeper into the nesting
	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	s.RLock()
	m := s.Data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	data = m[lastPart]
	s.RUnlock()

	return data
}

func (s *TSafeNestedMapInterface) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
