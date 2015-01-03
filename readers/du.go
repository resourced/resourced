package readers

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	gopsutil_disk "github.com/shirou/gopsutil/disk"
)

func NewDu() *Du {
	d := &Du{}
	d.Data = make(map[string]map[string]interface{})
	return d
}

type Du struct {
	Base
	Data map[string]map[string]interface{}
}

func (d *Du) Run() error {
	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		return err
	}

	for _, fs := range fslist.List {
		duStat, err := gopsutil_disk.DiskUsage(fs.DirName)
		if err == nil {
			d.Data[fs.DirName] = make(map[string]interface{})
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

func (d *Du) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
