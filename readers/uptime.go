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

// Uptime is a reader that presents its data in the form similar to `uptime`.
// Data source: https://github.com/cloudfoundry/gosigar/tree/master
type Uptime struct {
	Data map[string]interface{}
}

// Run gathers uptime information from gosigar.
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
	u.Data["TimeZone"] = currentTime.Format("MST")

	return err
}

// ToJson serialize Data field to JSON.
func (u *Uptime) ToJson() ([]byte, error) {
	return json.Marshal(u.Data)
}
