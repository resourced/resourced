package readers

import (
	"errors"
	readers_procfs "github.com/resourced/resourced/readers/procfs"
)

// NewGoStruct instantiates IReader
func NewGoStruct(name string) (IReader, error) {
	var structInstance IReader

	if name == "CpuInfo" {
		structInstance = NewCpuInfo()
	}
	if name == "Df" {
		structInstance = NewDf()
	}
	if name == "Du" {
		structInstance = NewDu()
	}
	if name == "DiskPartitions" {
		structInstance = NewDiskPartitions()
	}
	if name == "DiskIO" {
		structInstance = NewDiskIO()
	}
	if name == "DockerContainersMemory" {
		structInstance = NewDockerContainersMemory()
	}
	if name == "DockerContainersCpu" {
		structInstance = NewDockerContainersCpu()
	}
	if name == "Free" {
		structInstance = NewFree()
	}
	if name == "HostInfo" {
		structInstance = NewHostInfo()
	}
	if name == "HostUsers" {
		structInstance = NewHostUsers()
	}
	if name == "LoadAvg" {
		structInstance = NewLoadAvg()
	}
	if name == "Ps" {
		structInstance = NewPs()
	}
	if name == "NetIO" {
		structInstance = NewNetIO()
	}
	if name == "NetInterfaces" {
		structInstance = NewNetInterfaces()
	}
	if name == "ProcCpuInfo" {
		structInstance = readers_procfs.NewProcCpuInfo()
	}
	if name == "ProcDiskStats" {
		structInstance = readers_procfs.NewProcDiskStats()
	}
	if name == "ProcLoadAvg" {
		structInstance = readers_procfs.NewProcLoadAvg()
	}
	if name == "ProcMemInfo" {
		structInstance = readers_procfs.NewProcMemInfo()
	}
	if name == "ProcMounts" {
		structInstance = readers_procfs.NewProcMounts()
	}
	if name == "ProcNetDev" {
		structInstance = readers_procfs.NewProcNetDev()
	}
	if name == "ProcStat" {
		structInstance = readers_procfs.NewProcStat()
	}
	if name == "ProcUptime" {
		structInstance = readers_procfs.NewProcUptime()
	}
	if name == "ProcVmStat" {
		structInstance = readers_procfs.NewProcVmStat()
	}
	if name == "Uptime" {
		structInstance = NewUptime()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

type IReader interface {
	Run() error
	ToJson() ([]byte, error)
}
