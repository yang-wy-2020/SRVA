package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	_cpu "github.com/yang-wy-2020/SRVA/system/cpu"
	_disk "github.com/yang-wy-2020/SRVA/system/disk"
	_io "github.com/yang-wy-2020/SRVA/system/io"
	_memory "github.com/yang-wy-2020/SRVA/system/memory"
	_network "github.com/yang-wy-2020/SRVA/system/network"
	_sar "github.com/yang-wy-2020/SRVA/system/sar"
	_time "github.com/yang-wy-2020/SRVA/system/time"
	_tools "github.com/yang-wy-2020/SRVA/system/tools"
)

const cfg string = "./cfg/.config.json"

var _system _tools.System

func PermissionCheck() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}
	if currentUser.Uid != "0" {
		fmt.Println("You are not running as root.")
		os.Exit(1)
	}
}

func makeCpu(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["cpu"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	cpu_r := _cpu.GetCpuInformation(
		ret_file, cfg.CpuProcess, _tools.GetFileLineCount(ret_file))
	_system = _tools.System{
		CPU: cpu_r,
	}
}

func makeSystem(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["sar"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	system_r := _sar.GetSystemInformation(
		ret_file, _tools.GetFileLineCount(ret_file))
	_system = _tools.System{
		SYSTEM: system_r,
	}
}

func makeDIsk(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["disk"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	disk_r := _disk.GetDiskInformation(
		ret_file, _tools.GetFileLineCount(ret_file))
	_system = _tools.System{DISK: disk_r}
}

func makeIO(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["io"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	io_r := _io.GetIOInformation(
		ret_file, _tools.GetFileLineCount(ret_file))
	_system = _tools.System{IO: io_r}
}

func makeMem(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["mem"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	mem_r := _memory.GetMemInformation(
		ret_file, _tools.GetFileLineCount(ret_file))
	_system = _tools.System{MEMORY: mem_r}
}

func makeNetwork(cfg _tools.Data) {
	ret_file := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["sar"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	network_r := _network.GetNetworkInformation(
		ret_file, _tools.GetFileLineCount(ret_file), cfg.NetworkCard)
	_system = _tools.System{
		NETWORK: network_r,
	}
}

func makeTimeCheckMonitor(cfg _tools.Data) {
	time_sync := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["time_sync"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	time_monitor := _tools.FiltereServiceInformationCollect(
		cfg.MonitorService["time_monitor"], cfg.StartTime, cfg.EndTime, cfg.SavePath)
	time_r := _time.GetTimeInformation(time_sync, time_monitor,
		_tools.GetFileLineCount(time_sync), _tools.GetFileLineCount(time_monitor))
	_system = _tools.System{
		TIME: time_r,
	}
}

func main() {

	if len(os.Args) > 1 {
		var args []string
		for _, argVar := range os.Args {
			if argVar == "edit" {
				_tools.EditConfig(cfg)
			}
			if argVar == "all" {
				args = []string{"cpu", "io", "system", "disk", "network", "memory", "time"}
			}
		}
		if len(args) == 0 {
			args = os.Args[1:]
		}
		load_cfg := _tools.ReadConfig(cfg)
		log.Printf("select time:\n  %s -> %s\n", load_cfg.StartTime, load_cfg.EndTime)
		for _, arg := range args {
			switch arg {
			case "cpu":
				makeCpu(load_cfg)
			case "io":
				makeIO(load_cfg)
			case "system":
				makeSystem(load_cfg)
			case "disk":
				makeDIsk(load_cfg)
			case "memory":
				makeMem(load_cfg)
			case "network":
				makeNetwork(load_cfg)
			case "time":
				makeTimeCheckMonitor(load_cfg)
			default:
				_tools.Usage()
				os.Exit(1)
			}
			_system.NOTE = arg
			_tools.CreateLineChart(_system, load_cfg)
		}
	} else {
		_tools.Usage()
		os.Exit(1)
	}
}
