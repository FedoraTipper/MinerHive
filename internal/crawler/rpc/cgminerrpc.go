package rpc

import (
	"github.com/FedoraTipper/MinerHive/internal/constants"
)

type CGMinerRPCClient struct {
}

func (c *CGMinerRPCClient) GetStats(addr string) ([]byte, error) {
	request := &Request{
		Command: constants.StatsRPCCommand,
	}

	return makeCall(addr, request)
}
