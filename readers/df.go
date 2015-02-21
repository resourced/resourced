package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	"strings"
)

// NewDf is Df constructor.
func NewDf() *Df {
	d := &Df{}
	d.Data = make(map[string]map[string]interface{})
	return d
}

// Df is a reader that scrapes disk free data and presents it in the form similar to `df`.
// Data source: https://github.com/cloudfoundry/gosigar/tree/master
type Df struct {
	Data    map[string]map[string]interface{}
	FSPaths string
}

func (d *Df) runDefault() error {
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

func (d *Df) runCustomPaths() error {
	for _, path := range strings.Split(d.FSPaths, ",") {
		path = strings.TrimSpace(path)

		usage := sigar.FileSystemUsage{}
		err := usage.Get(path)

		if err == nil {
			d.Data[path] = make(map[string]interface{})
			d.Data[path]["Total"] = usage.Total
			d.Data[path]["Available"] = usage.Avail
			d.Data[path]["Used"] = usage.Used
			d.Data[path]["UsePercent"] = usage.UsePercent()
		}
	}
	return nil
}

// Run gathers df information.
func (d *Df) Run() error {
	err := d.runDefault()
	if err != nil {
		return err
	}

	if d.FSPaths != "" {
		err = d.runCustomPaths()
	}

	return err
}

// ToJson serialize Data field to JSON.
func (d *Df) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
