package http

import (
	"errors"
	"fmt"
	"github.com/FedoraTipper/AntHive/internal/constants"
	"io"
	"net/http"
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
