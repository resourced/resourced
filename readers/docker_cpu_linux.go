// +build linux

package readers

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	gopsutil_cpu "github.com/shirou/gopsutil/cpu"
	gopsutil_docker "github.com/shirou/gopsutil/docker"
)

func NewDockerContainersCpu() *DockerContainersCpu {
	m := &DockerContainersCpu{}
	m.Data = make(map[string]*gopsutil_cpu.CPUTimesStat)
	return m
}

// DockerContainersCpu gathers docker containers CPU data.
// Data sources:
// * https://github.com/shirou/gopsutil/tree/master/cpu
// * https://github.com/shirou/gopsutil/tree/master/docker
type DockerContainersCpu struct {
	Data map[string]*gopsutil_cpu.CPUTimesStat
}

func (m *DockerContainersCpu) Run() error {
	containers, err := libdocker.AllContainers()
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" {
			data, err := gopsutil_docker.CgroupCPUDocker(container.ID)
			if err == nil {
				m.Data[container.ID] = data
			}
		}
	}

	return nil
}

func (m *DockerContainersCpu) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
