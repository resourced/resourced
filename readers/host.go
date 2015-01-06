package readers

import (
	"encoding/json"
	gopsutil_host "github.com/shirou/gopsutil/host"
)

func NewHostInfo() *HostInfo {
	h := &HostInfo{}
	h.Data = make(map[string]interface{})
	return h
}

type HostInfo struct {
	Data map[string]interface{}
}

func (h *HostInfo) Run() error {
	data, err := gopsutil_host.HostInfo()
	if err != nil {
		return err
	}

	h.Data["Hostname"] = data.Hostname
	h.Data["Uptime"] = data.Uptime
	h.Data["Procs"] = data.Procs
	h.Data["OS"] = data.OS
	h.Data["Platform"] = data.Platform
	h.Data["PlatformFamily"] = data.PlatformFamily
	h.Data["PlatformVersion"] = data.PlatformVersion
	h.Data["VirtualizationSystem"] = data.VirtualizationSystem
	h.Data["VirtualizationRole"] = data.VirtualizationRole

	bootTime, err := gopsutil_host.BootTime()
	if err == nil {
		h.Data["BootTime"] = bootTime
	}

	return nil
}

func (h *HostInfo) ToJson() ([]byte, error) {
	return json.Marshal(h.Data)
}

// ----------------------------------------------------------------

func NewHostUsers() *HostUsers {
	h := &HostUsers{}
	h.Data = make(map[string]gopsutil_host.UserStat)
	return h
}

type HostUsers struct {
	Data map[string]gopsutil_host.UserStat
}

func (h *HostUsers) Run() error {
	dataSlice, err := gopsutil_host.Users()
	if err != nil {
		return err
	}

	for _, data := range dataSlice {
		h.Data[data.User] = data
	}
	return nil
}

func (h *HostUsers) ToJson() ([]byte, error) {
	return json.Marshal(h.Data)
}
