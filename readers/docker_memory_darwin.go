// +build darwin

package readers

import (
	"encoding/json"
)

func NewDockerMemory() *DockerMemory {
	m := &DockerMemory{}
	m.Data = make(map[string]string)
	return m
}

type DockerMemory struct {
	Base
	Data map[string]string
}

func (m *DockerMemory) Run() error {
	m.Data["Error"] = "Docker cgroup memory data is only available on Linux."
	return nil
}

func (m *DockerMemory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
