package constants

type MinerSeries string

const (
	X19 MinerSeries = "X19"
)

// TODO Move this somewhere else
var (
	SupportedMiners = map[MinerSeries]struct{}{
		X19: {},
	}
)

const (
	CGMiner_STATS           = "STATS"
	CGMINER_RESPONSE_STATUS = ""
)

const (
	S19_OKStatus          string = "S"
	Antminer_OKChipStatus string = "o"
)
