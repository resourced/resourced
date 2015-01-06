package readers

import (
	"encoding/json"
	gopsutil_disk "github.com/shirou/gopsutil/disk"
)

func NewDiskPartitions() *DiskPartitions {
	d := &DiskPartitions{}
	d.Data = make(map[string]map[string]gopsutil_disk.DiskPartitionStat)
	return d
}

type DiskPartitions struct {
	Data map[string]map[string]gopsutil_disk.DiskPartitionStat
}

func (d *DiskPartitions) Run() error {
	dataSlice, err := gopsutil_disk.DiskPartitions(true)
	if err != nil {
		return err
	}

	d.Data["PartitionsByDevice"] = make(map[string]gopsutil_disk.DiskPartitionStat)
	d.Data["PartitionsByMount"] = make(map[string]gopsutil_disk.DiskPartitionStat)

	for _, data := range dataSlice {
		d.Data["PartitionsByDevice"][data.Device] = data
		d.Data["PartitionsByMount"][data.Mountpoint] = data
	}

	return nil
}

func (d *DiskPartitions) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}

// ----------------------------------------------------------------
