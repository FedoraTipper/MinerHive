package crawler

import (
	"errors"
	"fmt"
	"github.com/FedoraTipper/AntHive/internal/constants"
	"github.com/FedoraTipper/AntHive/pkg/hash"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	authorizationHeaderSkeleton = `Digest realm="%s", nonce="%s", uri="%s", response="%s", qop=%s, nc=%s, cnonce="%s", username="%s"`
	authenticateHeaderKey       = "Www-Authenticate"

	qopAuth     = "auth"
	qopAuthInt  = "auth-int"
	qopAuthConf = "auth-conf"
)

type Authenticator struct {
	hasher                    hash.Hasher
	salt                      string
	authorisationHeaderValues AuthorisationHeaderValues
}

func NewAuthenticator(hashMethod, salt string) *Authenticator {
	hasher := hash.GetHasher(hashMethod)

	return &Authenticator{
		hasher: hasher,
		salt:   salt,
	}
}

type AuthorisationHeaderValues struct {
	Qop   string
	Nonce string
	Realm string
}

func ExtractWwwAuthenticateHeader(header http.Header) AuthorisationHeaderValues {
	headerValue := header.Get(authenticateHeaderKey)
	nonce := filterNonceValue(headerValue)
	realm := filterRealm(headerValue)
	qop := filterQOP(headerValue)

	return AuthorisationHeaderValues{
		Qop:   qop,
		Realm: realm,
		Nonce: nonce,
	}
}

func filterNonceValue(headerValue string) string {
	return regexFind(`nonce=".*?"`, headerValue, true)
}

func filterRealm(headerValue string) string {
	return regexFind(`realm=".*?"`, headerValue, true)
}

func filterQOP(headerValue string) string {
	return regexFind(`qop=".*?"`, headerValue, true)
}

func regexFind(regex, header string, quoteClean bool) string {
	re := regexp.MustCompile(regex)

	matchingKeyPair := re.FindString(header)

	value := strings.Split(matchingKeyPair, "=")[1]

	if quoteClean {
		return strings.ReplaceAll(value, `"`, "")
	}

	return value
}

func (a *Authenticator) GenerateAuthorizationHeader(uri constants.S19URI, httpMethod, realm, nonce, nonceCount, cnonce, qop, username, password string) (string, error) {
	A1 := a.generateA1(username, realm, password)
	A2, err := a.generateA2(httpMethod, uri, qop)

	if err != nil {
		return "", err
	}

	response := a.generateResponse(A1, A2, nonce, nonceCount, cnonce, qop)

	return fmt.Sprintf(authorizationHeaderSkeleton, realm, nonce, uri, response, qop, nonceCount, cnonce, username), nil
}

// FOLLOW RFC7616 SPEC
// https://datatracker.ietf.org/doc/html/rfc7616#section-2.1

func (a *Authenticator) generateResponse(A1, A2, nonce, nonceCount, cnonce, qop string) string {
	data := []byte(fmt.Sprintf("%s:%s:%s:%s:%s:%s", A1, nonce, nonceCount, cnonce, qop, A2))
	return fmt.Sprintf("%x", a.hasher.Hash(data))
}

func (a *Authenticator) generateA1(username, realm, password string) string {
	data := []byte(fmt.Sprintf("%s:%s:%s", username, realm, password))
	return fmt.Sprintf("%x", a.hasher.Hash(data))
}

func (a *Authenticator) generateA2(httpMethod string, uri constants.S19URI, qop string) (string, error) {
	var data string

	switch qop {
	case qopAuth:
		data = fmt.Sprintf("%s:%s", httpMethod, uri)
		break
	case qopAuthInt, qopAuthConf:
		return "", errors.New(fmt.Sprintf("Authentication method '%s' has not been implemented", qop))
	default:
		return "", errors.New(fmt.Sprintf("QOP authentication %s method unknown", qop))
	}

	return fmt.Sprintf("%x", a.hasher.Hash([]byte(data))), nil
}

func (a *Authenticator) GenerateRandomNonce() (string, error) {
	rand.Seed(time.Now().Unix())

	var data [][]byte

	if len(a.salt) > 0 {
		data = append(data, []byte(a.salt))
	}

	blk := make([]byte, 32)
	_, err := rand.Read(blk)

	if err != nil {
		return "", err
	}

	data = append(data, blk)

	return fmt.Sprintf("%x", a.hasher.Hash(data...)), nil
}
