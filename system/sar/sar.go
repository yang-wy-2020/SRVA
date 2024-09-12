package sar

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func stringToFloat64(str string) float64 {
	b2, _ := strconv.ParseFloat(str, 32)
	return b2
}

func checkNetworkInterface(name string) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return false
	}
	for _, intf := range interfaces {
		if intf.Name == name {
			return true
		}
	}
	return false
}

func GetSystemInformation(system_file string, fileLine int, NetCard []string) SystemInformation {
	var s SystemInformation
	currentUser, _ := os.Hostname()
	var values, netlist []string

	file_buf, err := os.Open(system_file)
	if err != nil {
		fmt.Println("open file failed:", err)
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)

	for _, intf := range NetCard {
		if checkNetworkInterface(intf) {
			netlist = append(netlist, intf)
		}
	}
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
		for index, intf := range netlist {
			if strings.Contains(line, intf) {
				pattern := fmt.Sprintf(`%s \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s`, intf)
				re := regexp.MustCompile(pattern)
				match := re.FindStringSubmatch(line)
				if len(match) > 1 {
					values = match[1:]
				}
				s.NetworkTotal = append(s.NetworkTotal, NetworkCard{
					Name: intf,
				})
				s.NetworkTotal[index].Rxpck = append(s.NetworkTotal[index].Rxpck, stringToFloat64(values[0]))
				s.NetworkTotal[index].Txpck = append(s.NetworkTotal[index].Txpck, stringToFloat64(values[1]))
				s.NetworkTotal[index].Rxkb = append(s.NetworkTotal[index].Rxkb, stringToFloat64(values[2]))
				s.NetworkTotal[index].Txkb = append(s.NetworkTotal[index].Txkb, stringToFloat64(values[3]))
			}
		}
	}
	return s
}
