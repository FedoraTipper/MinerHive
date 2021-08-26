package rpc

import "github.com/FedoraTipper/AntHive/internal/constants"

type CGMinerRPCClient struct {
}

func (c *CGMinerRPCClient) GetStats(addr string) ([]byte, error) {
	request := &Request{
		Command: constants.StatsRPCCommand,
	}

	return makeCall(addr, request)
}
