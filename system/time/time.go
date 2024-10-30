package time

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func stringToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("转换错误:", err)
	} else {
		return num
	}
	return 0
}

func GetTimeSync(_t *TimeCheckInformation, sync_status string, sync_line int) {
	file_buf, err := os.Open(sync_status)
	if err != nil {
		fmt.Println("open file failed:", err)
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)

	get_status := regexp.MustCompile(`eth.*:(\w+).*?eth.*:(\w+).*?phc2sys_status:(\w+_\w+)`)
	currentUser, _ := os.Hostname()
	_time := make([]string, 0)

	for i := 0; i < sync_line; i++ {
		line, _ := reader.ReadString('\n')
		if strings.Contains(line, "phc2sys_status") {
			_time = append(_time, strings.Split(line, currentUser)[0])
			status_ := get_status.FindStringSubmatch(line)
			_t.CardSyncState = append(_t.CardSyncState, NetCardState{
				Eth1SyncState: ptp_status[status_[1]],
				Eth2SyncState: ptp_status[status_[2]],
				LockState:     ptp_status[status_[3]],
			})
		}
	}
	_t.Time = _time
}

func GetTimeMonitor(_t *TimeCheckInformation, offset_status string, offset_line int) {
	file_buf, err := os.Open(offset_status)
	if err != nil {
		fmt.Println("open file failed:", err)
	}
	defer file_buf.Close()
	reader := bufio.NewReader(file_buf)

	get_monitor := regexp.MustCompile(`phc:.*`)
	currentUser, _ := os.Hostname()

	for i := 0; i < offset_line; i++ {
		line, err := reader.ReadString('\n')
		if strings.Contains(line, "phc:") {
			monitor := get_monitor.FindStringSubmatch(line)
			phc_offset := strings.Split(strings.Split(monitor[0], "(")[2], ")")[0]
			sys_offset := strings.Split(strings.Split(monitor[0], "(")[5], ")")[0]
			eth1_phctomaster := stringToInt64(strings.Split(phc_offset, "/")[0])
			eth2_phctomaster := stringToInt64(strings.Split(phc_offset, "/")[1])
			eth1_systemtophc := stringToInt64(strings.Split(sys_offset, "/")[0])
			eth2_systemtophc := stringToInt64(strings.Split(sys_offset, "/")[1])
			_t.Eth1TimeOffset = append(_t.Eth1TimeOffset, OffsetValue{
				Time:                 strings.Split(line, currentUser)[0],
				PhcToMasterOffset:    eth1_phctomaster,
				SystemToPhcOffset:    eth1_systemtophc,
				SystemToMasterOffset: eth1_phctomaster + eth1_systemtophc,
			})
			_t.Eth2TimeOffset = append(_t.Eth2TimeOffset, OffsetValue{
				Time:                 strings.Split(line, currentUser)[0],
				PhcToMasterOffset:    eth2_phctomaster,
				SystemToPhcOffset:    eth2_systemtophc,
				SystemToMasterOffset: eth2_phctomaster + eth2_systemtophc,
			})
		}
		if err != nil {
			break
		}
	}
}
func GetTimeInformation(sync_status, offset_status string, sync_count, offset_count int) TimeCheckInformation {
	var t TimeCheckInformation
	GetTimeSync(&t, sync_status, sync_count)
	GetTimeMonitor(&t, offset_status, offset_count)
	return t
}
