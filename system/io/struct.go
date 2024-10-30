package io

const (
	Service_g string = "qomolo_io_monitor.service"
)

type ProcessInformation struct {
	WriteKbSec, ReadKbSec float64
	IoCmd                 string
	IoPid                 int
}

type IoInformation struct {
	Time        []string
	ProcessInfo []ProcessInformation
}
