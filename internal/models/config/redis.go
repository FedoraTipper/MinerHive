package config

import "errors"

type RedisConfig struct {
	Host             string
	Port             string
	Username         string
	Password         string
	SelectedDatabase int
}

func (r *RedisConfig) Validate() []error {
	var errs []error

	if len(r.Host) == 0 {
		errs = append(errs, errors.New("Value for Host is empty"))
	}

	if len(r.Port) == 0 {
		errs = append(errs, errors.New("Value for Port is empty"))
	}

	return errs
}
