// +build darwin

package docker

import (
	"encoding/json"
)

func NewDockerContainersMemory() *DockerContainersMemory {
	m := &DockerContainersMemory{}
	m.Data = make(map[string]string)
	return m
}

type DockerContainersMemory struct {
	Data           map[string]string
	DockerHost     string
	CgroupBasePath string
}

func (m *DockerContainersMemory) Run() error {
	m.Data["Error"] = "Docker cgroup memory data is only available on Linux."
	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersMemory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
