package constants

type MinerSeries string

const (
	X19 MinerSeries = "X19"
)

// TODO Move this somewhere else
var (
	SupporteredMiners = map[MinerSeries]struct{}{
		X19: {},
	}
)

type MinerStatus int

const (
	MinerOK MinerStatus = iota + 1
	MinerStopped
)

const (
	S19_OKStatus string = "S"
)
