package readers

import (
	"encoding/json"
	"net"
)

func NewNetworkInterfaces() *NetworkInterfaces {
	n := &NetworkInterfaces{}
	n.Data = make(map[string]map[string][]string)
	return n
}

type NetworkInterfaces struct {
	Base
	Data map[string]map[string][]string
}

func (n *NetworkInterfaces) Run() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err == nil {
			n.Data[iface.Name] = make(map[string][]string)
			n.Data[iface.Name]["Addresses"] = make([]string, len(addrs))

			for i, addr := range addrs {
				n.Data[iface.Name]["Addresses"][i] = addr.String()
			}
		}
	}
	return nil
}

func (n *NetworkInterfaces) ToJson() ([]byte, error) {
	return json.Marshal(n.Data)
}
