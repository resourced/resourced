package libmap

import (
	"encoding/json"
	"strings"
	"sync"
)

func NewTSafeMapBytes(data map[string][]byte) *TSafeMapBytes {
	mp := &TSafeMapBytes{}
	if data == nil {
		mp.data = make(map[string][]byte)
	} else {
		mp.data = data
	}
	return mp
}

type TSafeMapBytes struct {
	data map[string][]byte
	sync.RWMutex
}

func (mp *TSafeMapBytes) Set(key string, value []byte) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = value
}

func (mp *TSafeMapBytes) Get(key string) []byte {
	mp.Lock()
	defer mp.Unlock()

	return mp.data[key]
}

func (mp *TSafeMapBytes) All() map[string][]byte {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string][]byte)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

func (mp *TSafeMapBytes) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

func NewTSafeMapStrings(data map[string][]string) *TSafeMapStrings {
	mp := &TSafeMapStrings{}
	if data == nil {
		mp.data = make(map[string][]string)
	} else {
		mp.data = data
	}
	return mp
}

type TSafeMapStrings struct {
	data map[string][]string
	sync.RWMutex
}

func (mp *TSafeMapStrings) Set(key string, value []string) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = value
}

func (mp *TSafeMapStrings) Append(key string, value string) {
	mp.Lock()
	defer mp.Unlock()
	mp.data[key] = append(mp.data[key], value)
}

func (mp *TSafeMapStrings) Get(key string) []string {
	mp.Lock()
	defer mp.Unlock()

	original, ok := mp.data[key]
	if !ok {
		return make([]string, 0)
	}

	copydata := make([]string, len(original))
	for i, value := range original {
		copydata[i] = value
	}

	return copydata
}

func (mp *TSafeMapStrings) Reset(key string) {
	mp.Lock()
	defer mp.Unlock()

	mp.data[key] = make([]string, 0)
}

func (mp *TSafeMapStrings) All() map[string][]string {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string][]string)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

func (mp *TSafeMapStrings) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

func NewTSafeMapCounter(data map[string]int) *TSafeMapCounter {
	mp := &TSafeMapCounter{}
	if data == nil {
		mp.data = make(map[string]int)
	} else {
		mp.data = data
	}
	return mp
}

type TSafeMapCounter struct {
	data map[string]int
	sync.RWMutex
}

func (mp *TSafeMapCounter) Incr(key string, value int) {
	mp.Lock()
	mp.data[key] = mp.data[key] + value
	mp.Unlock()
}

func (mp *TSafeMapCounter) Get(key string) int {
	mp.RLock()
	defer mp.RUnlock()

	data, ok := mp.data[key]
	if !ok {
		data = 0
	}

	return data
}

func (mp *TSafeMapCounter) Reset(key string) {
	mp.Lock()
	mp.data[key] = 0
	mp.Unlock()
}

func (mp *TSafeMapCounter) All() map[string]int {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string]int)
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

func (mp *TSafeMapCounter) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}

func NewTSafeNestedMapInterface(data map[string]interface{}) *TSafeNestedMapInterface {
	mp := &TSafeNestedMapInterface{}
	if data == nil {
		mp.data = make(map[string]interface{})
	} else {
		mp.data = data
	}
	return mp
}

type TSafeNestedMapInterface struct {
	data map[string]interface{}
	sync.RWMutex
}

func (mp *TSafeNestedMapInterface) initNestedMap(key string) {
	// Split key by dot, loop deeper into the nesting & create the maps
	keyParts := strings.Split(key, ".")

	mp.Lock()
	m := mp.data

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
	mp.Unlock()
}

func (mp *TSafeNestedMapInterface) Set(key string, value interface{}) {
	mp.initNestedMap(key)

	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	mp.Lock()
	m := mp.data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	m[lastPart] = value
	mp.Unlock()
}

func (mp *TSafeNestedMapInterface) Get(key string) interface{} {
	var data interface{}

	// Split key by dot and loop deeper into the nesting
	keyParts := strings.Split(key, ".")
	lastPart := keyParts[len(keyParts)-1]

	mp.RLock()
	m := mp.data

	for i, keyPart := range keyParts {
		if i == len(keyParts)-1 {
			break
		}

		m = m[keyPart].(map[string]interface{})
	}

	data = m[lastPart]
	mp.RUnlock()

	return data
}

func (mp *TSafeNestedMapInterface) All() map[string]interface{} {
	mp.Lock()
	defer mp.Unlock()

	copydata := make(map[string]interface{})
	for key, value := range mp.data {
		copydata[key] = value
	}

	return copydata
}

func (mp *TSafeNestedMapInterface) ToJson() ([]byte, error) {
	return json.Marshal(mp.data)
}
