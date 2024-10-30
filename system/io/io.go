package io

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetIOInformation(file string, counts int) IoInformation {
	var r IoInformation
	getIOValue(&r, file, counts)
	return r
}

func strToFloat64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error getting io info", err)
	}
	return num
}

func getIOValue(r *IoInformation, file string, counts int) {
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
	io_re := regexp.MustCompile(`.*io_.*monitor.*sh.*`)
	io_str := "io_monitor.sh"
	for i := 0; i < counts; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, io_str) {
			io_r := io_re.FindStringSubmatch(line)
			time := strings.Split(io_r[0], currentUser)[0]
			r.Time = append(r.Time, time)
			io_info := strings.Fields(io_r[0])
			io_rd := strToFloat64(io_info[9])
			io_wr := strToFloat64(io_info[10])
			io_cmd := io_info[len(io_info)-1]
			io_pid, _ := strconv.Atoi(io_info[7])

			r.ProcessInfo = append(r.ProcessInfo, ProcessInformation{
				WriteKbSec: io_wr,
				ReadKbSec:  io_rd,
				IoCmd:      io_cmd,
				IoPid:      io_pid,
			})
			// fmt.Println(time, io_rd, io_wr, io_cmd)
		}
	}
}
