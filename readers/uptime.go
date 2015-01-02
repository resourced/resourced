package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	"time"
)

func NewUptime() *Uptime {
	u := &Uptime{}
	u.Data = make(map[string]interface{})
	return u
}

type Uptime struct {
	Base
	Data map[string]interface{}
}

func (u *Uptime) Run() error {
	loadAvg := NewLoadAvg()
	err := loadAvg.Run()
	if err != nil {
		return err
	}

	u.Data = loadAvg.Data

	uptime := sigar.Uptime{}
	err = uptime.Get()
	if err != nil {
		return err
	}

	currentTime := time.Now()

	u.Data["CurrentTimeUnixNano"] = currentTime.UnixNano()
	u.Data["CurrentTime"] = currentTime.Format("15:04:05")
	u.Data["Uptime"] = uptime.Format()

	return err
}

func (u *Uptime) ToJson() ([]byte, error) {
	return json.Marshal(u.Data)
}
