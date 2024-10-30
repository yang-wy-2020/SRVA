package tools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/yang-wy-2020/SRVA/system/cpu"
	"github.com/yang-wy-2020/SRVA/system/disk"
	_io "github.com/yang-wy-2020/SRVA/system/io"
	"github.com/yang-wy-2020/SRVA/system/memory"
	"github.com/yang-wy-2020/SRVA/system/network"
	"github.com/yang-wy-2020/SRVA/system/sar"
	"github.com/yang-wy-2020/SRVA/system/time"
)

func ReadConfig(_config string) Data {
	var SelectTime Data
	content, err := ioutil.ReadFile(_config)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	} else {
		err = json.Unmarshal(content, &SelectTime)
		if err != nil {
			log.Fatal("check: ", _config, "Error during Unmarshal(): ", err)
		}
	}
	return SelectTime
}

func FileIsEmpty(filepath string) bool {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		return false
	}
	return len(data) == 0
}

func CheckPath(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// 创建目录，权限设置为0755
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("创建目录失败:", err)
			return false
		}
		fmt.Println("目录已创建:", dir)
	}
	return true
}

func ensureNewlineAtEOF(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 移动到文件末尾
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// 读取最后一个字符
	var buffer [1]byte
	_, err = file.Read(buffer[:])
	if err != nil && err != io.EOF {
		return err
	}

	// 检查是否是换行符
	if buffer[0] != '\n' && buffer[0] != '\r' {
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}
func writeFile(str string, file string, Path string) string {
	Cmd(fmt.Sprintf("sudo chown  %d:%d %s", os.Getuid(), os.Getgid(), Path))
	f, err := os.Create(Path + file)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		if len(str) < effectiveNumber {
			return Path + file
		}
		_, _ = f.WriteString(str)
	}
	defer f.Close()
	ensureNewlineAtEOF(Path + file)
	return Path + file
}
func Cmd(s string) string {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return ""
	}
	Str := strings.TrimRight(out.String(), "\n")
	return Str
}

func FiltereServiceInformationCollect(service_name, s_time, e_time string, savePath string) string {
	cmd := fmt.Sprintf("sudo %s -u %s --since '%s' --until '%s'",
		journal, service_name, s_time, e_time)
	ret := Cmd(cmd)
	return writeFile(ret, service_name, savePath)
}

func EditConfig(config string) {
	Cmd(fmt.Sprintf("sudo chown  %d:%d %s", os.Getuid(), os.Getgid(), config))
	cmd := exec.Command("vim", config)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal("err start Vim:", err)
		return
	}
	editedText, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal("err read file:", err)
		return
	}
	fmt.Println("edit after:")
	fmt.Println(string(editedText))
	var SelectTime Data
	err = json.Unmarshal(editedText, &SelectTime)
	if err != nil {
		fmt.Println("\x1b[31mcheck: ", config, "Error during Unmarshal(): ", err, "\x1b[0m")
	}
	os.Exit(-1)
}

func NotContainsString(slice []opts.LineData, str string) bool {
	for _, item := range slice {
		if item.Value.(string) == str {
			return false
		}
	}
	return true
}

func GetIoValue(i int, str string, num int, data _io.IoInformation) bool {
	var j int
	if i+num >= len(data.Time) {
		for j = i; j < len(data.Time); j++ {
			if data.ProcessInfo[j].IoCmd == str && data.Time[j] == data.Time[i] {
				return true
			}
		}
		return false
	}
	for j = i; j < i+num; j++ {
		if data.ProcessInfo[j].IoCmd == str && data.Time[j] == data.Time[i] {
			return true
		}
	}
	return false
}

func GetMemValue(i int, str string, num int, data memory.MemoryInformation) bool {
	var j int
	if i+num >= len(data.Time) {
		for j = i; j < len(data.Time); j++ {
			if data.ProcessInfo[j].MemoryCmd == str && data.Time[j] == data.Time[i] {
				return true
			}
		}
		return false
	}
	for j = i; j < i+num; j++ {
		if data.ProcessInfo[j].MemoryCmd == str && data.Time[j] == data.Time[i] {
			return true
		}
	}
	return false
}

func GetFileLineCount(filepath string) int {
	if FileIsEmpty(filepath) {
		log.Println("file is empty: ", filepath)
		os.Exit(1)
	}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
		return -1
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineCount := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break // End of file or error
		}
		lineCount++
	}
	return lineCount
}

func PrefixCreateChart(_charts *charts.Line, description string) *charts.Line {
	_charts.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  chartsWidth,            // 设置图表宽度
			Height: chartsHeight,           // 设置图表高度
			Theme:  types.ThemeInfographic, // 确保 ThemeInfographic 是有效的主题
		}),
		charts.WithTitleOpts(opts.Title{
			Title: description,
			// Subtitle: "",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			// 启用数据窗口组件，设置x轴可以缩放
			XAxisIndex: []int{0},
		}),
	)
	return _charts
}

func PostfixCreateChart(f *os.File, charts ...*charts.Line) {
	for _, chart := range charts {
		err := chart.Render(f)
		if err != nil {
			panic(err)
		}
	}
}

func CreateLineChart(sys System, cfg Data) {
	if !CheckPath(cfg.OutputPath) {
		fmt.Println(cfg.OutputPath + "is not exist!!")
	}
	f, err := os.Create(cfg.OutputPath + fmt.Sprintf("%s.html", sys.NOTE))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	switch sys.NOTE {
	case "cpu":
		LoadCharts := charts.NewLine()
		LoadCharts = PrefixCreateChart(LoadCharts, "LoadAvg 1m, 5m, 15m")
		one, five, fifteem := cpuLoadAvgGenerateLineItems(sys.CPU)
		LoadCharts.SetXAxis(sys.CPU.Time).
			AddSeries("LoadAvg 1m", one).AddSeries("LoadAvg 5m", five).AddSeries("LoadAvg 15m", fifteem)

		RateAndCoreCharts := charts.NewLine()
		RateAndCoreCharts = PrefixCreateChart(RateAndCoreCharts, "Cpu Rate And Core")
		RateAndCoreCharts.SetXAxis(sys.CPU.Time).
			AddSeries("Cpu rate", cpuGenerateLineItems(sys.CPU, "rate", 0)).
			AddSeries("Cpu core", cpuGenerateLineItems(sys.CPU, "core", 0))

		ProcessCharts := charts.NewLine()
		ProcessCharts = PrefixCreateChart(ProcessCharts, "Cpu Process Info")
		for i := 0; i < len(sys.CPU.ProcessInfo); i++ {
			ProcessCharts.SetXAxis(sys.CPU.Time).
				AddSeries(sys.CPU.ProcessInfo[i].Name+"-CpuUse", cpuGenerateLineItems(sys.CPU, "CpuUse", i)).
				AddSeries(sys.CPU.ProcessInfo[i].Name+"-SystemUse", cpuGenerateLineItems(sys.CPU, "SystemUse", i)).
				AddSeries(sys.CPU.ProcessInfo[i].Name+"-UserUse", cpuGenerateLineItems(sys.CPU, "UserUse", i))
		}
		PostfixCreateChart(f, LoadCharts, RateAndCoreCharts, ProcessCharts)
	case "io":
		{
			IoWrCharts := charts.NewLine()
			IoWrCharts = PrefixCreateChart(IoWrCharts, "Disk Write Io Rate")

			IoRdCharts := charts.NewLine()
			IoRdCharts = PrefixCreateChart(IoRdCharts, "Disk Read Io Rate")
			if len(sys.IO.ProcessInfo) == 0 {
				fmt.Println("NO IO DATA")
				os.Exit(0)
			}
			cmd_items := ioGenerateLineItems(sys.IO, "cmd", "0", 0)
			time_items := ioGenerateLineItems(sys.IO, "time", "0", 0)
			len_cmd := len(cmd_items)
			for i := 0; i < len_cmd; i++ {
				io_cmd := cmd_items[i].Value.(string)
				io_items := ioGenerateLineItems(sys.IO, "wr", io_cmd, len_cmd)
				IoWrCharts.SetXAxis(time_items).
					AddSeries(io_cmd+"-wr", io_items)

				io_items = ioGenerateLineItems(sys.IO, "rd", io_cmd, len_cmd)
				IoRdCharts.SetXAxis(time_items).
					AddSeries(io_cmd+"-rd", io_items)
			}
			PostfixCreateChart(f, IoWrCharts, IoRdCharts)
		}
	case "system":
		{
			SystemTotalCharts := charts.NewLine()
			SystemTotalCharts = PrefixCreateChart(SystemTotalCharts, "Cpu Total Information")
			SystemTotalCharts.SetXAxis(sys.SYSTEM.Time).
				AddSeries("User", systemGenerateLineItems(sys.SYSTEM, "user")).
				AddSeries("System", systemGenerateLineItems(sys.SYSTEM, "system")).
				AddSeries("Idle", systemGenerateLineItems(sys.SYSTEM, "idle"))
			PostfixCreateChart(f, SystemTotalCharts)

		}
	case "network":
		{
			NetwrodCardCharts := charts.NewLine()
			NetwrodCardCharts = PrefixCreateChart(NetwrodCardCharts, "NetWork Card Rx/Tx View")

			for i := 0; i < len(sys.NETWORK.NetCard); i++ {
				if len(sys.NETWORK.NetCard[i].NetworkTotal) == 0 {
					continue
				}
				rxp, txp, rx, tx := networkGenerateLineItems(sys.NETWORK, i)
				NetwrodCardCharts.SetXAxis(sys.NETWORK.NetCard[i].Time).
					AddSeries(sys.NETWORK.NetCard[i].Name+"-Rxp", rxp).
					AddSeries(sys.NETWORK.NetCard[i].Name+"-Txp", txp).
					AddSeries(sys.NETWORK.NetCard[i].Name+"-Rx", rx).
					AddSeries(sys.NETWORK.NetCard[i].Name+"-Tx", tx)
			}
			PostfixCreateChart(f, NetwrodCardCharts)
		}
	case "disk":
		{
			NvmeCharts := charts.NewLine()
			NvmeCharts = PrefixCreateChart(NvmeCharts, "Nvme Used View")
			time_items, used_items, free_items, disk_device := diskGenerateLineItems(sys.DISK, "nvme0n1p1")
			NvmeCharts.SetXAxis(time_items).
				AddSeries(disk_device+"-Used", used_items).
				AddSeries(disk_device+"-Free", free_items)
			PostfixCreateChart(f, NvmeCharts)

			time_items, used_items, free_items, disk_device = diskGenerateLineItems(sys.DISK, "root")
			RootCharts := charts.NewLine()
			RootCharts = PrefixCreateChart(RootCharts, "SysDisk Used View")
			RootCharts.SetXAxis(time_items).
				AddSeries(disk_device+"-Used", used_items).
				AddSeries(disk_device+"-Free", free_items)
			PostfixCreateChart(f, RootCharts)
		}
	case "memory":
		{
			MemPerCharts := charts.NewLine()
			MemPerCharts = PrefixCreateChart(MemPerCharts, "mem percent")

			MemVszCharts := charts.NewLine()
			MemVszCharts = PrefixCreateChart(MemVszCharts, "mem VirtualMemorySize")

			MemRssCharts := charts.NewLine()
			MemRssCharts = PrefixCreateChart(MemRssCharts, "mem ResidentSetSize")

			cmd_items := memGenerateLineItems(sys.MEMORY, "cmd", "0", 0)
			time_items := memGenerateLineItems(sys.MEMORY, "time", "0", 0)
			len_cmd := len(cmd_items)
			for i := 0; i < len_cmd; i++ {
				mem_cmd := cmd_items[i].Value.(string)
				mem_items := memGenerateLineItems(sys.MEMORY, "percent", mem_cmd, len_cmd)
				MemPerCharts.SetXAxis(time_items).
					AddSeries(mem_cmd+"-percent", mem_items)

				mem_items = memGenerateLineItems(sys.MEMORY, "vsz", mem_cmd, len_cmd)
				MemVszCharts.SetXAxis(time_items).
					AddSeries(mem_cmd+"-vsz", mem_items)

				mem_items = memGenerateLineItems(sys.MEMORY, "rss", mem_cmd, len_cmd)
				MemRssCharts.SetXAxis(time_items).
					AddSeries(mem_cmd+"-vsz", mem_items)
			}
			PostfixCreateChart(f, MemPerCharts, MemVszCharts, MemRssCharts)
		}
	case "time":
		TimeSyncCharts := charts.NewLine()
		TimeMonitorCharts := charts.NewLine()
		TimeSyncCharts = PrefixCreateChart(TimeSyncCharts, "Time sync status")
		TimeMonitorCharts = PrefixCreateChart(TimeMonitorCharts, "Time sync offset status")
		eth1, eth2, lock_status := timeGenerateLineItems(sys.TIME)
		for i := 0; i < len(sys.TIME.CardSyncState); i++ {
			TimeSyncCharts.SetXAxis(sys.TIME.Time).
				AddSeries("eth1_status", eth1).
				AddSeries("eth2_status", eth2).
				AddSeries("lock status", lock_status)
		}
		t1, eth1p2m, eth1s2p, eth1s2m, eth1_c := timeMonitorGenerateLineItems(sys.TIME, "eth1")
		_, eth2p2m2, eth2s2p, eth2s2m, eth2_c := timeMonitorGenerateLineItems(sys.TIME, "eth2")
		for e := 0; e < len(sys.TIME.Eth1TimeOffset); e++ {
			TimeMonitorCharts.SetXAxis(t1).
				AddSeries(eth1_c+"-PhcToMasterOffset", eth1p2m).
				AddSeries(eth1_c+"-SystemToPhcOffset", eth1s2p).
				AddSeries(eth1_c+"-SystemToMasterOffset", eth1s2m).
				AddSeries(eth2_c+"-PhcToMasterOffset", eth2p2m2).
				AddSeries(eth2_c+"-SystemToPhcOffset", eth2s2p).
				AddSeries(eth2_c+"-SystemToMasterOffset", eth2s2m)
		}
		PostfixCreateChart(f, TimeSyncCharts, TimeMonitorCharts)
	}
}

func timeGenerateLineItems(data time.TimeCheckInformation) ([]opts.LineData, []opts.LineData, []opts.LineData) {
	// sync
	eth1_sync_items := make([]opts.LineData, 0)
	eth2_sync_items := make([]opts.LineData, 0)
	lock_sync_items := make([]opts.LineData, 0)
	for i := 0; i < len(data.CardSyncState); i++ {
		eth1_sync_items = append(eth1_sync_items, opts.LineData{
			Value: data.CardSyncState[i].Eth1SyncState})
		eth2_sync_items = append(eth2_sync_items, opts.LineData{
			Value: data.CardSyncState[i].Eth2SyncState})
		lock_sync_items = append(lock_sync_items, opts.LineData{
			Value: data.CardSyncState[i].LockState})
	}
	return eth1_sync_items, eth2_sync_items, lock_sync_items
}
func timeMonitorGenerateLineItems(data time.TimeCheckInformation, eth string) (
	[]string, []opts.LineData, []opts.LineData, []opts.LineData, string) {
	p2m_items := make([]opts.LineData, 0)
	s2p_sync_items := make([]opts.LineData, 0)
	s2m_sync_items := make([]opts.LineData, 0)
	t := make([]string, 0)
	if eth == "eth1" {
		for i := 0; i < len(data.Eth1TimeOffset); i++ {
			t = append(t, data.Eth1TimeOffset[i].Time)
			p2m_items = append(p2m_items, opts.LineData{
				Value: data.Eth1TimeOffset[i].PhcToMasterOffset,
			})
			s2p_sync_items = append(s2p_sync_items, opts.LineData{
				Value: data.Eth1TimeOffset[i].SystemToPhcOffset,
			})
			s2m_sync_items = append(s2m_sync_items, opts.LineData{
				Value: data.Eth1TimeOffset[i].SystemToMasterOffset,
			})
		}
		return t, p2m_items, s2p_sync_items, s2m_sync_items, eth
	} else {
		for i := 0; i < len(data.Eth2TimeOffset); i++ {
			t = append(t, data.Eth2TimeOffset[i].Time)
			p2m_items = append(p2m_items, opts.LineData{
				Value: data.Eth2TimeOffset[i].PhcToMasterOffset,
			})
			s2p_sync_items = append(s2p_sync_items, opts.LineData{
				Value: data.Eth2TimeOffset[i].SystemToPhcOffset,
			})
			s2m_sync_items = append(s2m_sync_items, opts.LineData{
				Value: data.Eth2TimeOffset[i].SystemToMasterOffset,
			})
		}
		return t, p2m_items, s2p_sync_items, s2m_sync_items, eth
	}
}

func systemGenerateLineItems(data sar.SystemInformation, note string) []opts.LineData {
	var value interface{}
	items := make([]opts.LineData, 0)

	for i := 0; i < len(data.CpuTotal); i++ {
		if note == "user" {
			value = data.CpuTotal[i].User
		}
		if note == "system" {
			value = data.CpuTotal[i].System
		}
		if note == "idle" {
			value = data.CpuTotal[i].Idle
		}
		items = append(items, opts.LineData{Value: value})
	}

	return items
}

func networkGenerateLineItems(data network.NetworkInformation, index int) (
	[]opts.LineData, []opts.LineData, []opts.LineData, []opts.LineData) {
	rxp_items := make([]opts.LineData, 0)
	txp_items := make([]opts.LineData, 0)
	rx_items := make([]opts.LineData, 0)
	tx_items := make([]opts.LineData, 0)
	for _, Ncard := range data.NetCard[index].NetworkTotal {
		rxp_items = append(rxp_items, opts.LineData{Value: Ncard.Rxpck})
		txp_items = append(txp_items, opts.LineData{Value: Ncard.Txpck})
		rx_items = append(rx_items, opts.LineData{Value: Ncard.Rxkb})
		tx_items = append(tx_items, opts.LineData{Value: Ncard.Txkb})
	}
	return rxp_items, txp_items, rx_items, tx_items
}

func cpuGenerateLineItems(data cpu.CpuInformation, note string, index int) []opts.LineData {
	if note == "CpuUse" {
		cpu_items := make([]opts.LineData, 0)
		for _, use := range data.ProcessInfo[index].CpuUse {
			cpu_items = append(cpu_items, opts.LineData{Value: use})
		}
		return cpu_items
	}
	if note == "SystemUse" {
		system_items := make([]opts.LineData, 0)
		for _, use := range data.ProcessInfo[index].SystemUse {
			system_items = append(system_items, opts.LineData{Value: use})
		}
		return system_items
	}
	if note == "UserUse" {
		user_items := make([]opts.LineData, 0)
		for _, use := range data.ProcessInfo[index].UserUse {
			user_items = append(user_items, opts.LineData{Value: use})
		}
		return user_items
	}
	if note == "rate" {
		rate_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.Rate); i++ {
			rate_items = append(rate_items, opts.LineData{Value: data.Rate[i]})
		}
		return rate_items
	}
	if note == "core" {
		core_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.CpuCore); i++ {
			core_items = append(core_items, opts.LineData{Value: data.CpuCore[i]})
		}
		return core_items
	}
	return nil
}

func cpuLoadAvgGenerateLineItems(data cpu.CpuInformation) (
	[]opts.LineData, []opts.LineData, []opts.LineData) {

	One_items := make([]opts.LineData, 0)
	Five_items := make([]opts.LineData, 0)
	Fifteen_items := make([]opts.LineData, 0)
	for i := 0; i < len(data.LoadAvg); i++ {
		One_items = append(One_items, opts.LineData{Value: data.LoadAvg[i].AvgOne})
		Five_items = append(Five_items, opts.LineData{Value: data.LoadAvg[i].AvgFive})
		Fifteen_items = append(Fifteen_items, opts.LineData{Value: data.LoadAvg[i].AvgFifteen})
	}
	return One_items, Five_items, Fifteen_items
}

func diskGenerateLineItems(data disk.DiskInformation, note1 string) ([]opts.LineData, []opts.LineData, []opts.LineData, string) {
	var disk_device string
	if note1 == "nvme0n1p1" {
		time_items := make([]opts.LineData, 0)
		used_items := make([]opts.LineData, 0)
		free_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.DiskInfo); i++ {
			if data.DiskInfo[i].DiskDevice == "/dev/nvme0n1p1" {
				if NotContainsString(time_items, data.Time[i]) {
					time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				}
				used_items = append(used_items, opts.LineData{Value: data.DiskInfo[i].DiskUsed})
				free_items = append(free_items, opts.LineData{Value: data.DiskInfo[i].DiskFree})
				disk_device = data.DiskInfo[i].DiskDevice
			}
		}
		return time_items, used_items, free_items, disk_device
	}
	if note1 == "root" {
		time_items := make([]opts.LineData, 0)
		used_items := make([]opts.LineData, 0)
		free_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.DiskInfo); i++ {
			if data.DiskInfo[i].DiskDevice != "/dev/nvme0n1p1" {
				if NotContainsString(time_items, data.Time[i]) {
					time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				}
				used_items = append(used_items, opts.LineData{Value: data.DiskInfo[i].DiskUsed})
				free_items = append(free_items, opts.LineData{Value: data.DiskInfo[i].DiskFree})
				disk_device = data.DiskInfo[i].DiskDevice
			}
		}
		return time_items, used_items, free_items, disk_device
	}
	items := make([]opts.LineData, 0)
	return items, items, items, disk_device
}

func ioGenerateLineItems(data _io.IoInformation, note string, str string, num int) []opts.LineData {
	items := make([]opts.LineData, 0)
	if note == "cmd" {
		items = append(items, opts.LineData{Value: data.ProcessInfo[0].IoCmd})
		for i := 0; i < len(data.ProcessInfo); i++ {
			str := data.ProcessInfo[i].IoCmd
			if NotContainsString(items, str) {
				items = append(items, opts.LineData{Value: str})
			}
		}
	}
	if note == "time" {
		items = append(items, opts.LineData{Value: data.Time[0]})
		for i := 0; i < len(data.Time); i++ {
			str := data.Time[i]
			if NotContainsString(items, str) {
				items = append(items, opts.LineData{Value: str})
			}
		}
	}
	i := 0
	time_items := make([]opts.LineData, 0)
	if note == "wr" {
		for {
			if i < len(data.ProcessInfo) {
				if data.ProcessInfo[i].IoCmd == str {
					items = append(items, opts.LineData{Value: data.ProcessInfo[i].WriteKbSec})
					time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				} else {
					get_value := GetIoValue(i, str, num, data)
					if !get_value && NotContainsString(time_items, data.Time[i]) {
						items = append(items, opts.LineData{Value: 0})
						time_items = append(time_items, opts.LineData{Value: data.Time[i]})
					}
				}
				i++
			} else {
				break
			}
		}
	}
	if note == "rd" {
		for {
			if i < len(data.ProcessInfo) {
				if data.ProcessInfo[i].IoCmd == str {
					items = append(items, opts.LineData{Value: data.ProcessInfo[i].ReadKbSec})
					time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				} else {
					get_value := GetIoValue(i, str, num, data)
					if !get_value && NotContainsString(time_items, data.Time[i]) {
						items = append(items, opts.LineData{Value: 0})
						time_items = append(time_items, opts.LineData{Value: data.Time[i]})
					}
				}
				i++
			} else {
				break
			}
		}
	}
	return items
}

func memGenerateLineItems(data memory.MemoryInformation, note string, str string, num int) []opts.LineData {
	items := make([]opts.LineData, 0)
	if note == "cmd" {
		items = append(items, opts.LineData{Value: data.ProcessInfo[0].MemoryCmd})
		for i := 0; i < len(data.ProcessInfo); i++ {
			str := data.ProcessInfo[i].MemoryCmd
			if NotContainsString(items, str) {
				items = append(items, opts.LineData{Value: str})
			}
		}
	}
	if note == "time" {
		items = append(items, opts.LineData{Value: data.Time[0]})
		for i := 0; i < len(data.Time); i++ {
			str := data.Time[i]
			if NotContainsString(items, str) {
				items = append(items, opts.LineData{Value: str})
			}
		}
	}
	i := 0
	time_items := make([]opts.LineData, 0)
	if note == "percent" {
		for {
			if i < len(data.ProcessInfo) {
				if data.ProcessInfo[i].MemoryCmd == str {
					items = append(items, opts.LineData{Value: data.ProcessInfo[i].MemoryPercentage})
				} else {
					get_value := GetMemValue(i, str, num, data)
					if !get_value && NotContainsString(time_items, data.Time[i]) {
						items = append(items, opts.LineData{Value: 0})
					}
				}
				time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				i++
			} else {
				break
			}
		}
	}
	if note == "vsz" {
		for {
			if i < len(data.ProcessInfo) {
				if data.ProcessInfo[i].MemoryCmd == str {
					items = append(items, opts.LineData{Value: data.ProcessInfo[i].VirtualMemorySize})
				} else {
					get_value := GetMemValue(i, str, num, data)
					if !get_value && NotContainsString(time_items, data.Time[i]) {
						items = append(items, opts.LineData{Value: 0})
					}
				}
				time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				i++
			} else {
				break
			}
		}
	}
	if note == "rss" {
		for {
			if i < len(data.ProcessInfo) {
				if data.ProcessInfo[i].MemoryCmd == str {
					items = append(items, opts.LineData{Value: data.ProcessInfo[i].ResidentSetSize})
				} else {
					get_value := GetMemValue(i, str, num, data)
					if !get_value && NotContainsString(time_items, data.Time[i]) {
						items = append(items, opts.LineData{Value: 0})
					}
				}
				time_items = append(time_items, opts.LineData{Value: data.Time[i]})
				i++
			} else {
				break
			}
		}
	}
	return items
}
