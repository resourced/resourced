// +build linux

package docker

import (
	"encoding/json"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/resourced/resourced/libdocker"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("DockerContainersNetDev", NewDockerContainersNetDev)
}

func NewDockerContainersNetDev() readers.IReader {
	m := &DockerContainersNetDev{}
	m.Data = make(map[string]map[string]linuxproc.NetworkStat)
	return m
}

// DockerContainersNetDev gathers Docker memory data.
// Data source: https://github.com/shirou/gopsutil/blob/master/docker/docker_linux.go
type DockerContainersNetDev struct {
	Data           map[string]map[string]linuxproc.NetworkStat
	DockerHost     string
	CgroupBasePath string
}

// Run gathers cgroup memory information from cgroup itself.
// If you use container via systemd.slice, you could use
// containerid = docker-<container id>.scope and base=/sys/fs/cgroup/memory/system.slice/
func (m *DockerContainersNetDev) Run() error {
	containers, err := libdocker.AllInspectedContainers(m.DockerHost)
	if err != nil {
		return nil
	}

	for _, container := range containers {
		if container.ID != "" && container.State.Running {
			pid := container.State.Pid

			data, err := linuxproc.ReadNetworkStat(fmt.Sprintf("/proc/%v/net/dev", pid))
			if err == nil {
				m.Data[container.NiceImageName+"-"+container.ID] = make(map[string]linuxproc.NetworkStat)

				for _, perIface := range data {
					m.Data[container.NiceImageName+"-"+container.ID][perIface.Iface] = perIface
				}
			}
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *DockerContainersNetDev) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
