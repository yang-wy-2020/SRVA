package tools

import (
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/cpu"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/disk"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/io"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/memory"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/network"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/sar"
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
