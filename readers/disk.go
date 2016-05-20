package readers

import (
	"encoding/json"

	gopsutil_disk "github.com/shirou/gopsutil/disk"
)

func init() {
	Register("DiskPartitions", NewDiskPartitions)
	Register("DiskIO", NewDiskIO)
}

// NewDiskPartitions is DiskPartitions constructor.
func NewDiskPartitions() IReader {
	d := &DiskPartitions{}
	d.Data = make(map[string]map[string]gopsutil_disk.PartitionStat)
	return d
}

// DiskPartitions is a reader that gathers partition data.
// Data source: https://github.com/shirou/gopsutil/tree/master/disk
type DiskPartitions struct {
	Data map[string]map[string]gopsutil_disk.PartitionStat
}

// Run gathers partition information from gopsutil.
func (d *DiskPartitions) Run() error {
	dataSlice, err := gopsutil_disk.Partitions(true)
	if err != nil {
		return err
	}

	d.Data["PartitionsByDevice"] = make(map[string]gopsutil_disk.PartitionStat)
	d.Data["PartitionsByMount"] = make(map[string]gopsutil_disk.PartitionStat)

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
func NewDiskIO() IReader {
	d := &DiskIO{}
	d.Data = make(map[string]gopsutil_disk.IOCountersStat)
	return d
}

// DiskIO is a reader that gathers disk io data.
// Data source: https://github.com/shirou/gopsutil/tree/master/disk
type DiskIO struct {
	Data map[string]gopsutil_disk.IOCountersStat
}

// Run gathers disk IO information from gopsutil.
func (d *DiskIO) Run() error {
	data, err := gopsutil_disk.IOCounters()
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
