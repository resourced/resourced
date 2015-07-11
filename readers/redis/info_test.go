package redis

import (
	"testing"
)

func TestRedisInfoRun(t *testing.T) {
	r := &RedisInfo{}
	r.Data = make(map[string]string)
	if r.initConnection() == nil {
		err := r.Run()
		if err != nil {
			t.Errorf("fetching INFO data should always be successful. Error: %v", err)
		}

		if len(r.Data) == 0 {
			inJson, _ := r.ToJson()
			t.Errorf("fetching INFO data should always be successful. Data: %v", string(inJson))
		}
	}
}
