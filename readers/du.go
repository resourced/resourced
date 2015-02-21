package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	gopsutil_disk "github.com/shirou/gopsutil/disk"
	"strings"
)

func NewDu() *Du {
	d := &Du{}
	d.Data = make(map[string]map[string]interface{})
	return d
}

// Df is a reader that scrapes disk usage data and presents it in the form similar to `du`.
// Data sources:
// * https://github.com/cloudfoundry/gosigar/tree/master
// * https://github.com/shirou/gopsutil/tree/master/disk
type Du struct {
	Data    map[string]map[string]interface{}
	FSPaths string
}

// Run gathers du information from gosigar.
func (d *Du) runDefault() error {
	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		return err
	}

	for _, fs := range fslist.List {
		duStat, err := gopsutil_disk.DiskUsage(fs.DirName)
		if err == nil {
			d.Data[fs.DirName] = make(map[string]interface{})
			d.Data[fs.DirName]["DeviceName"] = fs.DevName
			d.Data[fs.DirName]["Path"] = duStat.Path
			d.Data[fs.DirName]["Total"] = duStat.Total
			d.Data[fs.DirName]["Free"] = duStat.Free
			d.Data[fs.DirName]["InodesTotal"] = duStat.InodesTotal
			d.Data[fs.DirName]["InodesFree"] = duStat.InodesFree
			d.Data[fs.DirName]["InodesUsed"] = duStat.InodesUsed
			d.Data[fs.DirName]["Used"] = duStat.Used

			if duStat.InodesTotal != 0 {
				d.Data[fs.DirName]["InodesUsedPercent"] = duStat.InodesUsedPercent
			}

			if duStat.Total != 0 {
				d.Data[fs.DirName]["UsedPercent"] = duStat.UsedPercent
			}
		}
	}
	return nil
}

func (d *Du) runCustomPaths() error {
	for _, path := range strings.Split(d.FSPaths, ",") {
		path = strings.TrimSpace(path)

		duStat, err := gopsutil_disk.DiskUsage(path)
		if err == nil {
			d.Data[path] = make(map[string]interface{})
			d.Data[path]["Path"] = duStat.Path
			d.Data[path]["Total"] = duStat.Total
			d.Data[path]["Free"] = duStat.Free
			d.Data[path]["InodesTotal"] = duStat.InodesTotal
			d.Data[path]["InodesFree"] = duStat.InodesFree
			d.Data[path]["InodesUsed"] = duStat.InodesUsed
			d.Data[path]["Used"] = duStat.Used

			if duStat.InodesTotal != 0 {
				d.Data[path]["InodesUsedPercent"] = duStat.InodesUsedPercent
			}

			if duStat.Total != 0 {
				d.Data[path]["UsedPercent"] = duStat.UsedPercent
			}
		}
	}
	return nil
}

// Run gathers du information.
func (d *Du) Run() error {
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
func (d *Du) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
