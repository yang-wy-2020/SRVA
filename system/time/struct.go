package time

var (
	ptp_status = map[string]int{
		"init":         0,
		"gm_check":     1,
		"wait_to_lock": 2,
		"restart":      3,
		"lock_bad":     4,
		"lock_good":    5,
		"off_off":      11,
		"on_off":       12,
		"off_on":       13,
		"locked_off":   21,
		"off_locked":   22,
	}
)

type OffsetValue struct {
	Time                 string
	PhcToMasterOffset    int64
	SystemToPhcOffset    int64
	SystemToMasterOffset int64
}
type NetCardState struct {
	Eth1SyncState int
	Eth2SyncState int
	LockState     int
}

type TimeCheckInformation struct {
	Time                           []string
	CardSyncState                  []NetCardState
	Eth1TimeOffset, Eth2TimeOffset []OffsetValue
}
