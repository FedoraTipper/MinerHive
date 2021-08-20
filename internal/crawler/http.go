package crawler

import (
	"fmt"
	"github.com/FedoraTipper/AntHive/internal/constants"
	"net/http"
)

const (
	authorisationHeader = "Authorization"
)

type HTTPClient struct{}

func GetNewWwwAuthorisationHeader(baseUrl string) (http.Header, error) {
	fullURL := fmt.Sprintf("%s%s", baseUrl, constants.HomeURI)
	resp, err := MakeRequest(http.MethodGet, fullURL, "")

	if err != nil {
		return nil, err
	}

	return resp.Header, nil
}

func GetStats(baseUrl, username, password, nonceCount, hashMethod, salt string, authorisationHeaderValues AuthorisationHeaderValues) (*http.Response, error) {
	authenticator := NewAuthenticator(hashMethod, salt)

	cnonce, err := authenticator.GenerateRandomNonce()

	if err != nil {
		return nil, err
	}

	authorizationHeader, err := authenticator.GenerateAuthorizationHeader(constants.StatsURI, http.MethodGet, authorisationHeaderValues.Realm, authorisationHeaderValues.Nonce,
		nonceCount, cnonce, authorisationHeaderValues.Qop, username, password)

	if err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf("%s%s", baseUrl, constants.StatsURI)

	return MakeRequest(http.MethodGet, fullURL, authorizationHeader)
}

func MakeRequest(method, url, authorisationDigestValue string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)

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
