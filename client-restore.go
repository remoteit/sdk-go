package api

import (
	"encoding/json"
	"net/http"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type RestoreClient interface {
	Restore(deviceID string, machineID string) ([]byte, errorx.Error)
}

type restoreClient struct {
	apiURL     string
	apiToken   string
	apiTimeout time.Duration
	userAgent  string
}

func NewRestoreClient(apiURL string, apiToken string, apiTimeout time.Duration, userAgent string) RestoreClient {
	return &restoreClient{
		apiURL:     apiURL,
		apiToken:   apiToken,
		apiTimeout: apiTimeout,
		userAgent:  userAgent,
	}
}

func (thisRef restoreClient) Restore(deviceID string, machineID string) ([]byte, errorx.Error) {
	type payloadT struct {
		DeviceId  string `json:"deviceId"`
		MachineId string `json:"machineId"`
	}

	request := payloadT{
		DeviceId:  deviceID,
		MachineId: machineID,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantPrepRequest
	}

	headers := map[string]string{
		"User-Agent": thisRef.userAgent,
		"token":      thisRef.apiToken,
	}

	response, data, err := doHTTPRequest(http.MethodPost, headers, thisRef.apiURL, payload, thisRef.apiTimeout)
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_Error
	}

	if response.StatusCode != 200 {
		switch response.StatusCode {
		case 400:
			return nil, apiContracts.ErrAPI_RestoreClient_DeviceActive
		case 401:
			return nil, apiContracts.ErrAPI_RestoreClient_TokenNotSpecified
		case 403:
			return nil, apiContracts.ErrAPI_RestoreClient_DeviceNotExists
		default:
			return nil, errorx.New(apiContracts.ErrAPI_RestoreClient_Generic, response.Status)
		}
	}

	return data, nil
}
