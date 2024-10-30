package memory

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetMemInformation(file string, counts int) MemoryInformation {
	var r MemoryInformation
	getMemValue(&r, file, counts)
	return r
}

func strToFloat64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error getting memory info", err)
	}
	return num
}

func strToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("Error getting memory vsz or rss", err)
	}
	return num
}

func getMemValue(r *MemoryInformation, file string, counts int) {
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
	mem_re := regexp.MustCompile(`.*mem_monitor.*sh.*`)
	mem_str := "mem_monitor.sh"
	for i := 0; i < counts; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, mem_str) {
			mem_r := mem_re.FindStringSubmatch(line)
			time := strings.Split(mem_r[0], currentUser)[0]
			r.Time = append(r.Time, time)
			mem_info := strings.Fields(mem_r[0])
			mem_vsz := strToInt64(mem_info[9])
			mem_rss := strToInt64(mem_info[10])
			mem_per := strToFloat64(mem_info[11])
			mem_cmd := mem_info[12]
			mem_pid := mem_info[6]
			r.ProcessInfo = append(r.ProcessInfo, ProcessInformation{
				MemoryPercentage:  mem_per,
				MemoryCmd:         mem_cmd,
				VirtualMemorySize: int(mem_vsz),
				ResidentSetSize:   int(mem_rss),
				MemoryPid:         mem_pid,
			})
		}
	}
}
