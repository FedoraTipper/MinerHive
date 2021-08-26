package config

import (
	"errors"
	"fmt"

	"github.com/FedoraTipper/AntHive/internal/constants"
)

type CrawlerConfig struct {
	CrawlerName    string
	CrawlInterval  int
	LoggingEnabled bool
	LoggingLevel   string
	Salt           string
	Miners         []MinerConfig
	Redis          RedisConfig
}

func (c *CrawlerConfig) Validate() []error {
	var errs []error

	if len(c.CrawlerName) == 0 {
		errs = append(errs, errors.New("Value for CrawlerName is empty"))
	}

	for i, miner := range c.Miners {
		minerErrs := miner.Validate()
		for _, err := range minerErrs {
			errs = append(errs, errors.New(fmt.Sprintf("MinerConfig error for Miner %d with error (%v)", i+1, err)))
		}
	}

	redisErrs := c.Redis.Validate()

	for _, err := range redisErrs {
		errs = append(errs, errors.New(fmt.Sprintf("RedisConfig error (%v)", err)))
	}

	return errs
}

type MinerConfig struct {
	MinerName string
	Host      string
	Port      uint
	Model     constants.MinerSeries
}

var (
	minerNames = map[string]int{}
)

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

	if len(m.Model) == 0 {
		errs = append(errs, errors.New("Value for Model is empty"))
	} else if _, exists := constants.SupportedMiners[m.Model]; !exists {
		errs = append(errs, errors.New(fmt.Sprintf("Model %v is not supported", m.Model)))
	}

	return errs
}
