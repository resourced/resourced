package readers

import (
	"encoding/json"
	gopsutil_net "github.com/resourced/resourced/vendor/gopsutil/net"
)

func NewNetIO() *NetIO {
	n := &NetIO{}
	n.Data = make(map[string]gopsutil_net.NetIOCountersStat)
	return n
}

type NetIO struct {
	Data map[string]gopsutil_net.NetIOCountersStat
}

// Run gathers network IO data from gopsutil.
func (n *NetIO) Run() error {
	dataSlice, err := gopsutil_net.NetIOCounters(true)
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

func NewNetInterfaces() *NetInterfaces {
	n := &NetInterfaces{}
	n.Data = make(map[string]gopsutil_net.NetInterfaceStat)
	return n
}

type NetInterfaces struct {
	Data map[string]gopsutil_net.NetInterfaceStat
}

// Run gathers network interfaces data from gopsutil.
func (n *NetInterfaces) Run() error {
	dataSlice, err := gopsutil_net.NetInterfaces()
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
