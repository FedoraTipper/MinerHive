package http

import (
	"fmt"
	"github.com/FedoraTipper/AntHive/internal/constants"
	"github.com/FedoraTipper/AntHive/internal/crawler/authentication"
	"github.com/FedoraTipper/AntHive/pkg/hash"
	"github.com/FedoraTipper/AntHive/pkg/hex"
	"log"
	"net/http"
)

type S19HTTPClient struct{}

func (c *S19HTTPClient) GetStatsResponse(baseUrl, username, password, salt string) (*http.Response, error) {
	wwwAuthorisationHeader, err := c.getNewWwwAuthorisationHeader(baseUrl)

	if err != nil {
		log.Fatalln(err)
	}

	authenticator := authentication.NewS19Authenticator(hash.Md5, salt, wwwAuthorisationHeader)

	authorizationHeader, err := authenticator.GenerateAuthorizationHeader(constants.StatsURI, http.MethodGet, hex.GenerateHexString(8), username, password)

	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s%s", baseUrl, constants.StatsURI)

	return MakeRequest(http.MethodGet, fullURL, authorizationHeader, nil)
}

func (c *S19HTTPClient) getNewWwwAuthorisationHeader(baseUrl string) (http.Header, error) {
	fullURL := fmt.Sprintf("%s%s", baseUrl, constants.HomeURI)
	resp, err := MakeRequest(http.MethodGet, fullURL, "", nil)

	if err != nil {
		return nil, err
	}

	return resp.Header, nil
}
