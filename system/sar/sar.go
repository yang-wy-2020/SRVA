package sar

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func stringToFloat64(str string) float64 {
	b2, _ := strconv.ParseFloat(str, 32)
	return b2
}

func GetSystemInformation(system_file string, fileLine int) SystemInformation {
	var s SystemInformation
	currentUser, _ := os.Hostname()
	var values []string

	file_buf, err := os.Open(system_file)
	if err != nil {
		fmt.Println("open file failed:", err)
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)

	for i := 0; i < fileLine; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, Cpu_flag) {
			re := regexp.MustCompile(`all\s+(\d+\.\d+)\s+(\d+\.\d+)\s+(\d+\.\d+)\s+(\d+\.\d+)\s+(\d+\.\d+)\s+(\d+\.\d+)`)
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				values = match[1:]
			}
			//  get time
			time := strings.Split(line, currentUser)[0]
			s.Time = append(s.Time, time)
			// cpu total
			s.CpuTotal = append(s.CpuTotal, SystemCpu{
				User:   stringToFloat64(values[0]),
				System: stringToFloat64(values[2]),
				IoWait: stringToFloat64(values[3]),
				Idle:   stringToFloat64(values[5]),
			})
		}
	}
	return s
}
