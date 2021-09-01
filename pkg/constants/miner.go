package constants

type MinerStatus int

const (
	MinerStatusOk MinerStatus = iota + 1
	MinerStatusStopped
	MinerStatusUnknown
)
