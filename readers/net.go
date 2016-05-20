package readers

import (
	"encoding/json"

	gopsutil_net "github.com/shirou/gopsutil/net"
)

func init() {
	Register("NetIO", NewNetIO)
	Register("NetInterfaces", NewNetInterfaces)
}

func NewNetIO() IReader {
	n := &NetIO{}
	n.Data = make(map[string]gopsutil_net.IOCountersStat)
	return n
}

type NetIO struct {
	Data map[string]gopsutil_net.IOCountersStat
}

// Run gathers network IO data from gopsutil.
func (n *NetIO) Run() error {
	dataSlice, err := gopsutil_net.IOCounters(true)
	if err != nil {
		return err
	}

	for _, data := range dataSlice {
		n.Data[data.Name] = data
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (n *NetIO) ToJson() ([]byte, error) {
	return json.Marshal(n.Data)
}

// ------------------------------------------------------

func NewNetInterfaces() IReader {
	n := &NetInterfaces{}
	n.Data = make(map[string]gopsutil_net.InterfaceStat)
	return n
}

type NetInterfaces struct {
	Data map[string]gopsutil_net.InterfaceStat
}

// Run gathers network interfaces data from gopsutil.
func (n *NetInterfaces) Run() error {
	dataSlice, err := gopsutil_net.Interfaces()
	if err != nil {
		return err
	}

	for _, data := range dataSlice {
		n.Data[data.Name] = data
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (n *NetInterfaces) ToJson() ([]byte, error) {
	return json.Marshal(n.Data)
}
