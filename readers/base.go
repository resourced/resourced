// Package readers provides objects that gathers resource data from a host.
package readers

import (
	"errors"
	resourced_config "github.com/resourced/resourced/config"
	readers_docker "github.com/resourced/resourced/readers/docker"
	readers_mysql "github.com/resourced/resourced/readers/mysql"
	readers_procfs "github.com/resourced/resourced/readers/procfs"
	readers_redis "github.com/resourced/resourced/readers/redis"
	"reflect"
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
	if name == "DockerContainers" {
		structInstance = readers_docker.NewDockerContainers()
	}
	if name == "DockerImages" {
		structInstance = readers_docker.NewDockerImages()
	}
	if name == "DockerContainersMemory" {
		structInstance = readers_docker.NewDockerContainersMemory()
	}
	if name == "DockerContainersCpu" {
		structInstance = readers_docker.NewDockerContainersCpu()
	}
	if name == "DockerContainersNetDev" {
		structInstance = readers_docker.NewDockerContainersNetDev()
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
	if name == "MysqlProcesslist" {
		structInstance = readers_mysql.NewMysqlProcesslist()
	}
	if name == "MysqlInformationSchemaTables" {
		structInstance = readers_mysql.NewMysqlInformationSchemaTables()
	}
	if name == "MysqlDumpSlow" {
		structInstance = readers_mysql.NewMysqlDumpSlow()
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
	if name == "ProcNetDevPid" {
		structInstance = readers_procfs.NewProcNetDevPid()
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
	if name == "RedisInfo" {
		structInstance = readers_redis.NewRedisInfo()
	}
	if name == "Ps" {
		structInstance = NewPs()
	}
	if name == "Uptime" {
		structInstance = NewUptime()
	}
	if name == "DMI" {
		structInstance = NewDMI()
	}

	if structInstance == nil {
		return nil, errors.New("GoStruct is undefined.")
	}

	return structInstance, nil
}

// NewGoStructByConfig instantiates IReader given Config struct
func NewGoStructByConfig(config resourced_config.Config) (IReader, error) {
	reader, err := NewGoStruct(config.GoStruct)
	if err != nil {
		return nil, err
	}

	// Populate IReader fields dynamically
	if len(config.GoStructFields) > 0 {
		for structFieldInString, value := range config.GoStructFields {
			goStructField := reflect.ValueOf(reader).Elem().FieldByName(structFieldInString)

			if goStructField.IsValid() && goStructField.CanSet() {
				valueOfValue := reflect.ValueOf(value)
				goStructField.Set(valueOfValue)
			}
		}
	}

	return reader, err
}

// IReader is generic interface for all readers.
type IReader interface {
	Run() error
	ToJson() ([]byte, error)
}
