# qomolo-system-analysis

> analysis system load status

## Description
```bash
采用ubuntu service方式进行监控
基于系统监控日志的可视化工具：
cpu: cpustat
    description:
        cpu process monitor & system cpu core & cpu rate
    cmd: 
        cpustat 5 1 -l -n 10 -x 
io: pidstat
    description：
        system input/outpu monitor
    cmd:
        sudo pidstat -d 10 1 | grep -v '^Average' |  sed '1,3 d' | sort -rnk 5 | awk '{if ($4>500 || $5>500) print $0}'
mem: pidstat
    description:
        system mem Percentage 
    cmd: 
        sudo pidstat -r 30 1 | grep -v '^Average' | sed '1,3 d' |sort -rnk 8 | awk '{if($8 > 1) printf "%-15s%-15s%-15s%-15s%-15s%-15s%-15s%-15s\n", $3,$4,$5,$6,$7,$8,$9,$10}'
system: sar
    description:
        cpu total use & network card total use
    cmd:
        sudo sar -r -u -n DEV 10
```

## Install 
```bash
sudo apt update 
sudo apt install qomolo-system-analysis -y
bash   # 刷新环境变量
```

## Usage 
```bash
args1:
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
```

## Remove
```bash
sudo apt remove qomolo-system-analysis -y 
```

## Build 
```bash
# 移除CGO依赖，增加兼容性
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main main.go

# 交叉编译 arm版本
CGO_ENABLED=0 GOOS=linux GOARCH=arm  go build -o main main.go
```# SystemResourceViewAnalysis
