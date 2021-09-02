package config

import (
	"errors"
	"fmt"
	"net"
)

type RedisConfig struct {
	Host             string
	Port             uint
	Username         string
	Password         string
	SelectedDatabase int
}

func (r *RedisConfig) Validate() []error {
	var errs []error

	if len(r.Host) == 0 {
		errs = append(errs, errors.New("Value for Host is empty"))
	}

	if r.Port == 0 {
		errs = append(errs, errors.New("Value for Port is unassigned"))
	}

	if len(errs) == 0 {
		conn, err := net.Dial("tcp", r.GetAddress())
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Unable to dial RedisDB (%s)", r.GetAddress())))
		} else {
			conn.Close()
		}
	}

	return errs
}

func (r *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
