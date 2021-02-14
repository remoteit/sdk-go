package api

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func isNullOrEmpty(value string) bool {
	return len(strings.TrimSpace(value)) <= 0
}

func doHTTPRequest(method string, headers map[string]string, url string, payload []byte, timeout time.Duration) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	return response, data, err
}
