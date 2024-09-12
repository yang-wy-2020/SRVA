package tools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/cpu"
	"gitlab.qomolo.com/xiangyang.chen/qomolo-system-analysis/system/sar"
)

func ReadConfig(_config string) Data {
	var SelectTime Data
	content, err := ioutil.ReadFile(_config)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	} else {
		err = json.Unmarshal(content, &SelectTime)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}
	}
	return SelectTime
}
func writeFile(str string, file string) string {
	f, err := os.Create(output + file)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, _ = f.WriteString(str)
	}
	defer f.Close()
	return output + file
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

func FiltereServiceInformationCollect(service_name, s_time, e_time string) string {
	cmd := fmt.Sprintf("sudo %s -u %s --since '%s' --until '%s'",
		journal, service_name, s_time, e_time)
	ret := Cmd(cmd)
	return writeFile(ret, service_name)
}

func EditConfig(config string) {
	Cmd(fmt.Sprintf("sudo chown  %d:%d %s", os.Getuid(), os.Getgid(), config))
	cmd := exec.Command("vim", config)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("err start Vim:", err)
		return
	}
	editedText, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println("err read file:", err)
		return
	}
	fmt.Println("edit after:")
	fmt.Println(string(editedText))
	os.Exit(-1)
}

func GetFileLineCount(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
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

func CreateLineChart(sys System) {

	f, err := os.Create(fmt.Sprintf("%s.html", sys.NOTE))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	switch sys.NOTE {
	case "cpu":
		LoadCharts := charts.NewLine()
		LoadCharts.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				Width:  chartsWidth,            // 设置图表宽度
				Height: chartsHeight,           // 设置图表高度
				Theme:  types.ThemeInfographic, // 确保 ThemeInfographic 是有效的主题
			}),
			charts.WithTitleOpts(opts.Title{
				Title: "LoadAvg 1m, 5m, 15m",
				// Subtitle: "",
			}),
			charts.WithDataZoomOpts(opts.DataZoom{
				// 启用数据窗口组件，设置x轴可以缩放
				XAxisIndex: []int{0},
			}),
		)

		LoadCharts.SetXAxis(sys.CPU.Time).
			AddSeries("LoadAvg 1m", cpuGenerateLineItems(sys.CPU, "one")).
			AddSeries("LoadAvg 5m", cpuGenerateLineItems(sys.CPU, "five")).
			AddSeries("LoadAvg 15m", cpuGenerateLineItems(sys.CPU, "fifteen"))

		err = LoadCharts.Render(f)
		if err != nil {
			panic(err)
		}

		RateAndCoreCharts := charts.NewLine()
		RateAndCoreCharts.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				Width:  chartsWidth,
				Height: chartsHeight,
				Theme:  types.ThemeInfographic,
			}),
			charts.WithTitleOpts(opts.Title{
				Title: "Cpu rate and core",
			}),
			charts.WithDataZoomOpts(opts.DataZoom{
				// 启用数据窗口组件，设置x轴可以缩放
				XAxisIndex: []int{0},
			}),
		)
		RateAndCoreCharts.SetXAxis(sys.CPU.Time).
			AddSeries("Cpu rate", cpuGenerateLineItems(sys.CPU, "rate")).
			AddSeries("Cpu core", cpuGenerateLineItems(sys.CPU, "core"))

		err = RateAndCoreCharts.Render(f)
		if err != nil {
			panic(err)
		}
	case "io":
		{
			fmt.Println("io view")
		}
	case "system":
		{
			LoadCharts := charts.NewLine()
			LoadCharts.SetGlobalOptions(
				charts.WithInitializationOpts(opts.Initialization{
					Width:  chartsWidth,            // 设置图表宽度
					Height: chartsHeight,           // 设置图表高度
					Theme:  types.ThemeInfographic, // 确保 ThemeInfographic 是有效的主题
				}),
				charts.WithTitleOpts(opts.Title{
					Title: "CPU TOTAL INFORMATION",
					// Subtitle: "",
				}),
				charts.WithDataZoomOpts(opts.DataZoom{
					// 启用数据窗口组件，设置x轴可以缩放
					XAxisIndex: []int{0},
				}),
			)

			LoadCharts.SetXAxis(sys.SYSTEM.Time).
				AddSeries("User", systemGenerateLineItems(sys.SYSTEM, "user", 0)).
				AddSeries("System", systemGenerateLineItems(sys.SYSTEM, "system", 0)).
				AddSeries("Idle", systemGenerateLineItems(sys.SYSTEM, "idle", 0))

			// TODO network card view
			NetwrodCardCharts := charts.NewLine()
			NetwrodCardCharts.SetGlobalOptions(
				charts.WithInitializationOpts(opts.Initialization{
					Width:  chartsWidth,
					Height: chartsHeight,
					Theme:  types.ThemeInfographic,
				}),
				charts.WithTitleOpts(opts.Title{
					Title: "NetWork Card view",
				}),
				charts.WithDataZoomOpts(opts.DataZoom{
					// 启用数据窗口组件，设置x轴可以缩放
					XAxisIndex: []int{0},
				}),
			)
			for i := 0; i < len(sys.SYSTEM.NetworkTotal); i++ {
				NetwrodCardCharts.SetXAxis(sys.SYSTEM.Time).
					AddSeries(sys.SYSTEM.NetworkTotal[i].Name+"-Rxp", systemGenerateLineItems(sys.SYSTEM, "rxp", i)).
					AddSeries(sys.SYSTEM.NetworkTotal[i].Name+"-Txp", systemGenerateLineItems(sys.SYSTEM, "txp", i)).
					AddSeries(sys.SYSTEM.NetworkTotal[i].Name+"-Rx", systemGenerateLineItems(sys.SYSTEM, "rx", i)).
					AddSeries(sys.SYSTEM.NetworkTotal[i].Name+"-Tx", systemGenerateLineItems(sys.SYSTEM, "tx", i))
			}

			err = NetwrodCardCharts.Render(f)
			if err != nil {
				panic(err)
			}
			err = LoadCharts.Render(f)
			if err != nil {
				panic(err)
			}
		}
	}

}

func systemGenerateLineItems(data sar.SystemInformation, note string, index int) []opts.LineData {
	var value interface{}
	items := make([]opts.LineData, 0)
	if note == "rxp" {
		Ncard_items := make([]opts.LineData, 0)
		for _, Ncard := range data.NetworkTotal[index].Rxpck {
			Ncard_items = append(Ncard_items, opts.LineData{Value: Ncard})
		}
		return Ncard_items
	}
	if note == "txp" {
		Ncard_items := make([]opts.LineData, 0)
		for _, Ncard := range data.NetworkTotal[index].Txpck {
			Ncard_items = append(Ncard_items, opts.LineData{Value: Ncard})
		}
		return Ncard_items
	}
	if note == "rx" {
		Ncard_items := make([]opts.LineData, 0)
		for _, Ncard := range data.NetworkTotal[index].Rxkb {
			Ncard_items = append(Ncard_items, opts.LineData{Value: Ncard})
		}
		return Ncard_items
	}
	if note == "tx" {
		Ncard_items := make([]opts.LineData, 0)
		for _, Ncard := range data.NetworkTotal[index].Txkb {
			Ncard_items = append(Ncard_items, opts.LineData{Value: Ncard})
		}
		return Ncard_items
	}

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

func cpuGenerateLineItems(data cpu.CpuInformation, note string) []opts.LineData {
	var rate, value, core interface{}

	if note == "rate" {
		rate_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.Rate); i++ {
			rate = data.Rate[i]
			rate_items = append(rate_items, opts.LineData{Value: rate})
		}
		return rate_items
	}

	if note == "core" {
		core_items := make([]opts.LineData, 0)
		for i := 0; i < len(data.CpuCore); i++ {
			core = data.CpuCore[i]
			core_items = append(core_items, opts.LineData{Value: core})
		}
		return core_items
	}
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data.LoadAvg); i++ {
		if note == "one" {
			value = data.LoadAvg[i].AvgOne
		}
		if note == "five" {
			value = data.LoadAvg[i].AvgFive
		}
		if note == "fifteen" {
			value = data.LoadAvg[i].AvgFifteen
		}
		items = append(items, opts.LineData{Value: value})
	}
	return items
}
