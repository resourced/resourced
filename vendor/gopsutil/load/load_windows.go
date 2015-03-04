// +build windows

package load

import (
	common "github.com/resourced/resourced/vendor/gopsutil/common"
)

func LoadAvg() (*LoadAvgStat, error) {
	ret := LoadAvgStat{}

	return &ret, common.NotImplementedError
}
