package http

// DEPRECATED. REPLACED WITH RPC CLIENT

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FedoraTipper/MinerHive/internal/constants"
)

const (
	authorisationHeader = "Authorization"
)

type IHTTPClient interface {
	GetStatsResponse(baseUrl, username, password, salt string) (*http.Response, error)
}

func GetHTTPClient(model constants.MinerSeries) (IHTTPClient, error) {
	switch model {
	case constants.X19:
		return &S19HTTPClient{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown model: %s", model))
	}
}

func MakeRequest(method, url, authorisationDigestValue string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	if len(authorisationDigestValue) > 0 {
		req.Header.Set(authorisationHeader, authorisationDigestValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	client.CloseIdleConnections()

	return resp, nil
}
