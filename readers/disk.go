package readers

import (
	"encoding/json"
	gopsutil_disk "github.com/resourced/resourced/vendor/gopsutil/disk"
)

// NewDiskPartitions is DiskPartitions constructor.
func NewDiskPartitions() *DiskPartitions {
	d := &DiskPartitions{}
	d.Data = make(map[string]map[string]gopsutil_disk.DiskPartitionStat)
	return d
}

// DiskPartitions is a reader that gathers partition data.
// Data source: https://github.com/shirou/gopsutil/tree/master/disk
type DiskPartitions struct {
	Data map[string]map[string]gopsutil_disk.DiskPartitionStat
}

// Run gathers partition information from gopsutil.
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

// ToJson serialize Data field to JSON.
func (d *DiskPartitions) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}

// ----------------------------------------------------------------

// NewDiskIO is DiskIO constructor.
func NewDiskIO() *DiskIO {
	d := &DiskIO{}
	d.Data = make(map[string]gopsutil_disk.DiskIOCountersStat)
	return d
}

// DiskIO is a reader that gathers disk io data.
// Data source: https://github.com/shirou/gopsutil/tree/master/disk
type DiskIO struct {
	Data map[string]gopsutil_disk.DiskIOCountersStat
}

// Run gathers disk IO information from gopsutil.
func (d *DiskIO) Run() error {
	data, err := gopsutil_disk.DiskIOCounters()
	if err != nil {
		return err
	}

	d.Data = data
	return nil
}

// ToJson serialize Data field to JSON.
func (d *DiskIO) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
