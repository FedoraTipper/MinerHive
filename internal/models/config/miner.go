package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/FedoraTipper/MinerHive/internal/constants"
)

type MinerConfig struct {
	MinerName string
	Host      string
	Port      uint
	Model     constants.MinerSeries
}

var (
	minerNames = map[string]int{}
)

func (m *MinerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func (m *MinerConfig) Validate() []error {
	var errs []error

	if len(m.MinerName) == 0 {
		errs = append(errs, errors.New("Value for MinerName is empty"))
	} else if _, exists := minerNames[m.MinerName]; exists {
		errs = append(errs, errors.New(fmt.Sprintf("Duplicate MinerName %s found. Each miner should have an unique MinerName.", m.MinerName)))
	} else {
		minerNames[m.MinerName] = 1
	}

	if m.Port == 0 {
		errs = append(errs, errors.New("Value for Port is unassigned"))
	}

	if len(errs) == 0 {
		conn, err := net.Dial("tcp", m.GetAddress())
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Unable to dial miner (%s)", m.GetAddress())))
		} else {
			conn.Close()
		}
	}

	if len(m.Model) == 0 {
		errs = append(errs, errors.New("Value for Model is empty"))
	} else if _, exists := constants.SupportedMiners[m.Model]; !exists {
		errs = append(errs, errors.New(fmt.Sprintf("Model %v is not supported", m.Model)))
	}

	return errs
}
