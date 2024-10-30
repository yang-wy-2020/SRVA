package sar

var (
	Cpu_flag string = "all"
)

type SystemCpu struct {
	User   float64
	System float64
	IoWait float64
	Idle   float64
}

type SystemInformation struct {
	Time     []string
	CpuTotal []SystemCpu
}
