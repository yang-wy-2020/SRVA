package disk

const (
	Service_Disk string = "qomolo_disk_monitor.service"
)

type Disk_Info struct {
	DiskDevice         string
	DiskUsed, DiskFree float64
}

type DiskInformation struct {
	Time     []string
	DiskInfo []Disk_Info
}
