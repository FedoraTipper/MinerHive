package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/FedoraTipper/AntHive/internal/constants"
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
		return &S19RPCClient{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown model: %s", model))
	}
}

func makeCall(addr string, request *Request) ([]byte, error) {
	client, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	blob, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, err = client.Write(blob)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(client)

	if err != nil {
		return nil, err
	}

	b = bytes.Trim(b, "\x00")

	return b, nil
}
