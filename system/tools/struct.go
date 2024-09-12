package tools

import (
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/cpu"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/io"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/sar"
)

const (
	journal                   string = "journalctl"
	output                    string = "/data/tmp/"
	chartsWidth, chartsHeight string = "1300px", "650px"
)

type System struct {
	CPU    cpu.CpuInformation
	IO     io.IoInformation
	SYSTEM sar.SystemInformation
	NOTE   string
}

type Data struct {
	StartTime      string
	EndTime        string
	ModelsList     []string
	NetworkCard    []string
	MonitorService map[string]string
}
