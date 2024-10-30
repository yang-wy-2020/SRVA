package disk

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetDiskAvail() {
}

func GetDiskInformation(file string, counts int) DiskInformation {
	var r DiskInformation
	getDiskUsedAndFree(&r, file, counts)
	return r
}

func strToFloat(str string) float64 {
	str = str[0 : len(str)-1]
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error getting disk values", err)
	}
	return num
}

func getDiskUsedAndFree(r *DiskInformation, file string, counts int) {
	file_buf, err := os.Open(file)
	if err != nil {
		fmt.Println("open file failed:", err)
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)
	currentUser, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting current user:", err)
	}
	device := regexp.MustCompile(`.*dev.*`)
	device_str := "dev"
	for i := 0; i < counts; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, device_str) {
			disk_r := device.FindStringSubmatch(line)
			time := strings.Split(disk_r[0], currentUser)[0]
			r.Time = append(r.Time, time)
			disk_info := strings.Fields(disk_r[0])
			disk_device := disk_info[5]
			// r.DiskDevice = append(r.DiskDevice, disk_device)
			disk_used := strToFloat(disk_info[6])
			// r.DiskUsed = append(r.DiskUsed, Disk_Used{DiskUsed: disk_used})
			disk_free := strToFloat(disk_info[7])
			// r.DiskFree = append(r.DiskFree, Disk_Free{DiskFree: disk_free})
			r.DiskInfo = append(r.DiskInfo, Disk_Info{
				DiskDevice: disk_device,
				DiskUsed:   disk_used,
				DiskFree:   disk_free})
			// fmt.Println(disk_device, disk_used, disk_free)
		}
	}
}
