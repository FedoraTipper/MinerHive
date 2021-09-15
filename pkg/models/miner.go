package models

import (
	"encoding/json"

	"github.com/FedoraTipper/MinerHive/pkg/constants"
)

type MinerStats struct {
	MinerName    string
	CrawlerName  string
	MinerModel   string
	Status       constants.MinerStatus
	MinerVersion string
	Uptime       int
	HashRateUnit string
	Fans         []Fan
	HashBoards   []Hashboard
	Nonce        string // Generated Nonce to avoid consuming the same obj in redis
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
	PICTemperature          []int
	PCBTemperature          []int
	ChipTemperature         []int
}

func (m *MinerStats) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MinerStats) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return nil
}
