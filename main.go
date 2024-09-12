package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	_cpu "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/cpu"
	_sar "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/sar"
	_tools "gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/tools"
)

// const cfg string = "/opt/qomolo/utils/qomolo-system-analysis/cfg/config.json"
const cfg string = "./cfg/.config.json"

func main() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	if currentUser.Uid != "0" {
		fmt.Println("You are not running as root.")
		// os.Exit(1)
	}

	load_cfg := _tools.ReadConfig(cfg)
	var _system _tools.System
	var flag string = "defalut"
	if len(os.Args) > 1 {
		flag = os.Args[1]
	}

	switch flag {
	case "defalut":
		fmt.Println(`args1:
    cpu
        - cpu view
    io 
        - io view
    system
        - cpu total view 
    network
        - network card view
    edit
        - edit config
e.g.
sudo qomolo-sys-analysis $(args1)
		`)
		os.Exit(1)
	case "edit":
		_tools.EditConfig(cfg)
	case "cpu": // get cpu info
		ret_file := _tools.FiltereServiceInformationCollect(
			load_cfg.MonitorService["cpu"], load_cfg.StartTime, load_cfg.EndTime)
		cpu_r := _cpu.GetCpuInformation(
			ret_file, load_cfg.ModelsList, _tools.GetFileLineCount(ret_file))
		_system = _tools.System{
			CPU: cpu_r,
		}
	case "io": // get io info
		_system = _tools.System{}
	case "system":
		ret_file := _tools.FiltereServiceInformationCollect(
			load_cfg.MonitorService["sar"], load_cfg.StartTime, load_cfg.EndTime)
		system_r := _sar.GetSystemInformation(
			ret_file, _tools.GetFileLineCount(ret_file), load_cfg.NetworkCard)
		_system = _tools.System{
			SYSTEM: system_r,
		}
	}
	log.Printf("select time:\n  %s -> %s\n", load_cfg.StartTime, load_cfg.EndTime)
	_system.NOTE = os.Args[1]
	_tools.CreateLineChart(_system)
}
