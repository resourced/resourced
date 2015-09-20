// +build linux

package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	"github.com/resourced/resourced/readers"
	gopsutil_cpu "github.com/shirou/gopsutil/cpu"
	gopsutil_docker "github.com/shirou/gopsutil/docker"
	"strings"
)

func init() {
	readers.Register("DockerContainersCpu", NewDockerContainersCpu)
}

func NewDockerContainersCpu() readers.IReader {
	m := &DockerContainersCpu{}
	m.Data = make(map[string]*gopsutil_cpu.CPUTimesStat)
	return m
}

// DockerContainersCpu gathers docker containers CPU data.
// Data sources:
// * https://github.com/shirou/gopsutil/tree/master/cpu
// * https://github.com/shirou/gopsutil/blob/master/docker/docker_linux.go
type DockerContainersCpu struct {
	Data           map[string]*gopsutil_cpu.CPUTimesStat
	DockerHost     string
	CgroupBasePath string
}

// Run gathers cgroup CPU information from cgroup itself.
// If you use container via systemd.slice, you could use
// containerid = docker-<container id>.scope and base=/sys/fs/cgroup/cpuacct/system.slice/
func (m *DockerContainersCpu) Run() error {
	containers, err := libdocker.AllContainers(m.DockerHost)
	if err != nil {
		return nil
	}

	// Check if using systemd.
	useSystemd := false
	if strings.Contains(m.CgroupBasePath, "/system.slice") {
		useSystemd = true
	}

	for _, container := range containers {
		if container.ID != "" {
			containerDir := container.ID
			if useSystemd {
				containerDir = "docker-" + container.ID + ".scope"
			}

			data, err := gopsutil_docker.CgroupCPU(containerDir, m.CgroupBasePath)
			if err == nil && len(container.Names) > 0 {
				m.Data[container.Names[0]] = data
			}
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersCpu) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
