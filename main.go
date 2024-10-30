package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	_cpu "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/cpu"
	_disk "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/disk"
	_io "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/io"
	_memory "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/memory"
	_network "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/network"
	_sar "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/sar"
	_tools "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/tools"
)

const cfg string = "/opt/qomolo/utils/qomolo-system-analysis/cfg/config.json"

// const cfg string = "./cfg/.config.json"

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

func main() {

	if len(os.Args) > 1 {
		var args []string
		for _, argVar := range os.Args {
			if argVar == "edit" {
				_tools.EditConfig(cfg)
			}
			if argVar == "all" {
				args = []string{"cpu", "io", "system", "disk", "network", "memory"}
			}
		}
		if len(args) == 0 {
			args = os.Args[1:]
		}
		load_cfg := _tools.ReadConfig(cfg)
		fmt.Println(load_cfg.OutputPath)
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
			default:
				_tools.Usage()
				os.Exit(1)
			}
			log.Printf("select time:\n  %s -> %s\n", load_cfg.StartTime, load_cfg.EndTime)
			_system.NOTE = arg
			_tools.CreateLineChart(_system, load_cfg)
		}
	} else {
		_tools.Usage()
		os.Exit(1)
	}
}
