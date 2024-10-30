package network

import (
	"bufio"
	"fmt"
	"io"
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
func GetNetworkInformation(system_file string, fileLine int, NetCard []string) NetworkInformation {
	var s NetworkInformation
	var values []string
	currentUser, err := os.Hostname()
	if err != nil {
		return s
	}

	file_buf, err := os.Open(system_file)
	if err != nil {
		return s
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)

	// 过滤有效的网络接口
	var netlist []string
	for _, intf := range NetCard {
		if checkNetworkInterface(intf) {
			netlist = append(netlist, intf)
		}
	}

	// 初始化网络接口切片
	s.NetCard = make([]NetWorkCard, len(netlist))
	for index, intf := range netlist {
		s.NetCard[index].Name = intf
		s.NetCard[index].Time = make([]string, 0)
		s.NetCard[index].NetworkTotal = make([]NetWorkIndicators, 0)
	}
	for i := 0; i < fileLine; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return s
			}
			break
		}
		_time := strings.Split(line, currentUser)[0]
		for index, intf := range netlist {
			if strings.Contains(line, intf) {
				// pattern := fmt.Sprintf(`%s \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s+ (\d+\.\d+) \s`, intf)
				pattern := fmt.Sprintf(`.* %s .*`, intf)
				re := regexp.MustCompile(pattern)
				match := re.FindStringSubmatch(line)
				if len(match) > 0 {
					net_info := strings.Fields(match[0])
					values = net_info[8:] // 确保 values 切片包含足够的元素
					netInfo := NetWorkIndicators{
						Rxpck: stringToFloat64(values[0]),
						Txpck: stringToFloat64(values[1]),
						Rxkb:  stringToFloat64(values[2]),
						Txkb:  stringToFloat64(values[3]),
					}
					s.NetCard[index].Time = append(s.NetCard[index].Time, _time)
					s.NetCard[index].NetworkTotal = append(s.NetCard[index].NetworkTotal, netInfo)
				}
			}
		}
	}
	return s
}
