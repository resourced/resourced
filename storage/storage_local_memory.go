package storage

import (
	"encoding/json"
	"sync"
)

func NewStorage() *Storage {
	s := &Storage{}
	s.Data = make(map[string][]byte)
	return s
}

type Storage struct {
	Data map[string][]byte
	sync.RWMutex
}

func (s *Storage) Set(key string, value []byte) {
	s.Lock()
	s.Data[key] = value
	s.Unlock()
}

func (s *Storage) Get(key string) []byte {
	var data []byte

	s.Lock()
	data = s.Data[key]
	s.Unlock()

	return data
}

func (s *Storage) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
