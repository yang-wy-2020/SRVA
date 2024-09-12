package cpu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetCpuInformation(file string, models []string, counts int) CpuInformation {
	var r CpuInformation
	getCpuRateAndRateAndCpuCore(&r, file, counts)
	getModelsCpuInfomation(&r, models)
	return r
}

func stringToFloat64(str string) float64 {
	b2, _ := strconv.ParseFloat(str, 32)
	return b2
}

func StringToUint8(str string) uint8 {
	// 将字符串转换为整数
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	// 检查数值是否在 uint8 的范围内
	if i < 0 || i > 255 {
		return 0
		// 转换为 uint8 并返回
	}
	return uint8(i)

}

func getModelsCpuInfomation(r *CpuInformation, models []string) {

}

func getCpuRateAndRateAndCpuCore(r *CpuInformation, file string, counts int) {
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
	rate := regexp.MustCompile(`Freq Avg ([\d\.]+) GHz`)
	loadavg := regexp.MustCompile(`Load Avg ([0-9]+\.[0-9]+) ([0-9]+\.[0-9]+) ([0-9]+\.[0-9]+)`)
	core := regexp.MustCompile(`(\d+) CPUs`)

	for i := 0; i < counts; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, LoadAvg_g) {
			// add data loadAvgs
			loadavg_r := loadavg.FindStringSubmatch(line)
			if len(loadavg_r) < 4 {
				log.Fatal("No matches found")
			}
			r.LoadAvg = append(r.LoadAvg, LoadAvgInformation{
				AvgOne:     stringToFloat64(loadavg_r[1]),
				AvgFive:    stringToFloat64(loadavg_r[2]),
				AvgFifteen: stringToFloat64(loadavg_r[3])})

			// add data time
			time := strings.Split(line, currentUser)[0]
			r.Time = append(r.Time, time)

			// add data cpu rate
			rate_r := rate.FindStringSubmatch(line)
			r.Rate = append(r.Rate, stringToFloat64(rate_r[1]))

			// add data cpu core
			core_r := core.FindStringSubmatch(line)[0]
			core_real := strings.Split(core_r, " ")[0]
			r.CpuCore = append(r.CpuCore, StringToUint8(core_real))
		}
	}
}
