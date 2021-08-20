package constants

type Miner struct {
	MinerName       string
	MinerSerial     string
	CrawlerName     string
	MinerModel      string
	Status          string
	MinerVersion    string
	CurrentHashRate float64
	HashRateAvg30   float64
	Fans            []Fan
	HashBoards      []Hashboard
}

type Fan struct {
	FanNumber int
	RPM       int
}

type Hashboard struct {
	BoardNumber     int
	NoOfChips       int
	Errors          int
	ChipFrequency   int
	CurrentHashRate float64
	RatedHashRate   float64
	PICTemperature  []int
	PCBTemperature  []int
	ChipTemperature []int
}
