// +build linux

package docker

import (
	"encoding/json"
	"github.com/resourced/resourced/libdocker"
	"github.com/resourced/resourced/readers"
	gopsutil_docker "github.com/shirou/gopsutil/docker"
	"strings"
)

func init() {
	readers.Register("DockerContainersMemory", NewDockerContainersMemory)
}

func NewDockerContainersMemory() readers.IReader {
	m := &DockerContainersMemory{}
	m.Data = make(map[string]*gopsutil_docker.CgroupMemStat)
	return m
}

// DockerContainersMemory gathers Docker memory data.
// Data source: https://github.com/shirou/gopsutil/blob/master/docker/docker_linux.go
type DockerContainersMemory struct {
	Data           map[string]*gopsutil_docker.CgroupMemStat
	DockerHost     string
	CgroupBasePath string
}

// Run gathers cgroup memory information from cgroup itself.
// If you use container via systemd.slice, you could use
// containerid = docker-<container id>.scope and base=/sys/fs/cgroup/memory/system.slice/
func (m *DockerContainersMemory) Run() error {
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

			data, err := gopsutil_docker.CgroupMem(containerDir, m.CgroupBasePath)
			if err == nil && len(container.Names) > 0 {
				m.Data[container.Names[0]] = data
			}
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersMemory) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
