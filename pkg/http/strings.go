package http

import (
	"fmt"
	"strings"
)

func FormURL(host, port string) string {
	url := fmt.Sprintf("%s:%s", host, port)

	if !strings.Contains(host, "http") {
		url = "http://" + url
	}

	return url
}
