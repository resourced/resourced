package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
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

func NewResourcedMasterMetadataStorage(root string) *ResourcedMasterMetadataStorage {
	s := &ResourcedMasterMetadataStorage{}
	s.Root = root
	return s
}

type ResourcedMasterMetadataStorage struct {
	Root string
}

func (s *ResourcedMasterMetadataStorage) Set(key string, data []byte) error {
	req, err := http.NewRequest("POST", path.Join(s.Root, key), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return nil
}

func (s *ResourcedMasterMetadataStorage) Get(key string) ([]byte, error) {
	resp, err := http.Get(path.Join(s.Root, key))
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return ioutil.ReadAll(resp.Body)
}
