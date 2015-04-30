// +build darwin

package docker

import (
	"encoding/json"
)

func NewDockerContainersNetDev() *DockerContainersNetDev {
	m := &DockerContainersNetDev{}
	m.Data = make(map[string]string)
	return m
}

type DockerContainersNetDev struct {
	Data           map[string]string
	DockerHost     string
	CgroupBasePath string
}

func (m *DockerContainersNetDev) Run() error {
	m.Data["Error"] = "Docker pid/net/dev data is only available on Linux."
	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersNetDev) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
