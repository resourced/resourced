package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
)

func NewDf() *Df {
	n := &Df{}
	n.Data = make(map[string]map[string]interface{})
	return n
}

type Df struct {
	Base
	Data map[string]map[string]interface{}
}

func (d *Df) Run() error {
	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		return err
	}

	for _, fs := range fslist.List {
		usage := sigar.FileSystemUsage{}
		err := usage.Get(fs.DirName)

		if err == nil {
			d.Data[fs.DirName] = make(map[string]interface{})
			d.Data[fs.DirName]["DeviceName"] = fs.DevName
			d.Data[fs.DirName]["Total"] = usage.Total
			d.Data[fs.DirName]["Available"] = usage.Avail
			d.Data[fs.DirName]["Used"] = usage.Used
			d.Data[fs.DirName]["UsePercent"] = usage.UsePercent()
		}
	}
	return nil
}

func (d *Df) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
