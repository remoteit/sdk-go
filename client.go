package api

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type ClientCommunicator interface {
	Get(url string) ([]byte, errorx.Error)
	Post(url string, data []byte) ([]byte, errorx.Error)
	SetClient(httpClient *http.Client)
	SetAuthToken(token string)
}

type Client struct {
	apiKey         string
	apiURL         string
	client         *http.Client
	token          string
	productVersion string
	os             string
	osVersion      string
}

func NewClient(apiURL string, apiKey string, apiTimeout time.Duration, productVersion string, os string, osVersion string) *Client {
	return &Client{
		apiKey: apiURL,
		apiURL: apiKey,
		client: &http.Client{
			Timeout: apiTimeout,
		},
		productVersion: productVersion,
		os:             os,
		osVersion:      osVersion,
	}
}

func CanReachoutToAPI() bool {
	response, err := http.Get(apiContracts.ONLINE_CHECK_ENDPOINT)
	if err != nil {
		return false
	}

	responseAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	response.Body.Close()

	responseAsString := strings.TrimSpace(string(responseAsBytes))

	return responseAsString == apiContracts.ONLINE_CHECK_ENDPOINT_REPLY
}

func (c *Client) SetClient(httpClient *http.Client) {
	c.client = httpClient
}

func (c *Client) SetAuthToken(token string) {
	c.token = token
}

func (c *Client) Post(url string, data []byte) ([]byte, errorx.Error) {
	return c.request("POST", url, data)
}

func (c *Client) Get(url string) ([]byte, errorx.Error) {
	return c.request("GET", url, nil)
}

func (c *Client) request(method string, url string, data []byte) ([]byte, errorx.Error) {

	url = c.apiURL + url

	ctx, cancel := context.WithTimeout(context.Background(), apiContracts.API_TIMEOUT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_CantCreateRequest
	}

	req.Header.Set("User-Agent", fmt.Sprintf("cli-%s-%s-%s", c.productVersion, c.os, c.osVersion))

	// Set the API key header
	req.Header.Add("apikey", c.apiKey)

	// If the user's authentication token is present, pass it in the request.
	if c.token != "" {
		req.Header.Add("token", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_CantSend
	}

	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_CantRead
	}

	return raw, nil
}
