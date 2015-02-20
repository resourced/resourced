// +build linux

package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	gopsutil_docker "github.com/shirou/gopsutil/docker"
)

func NewDockerContainersMemory() *DockerContainersMemory {
	m := &DockerContainersMemory{}
	m.Data = make(map[string]*gopsutil_docker.CgroupMemStat)
	return m
}

// DockerContainersMemory gathers Docker memory data.
// Data source: https://github.com/shirou/gopsutil/blob/master/docker/docker_linux.go
type DockerContainersMemory struct {
	Data           map[string]*gopsutil_docker.CgroupMemStat
	CgroupBasePath string
}

// Run gathers cgroup memory information from cgroup itself.
func (m *DockerContainersMemory) Run() error {
	containers, err := libdocker.AllContainers("")
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" {
			data, err := gopsutil_docker.CgroupMem(container.ID, m.CgroupBasePath)
			if err == nil {
				m.Data[container.ID] = data
			}
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersMemory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
