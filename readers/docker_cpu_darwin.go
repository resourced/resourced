// +build darwin

package readers

import (
	"encoding/json"
)

func NewDockerContainersCpu() *DockerContainersCpu {
	m := &DockerContainersCpu{}
	m.Data = make(map[string]string)
	return m
}

type DockerContainersCpu struct {
	Data           map[string]string
	CgroupBasePath string
}

func (m *DockerContainersCpu) Run() error {
	m.Data["Error"] = "Docker cgroup memory data is only available on Linux."
	return nil
}

func (m *DockerContainersCpu) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
