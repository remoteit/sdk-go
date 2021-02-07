package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type RestoreClient struct {
	endpointURL string
	token       string
	client      *http.Client
}

func NewRestoreClient(
	apiURL string, apiKey string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string,
	restoreAPIURL string,
	username string, authHash string,
) *RestoreClient {

	token := ""

	response, err := NewClient(apiURL, apiKey, apiTimeout, productVersion, os, osVersion).AuthHashSignin(
		username, authHash,
	)
	if err != nil {
		// FIXME:
	} else {
		token = response.Token
	}

	return &RestoreClient{
		endpointURL: restoreAPIURL,
		token:       token,
		client: &http.Client{
			Timeout: apiTimeout,
		},
	}
}

func (thisRef RestoreClient) Execute(deviceID string, machineID string) ([]byte, errorx.Error) {
	// build payload
	type payload struct {
		DeviceId  string `json:"deviceId"`
		MachineId string `json:"machineId"`
	}

	request := payload{
		DeviceId:  deviceID,
		MachineId: machineID,
	}
	requestAsBytes, err := json.Marshal(request)
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantPrepRequest
	}

	// build request
	req, err := http.NewRequest("POST", thisRef.endpointURL, bytes.NewBuffer(requestAsBytes))
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantPrepRequest
	}

	req.Header.Set("token", thisRef.token)

	// call API
	resp, err := thisRef.client.Do(req)
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantSendRequest
	}

	defer resp.Body.Close()

	// check error code
	if resp.StatusCode != 200 { // 200 Success
		switch resp.StatusCode {
		case 400:
			return nil, apiContracts.ErrAPI_RestoreClient_DeviceActive
		case 401:
			return nil, apiContracts.ErrAPI_RestoreClient_TokenNotSpecified
		case 403:
			return nil, apiContracts.ErrAPI_RestoreClient_DeviceNotExists
		default:
			return nil, errorx.New(apiContracts.ErrAPI_RestoreClient_Generic, resp.Status)
		}
	}

	// read response
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantReadResponse
	}

	return raw, nil
}
