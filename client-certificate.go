package api

import (
	"encoding/json"
	"net/http"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type CertificateClient interface {
	Generate(request apiContracts.CertificateRequest) (*apiContracts.CertificateResponse, errorx.Error)
}

type certificateClient struct {
	apiURL     string
	apiToken   string
	apiTimeout time.Duration
	userAgent  string
}

func NewCertificateClient(apiURL string, apiToken string, apiTimeout time.Duration, userAgent string) CertificateClient {
	return &certificateClient{
		apiURL:     apiURL,
		apiToken:   apiToken,
		apiTimeout: apiTimeout,
		userAgent:  userAgent,
	}
}

func (thisRef certificateClient) Generate(request apiContracts.CertificateRequest) (*apiContracts.CertificateResponse, errorx.Error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, apiContracts.ErrAPI_RestoreClient_CantPrepRequest
	}

	headers := map[string]string{
		"User-Agent": thisRef.userAgent,
		"token":      thisRef.apiToken,
	}

	httpResponse, data, err := doHTTPRequest(http.MethodPost, headers, thisRef.apiURL, payload, thisRef.apiTimeout)
	if err != nil {
		return nil, errorx.NewFromErr(apiContracts.ErrAPI_CertClient_Generic, err)
	}

	if httpResponse.StatusCode != 200 {
		switch httpResponse.StatusCode {
		case 401:
			return nil, apiContracts.ErrAPI_CertClient_TokenNotSpecified
		default:
			return nil, errorx.New(apiContracts.ErrAPI_CertClient_Generic, httpResponse.Status)
		}
	}

	//
	response := apiContracts.CertificateResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_Error
	}

	return &response, nil
}
