package libmap

import (
	"encoding/json"
	"strings"
	"sync"
)

func NewTSafeMapBytes(data map[string][]byte) *TSafeMapBytes {
	s := &TSafeMapBytes{}
	if data == nil {
		s.data = make(map[string][]byte)
	} else {
		s.data = data
	}
	return s
}

type TSafeMapBytes struct {
	data map[string][]byte
	sync.RWMutex
}

func (s *TSafeMapBytes) Set(key string, value []byte) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *TSafeMapBytes) Get(key string) []byte {
	s.Lock()
	defer s.Unlock()

	return s.data[key]
}

func (s *TSafeMapBytes) All() map[string][]byte {
	s.Lock()
	defer s.Unlock()

	copydata := make(map[string][]byte)
	for key, value := range s.data {
		copydata[key] = value
	}

	return copydata
}

func (s *TSafeMapBytes) ToJson() ([]byte, error) {
	return json.Marshal(s.data)
}

func NewTSafeMapCounter(data map[string]int) *TSafeMapCounter {
	s := &TSafeMapCounter{}
	if data == nil {
		s.data = make(map[string]int)
	} else {
		s.data = data
	}
	return s
}

type TSafeMapCounter struct {
	data map[string]int
	sync.RWMutex
}

func (s *TSafeMapCounter) Incr(key string, value int) {
	s.Lock()
	s.data[key] = s.data[key] + value
	s.Unlock()
}

func (s *TSafeMapCounter) Get(key string) int {
	s.RLock()
	defer s.RUnlock()

	data, ok := s.data[key]
	if !ok {
		data = 0
	}

	return data
}

func (s *TSafeMapCounter) Reset(key string) {
	s.Lock()
	s.data[key] = 0
	s.Unlock()
}

func (s *TSafeMapCounter) All() map[string]int {
	s.Lock()
	defer s.Unlock()

	copydata := make(map[string]int)
	for key, value := range s.data {
		copydata[key] = value
	}

	return copydata
}

func (s *TSafeMapCounter) ToJson() ([]byte, error) {
	return json.Marshal(s.data)
}

func NewTSafeNestedMapInterface(data map[string]interface{}) *TSafeNestedMapInterface {
	s := &TSafeNestedMapInterface{}
	if data == nil {
		s.data = make(map[string]interface{})
	} else {
		s.data = data
	}
	return s
}

type TSafeNestedMapInterface struct {
	data map[string]interface{}
	sync.RWMutex
}

func (s *TSafeNestedMapInterface) initNestedMap(key string) {
	// Split key by dot, loop deeper into the nesting & create the maps
	keyParts := strings.Split(key, ".")

	s.Lock()
	m := s.data

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
	m := s.data

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
	m := s.data

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

func (s *TSafeNestedMapInterface) All() map[string]interface{} {
	s.Lock()
	defer s.Unlock()

	copydata := make(map[string]interface{})
	for key, value := range s.data {
		copydata[key] = value
	}

	return copydata
}

func (s *TSafeNestedMapInterface) ToJson() ([]byte, error) {
	return json.Marshal(s.data)
}
