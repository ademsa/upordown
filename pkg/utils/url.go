package utils

import (
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

// CheckTargetURL to send GET request to endpoint and return status code and body
func CheckTargetURL(url string) (string, int, int64, *html.Node, error) {
	prefixSplit := strings.Split(url, "://")
	prefix := "http://"
	if len(prefixSplit) == 1 {
		url = prefix + url
	} else {
		prefix = prefixSplit[0] + "://"
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	defer httpClient.CloseIdleConnections()

	st := time.Now()
	response, err := httpClient.Get(url)
	et := time.Now()
	if err != nil {
		return "", 0, time.Duration(0).Milliseconds(), nil, err
	}
	defer response.Body.Close()

	rt := et.Sub(st).Milliseconds()
	baseURL := prefix + response.Request.Host

	body, err := html.Parse(response.Body)
	if err != nil {
		return baseURL, response.StatusCode, rt, nil, err
	}

	return baseURL, response.StatusCode, rt, body, nil
}
