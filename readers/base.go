package readers

import (
	"errors"
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
	if name == "Meminfo" {
		structInstance = NewMeminfo()
	}
	if name == "NetIO" {
		structInstance = NewNetIO()
	}
	if name == "NetInterfaces" {
		structInstance = NewNetInterfaces()
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
