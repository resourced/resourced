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
