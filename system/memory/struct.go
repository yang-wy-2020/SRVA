package memory

const (
	Service_g string = "qomolo_mem_monitor.service"
)

type ProcessInformation struct {
	MemoryPercentage                   float64
	VirtualMemorySize, ResidentSetSize int
	MemoryCmd                          string
	MemoryPid                          string
}

type MemoryInformation struct {
	Time        []string
	ProcessInfo []ProcessInformation
}
