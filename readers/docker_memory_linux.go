// +build linux

package readers

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	gopsutil_docker "github.com/shirou/gopsutil/docker"
)

func NewDockerMemory() *DockerMemory {
	m := &DockerMemory{}
	m.Data = make(map[string]*gopsutil_docker.CgroupMemStat)
	return m
}

type DockerMemory struct {
	Base
	Data map[string]*gopsutil_docker.CgroupMemStat
}

func (m *DockerMemory) Run() error {
	containers, err := libdocker.AllContainers()
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" {
			data, err := gopsutil_docker.CgroupMemDocker(container.ID)
			if err == nil {
				m.Data[container.ID] = data
			}
		}
	}

	return nil
}

func (m *DockerMemory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
