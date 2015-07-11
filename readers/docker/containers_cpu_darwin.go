// +build darwin

package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("DockerContainersCpu", NewDockerContainersCpu)
}

func NewDockerContainersCpu() readers.IReader {
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

// ToJson serialize Data field to JSON.
func (m *DockerContainersCpu) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
