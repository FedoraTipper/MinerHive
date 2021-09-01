package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/FedoraTipper/AntHive/internal/constants"
	"go.uber.org/zap"
)

type IRPCCLient interface {
	GetStats(addr string) ([]byte, error)
}

type Request struct {
	Command   string `json:"command"`
	Parameter string `json:"parameter,omitempty"`
}

func GetRPCClient(model constants.MinerSeries) (IRPCCLient, error) {
	switch model {
	case constants.X19:
		return &CGMinerRPCClient{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown model: %s", model))
	}
}

func makeCall(addr string, request *Request) ([]byte, error) {
	zap.S().Debugw("Dialling RPC over tcp", "addr", addr)
	client, err := net.Dial("tcp", addr)

	if err != nil {
		zap.S().Errorw("Unable to dial RPC API", "addr", addr, "Error", err)
		return nil, err
	}

	defer func() {
		err := client.Close()
		if err != nil {
			zap.S().Errorw(fmt.Sprintf("Unable to close client connection to %s", addr), "Error", err)
		}
	}()

	jsonRequest, err := json.Marshal(request)

	if err != nil {
		zap.S().Errorw("Unable to marshal request to make RPC call", "request", request, "Error", err)
		return nil, err
	}

	zap.S().Debugw("Successfully marshaled RPC request to json", "JSON Request", jsonRequest)

	_, err = client.Write(jsonRequest)

	if err != nil {
		zap.S().Errorw("Unable to write to RPC connection", "JSON Request", jsonRequest, "Error", err)
		return nil, err
	}

	zap.S().Debug("Successfully wrote json data to RPC API")

	b, err := ioutil.ReadAll(client)

	if err != nil {
		zap.S().Errorw("Unable to read response from RPC API", "Error", err)
		return nil, err
	}

	zap.S().Debug("Successfully read RPC response")

	b = bytes.Trim(b, "\x00")

	return b, nil
}
