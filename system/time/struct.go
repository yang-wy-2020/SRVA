package time

type OffsetValue struct {
	PhcToMasterOffset    int
	SystemToPhcOffset    int
	SystemToMasterOffset int
}

type TimeCheckInformation struct {
	LockState  string
	TimeOffset []OffsetValue
}
