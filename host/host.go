package host

import (
	"os"
)

func NewHostByHostname() (*Host, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	h := NewHost(hostname)
	return h, nil
}

func NewHost(name string) *Host {
	h := &Host{}
	h.Name = name

	return h
}

type Host struct {
	Name              string
	Tags              []string
	NetworkInterfaces map[string]map[string]interface{}
}
