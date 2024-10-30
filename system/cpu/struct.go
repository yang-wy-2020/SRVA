package cpu

const (
	LoadAvg_g string = "Load Avg"
)

type LoadAvgInformation struct {
	AvgOne, AvgFive, AvgFifteen float64
}

type ProcessInformation struct {
	Name                       string
	CpuUse, UserUse, SystemUse []float64
}
type CpuInformation struct {
	Time        []string
	LoadAvg     []LoadAvgInformation
	ProcessInfo []ProcessInformation
	Rate        []float64
	CpuCore     []uint8
}
