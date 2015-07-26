package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

type MetadataStorages struct {
	ResourcedMaster *ResourcedMasterMetadataStorage
}

func NewResourcedMasterMetadataStorage(root, accessToken string) *ResourcedMasterMetadataStorage {
	s := &ResourcedMasterMetadataStorage{}
	s.Root = root
	s.AccessToken = accessToken
	return s
}

type ResourcedMasterMetadataStorage struct {
	Root        string
	AccessToken string
}

func (s *ResourcedMasterMetadataStorage) Set(key string, data []byte) error {
	if strings.HasPrefix(key, "/") {
		strings.Replace(key, "/", "", 1)
	}

	postPath = "api/metadata"

	url := strings.Join([]string{s.Root, postPath, key}, "/")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(s.AccessToken, "")

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
	if strings.HasPrefix(key, "/") {
		strings.Replace(key, "/", "", 1)
	}

	getPath = "api/metadata"

	url := strings.Join([]string{s.Root, postPath, key}, "/")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return ioutil.ReadAll(resp.Body)
}
