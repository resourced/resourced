// +build darwin

package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("DockerContainersMemory", NewDockerContainersMemory)
}

func NewDockerContainersMemory() readers.IReader {
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
