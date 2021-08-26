package rpc

import "github.com/FedoraTipper/AntHive/internal/constants"

type S19RPCClient struct {
}

func (c *S19RPCClient) GetStats(addr string) ([]byte, error) {
	request := &Request{
		Command: constants.StatsRPCCommand,
	}

	return makeCall(addr, request)
}
