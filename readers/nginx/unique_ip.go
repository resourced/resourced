package redis

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("NginxUniqueIP", NewNginxUniqueIP)
}

func NewNginxUniqueIP() readers.IReader {
	nginx := &NginxUniqueIP{}
	nginx.Data = make(map[string]int64)

	return nginx
}

type NginxUniqueIP struct {
	readers.TailFile
	Data map[string]int64
	sync.RWMutex
}

func (nginx *NginxUniqueIP) Run() error {
	tailer, err := nginx.Tailer()
	if err != nil {
		return err
	}

	go func() {
		for line := range tailer.Lines {
			lineChunks := strings.Fields(line.Text)
			ip := lineChunks[0]

			nginx.Lock()
			_, ok := nginx.Data[ip]
			if !ok {
				nginx.Data[ip] = int64(0)
			}
			nginx.Data[ip] = nginx.Data[ip] + 1
			nginx.Unlock()

			jsn, _ := nginx.ToJson()
			println(string(jsn))
		}
	}()

	jsn, _ := nginx.ToJson()
	println(string(jsn))

	return nil
}

// ToJson serialize Data field to JSON.
func (nginx *NginxUniqueIP) ToJson() ([]byte, error) {
	return json.Marshal(nginx.Data)
}
