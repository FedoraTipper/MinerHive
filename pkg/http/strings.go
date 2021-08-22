package http

import (
	"fmt"
	"strings"
)

func FormURL(host string, port uint) string {
	url := fmt.Sprintf("%s:%d", host, port)

	if !strings.Contains(host, "http") {
		url = "http://" + url
	}

	return url
}
