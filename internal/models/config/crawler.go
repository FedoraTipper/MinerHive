package config

import (
	"errors"
	"fmt"
)

type CrawlerConfig struct {
	CrawlerName   string
	CrawlInterval int
	LoggingLevel  string
	LoggingFile   string
	Miners        []MinerConfig
	Redis         RedisConfig
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
