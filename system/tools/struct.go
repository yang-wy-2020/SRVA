package tools

import (
	"github.com/yang-wy-2020/SRVA/system/cpu"
	"github.com/yang-wy-2020/SRVA/system/disk"
	"github.com/yang-wy-2020/SRVA/system/io"
	"github.com/yang-wy-2020/SRVA/system/memory"
	"github.com/yang-wy-2020/SRVA/system/network"
	"github.com/yang-wy-2020/SRVA/system/sar"
	"github.com/yang-wy-2020/SRVA/system/time"
)

const (
	journal                   string = "journalctl"
	chartsWidth, chartsHeight string = "1300px", "650px"
	effectiveNumber           int    = 200
)

type System struct {
	CPU     cpu.CpuInformation
	IO      io.IoInformation
	SYSTEM  sar.SystemInformation
	NETWORK network.NetworkInformation
	DISK    disk.DiskInformation
	MEMORY  memory.MemoryInformation
	TIME    time.TimeCheckInformation
	NOTE    string
}

type Data struct {
	StartTime      string
	EndTime        string
	CpuProcess     []string
	NetworkCard    []string
	MonitorService map[string]string
	SavePath       string
	OutputPath     string
}
