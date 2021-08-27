package config

import (
	"errors"
	"fmt"
)

type HiveConfig struct {
	LoggingLevel string
	LoggingFile  string
	Port         uint
	Token        string
	Redis        RedisConfig
}

func (h *HiveConfig) Validate() []error {
	var errs []error

	if len(h.Token) == 0 {
		errs = append(errs, errors.New("Value for Token is empty"))
	}

	redisErrs := h.Redis.Validate()

	for _, err := range redisErrs {
		errs = append(errs, errors.New(fmt.Sprintf("RedisConfig error (%v)", err)))
	}

	return errs
}
