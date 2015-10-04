package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
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

// Set stores bytes in-memory with gzip lvl 9 compression.
func (s *Storage) Set(key string, value []byte) {
	var b bytes.Buffer

	w, _ := gzip.NewWriterLevel(&b, 9)
	w.Write(value)
	w.Flush()
	w.Close()

	s.Lock()
	s.Data[key] = b.Bytes()
	s.Unlock()
}

// Get returns bytes after un-gzipping.
func (s *Storage) Get(key string) []byte {
	var compressedData []byte

	s.RLock()
	compressedData = s.Data[key]
	s.RUnlock()

	r, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	return data
}

func (s *Storage) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
