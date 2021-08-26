package models

import (
	"encoding/json"

	"github.com/FedoraTipper/AntHive/pkg/constants"
)

type Miner struct {
	MinerName    string
	CrawlerName  string
	MinerModel   string
	Status       constants.MinerStatus
	MinerVersion string
	Uptime       int
	HashRateUnit string
	Fans         []Fan
	HashBoards   []Hashboard
}

type Fan struct {
	FanNumber int
	RPM       int
}

type Hashboard struct {
	BoardNumber             int
	NoOfChips               int
	MalfunctioningChipsList []int
	HWErrors                int
	ChipFrequency           int
	CurrentHashRate         float64
	RatedHashRate           float64
	PICTemperature          []int
	PCBTemperature          []int
	ChipTemperature         []int
	SerialNumber            string
}

func (m *Miner) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Miner) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return nil
}
