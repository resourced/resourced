package storage

import (
	"sync"
)

func NewStorage() *Storage {
	s := &Storage{}
	s.data = make(map[string][]byte)
	return s
}

type Storage struct {
	data map[string][]byte
	sync.RWMutex
}

func (s *Storage) Set(key string, value []byte) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s *Storage) Get(key string) []byte {
	var data []byte

	s.Lock()
	data = s.data[key]
	s.Unlock()

	return data
}
