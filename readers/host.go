package readers

import (
	"encoding/json"

	gopsutil_host "github.com/shirou/gopsutil/host"
)

func init() {
	Register("HostInfo", NewHostInfo)
	Register("HostUsers", NewHostUsers)
}

func NewHostInfo() IReader {
	h := &HostInfo{}
	h.Data = make(map[string]interface{})
	return h
}

type HostInfo struct {
	Data map[string]interface{}
}

// Run gathers host information from gopsutil.
func (h *HostInfo) Run() error {
	data, err := gopsutil_host.Info()
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

// ToJson serialize Data field to JSON.
func (h *HostInfo) ToJson() ([]byte, error) {
	return json.Marshal(h.Data)
}

// ----------------------------------------------------------------

func NewHostUsers() IReader {
	h := &HostUsers{}
	h.Data = make(map[string]gopsutil_host.UserStat)
	return h
}

type HostUsers struct {
	Data map[string]gopsutil_host.UserStat
}

// Run gathers user information from gopsutil.
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

// ToJson serialize Data field to JSON.
func (h *HostUsers) ToJson() ([]byte, error) {
	return json.Marshal(h.Data)
}
