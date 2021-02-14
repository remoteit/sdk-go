package api

import (
	"encoding/json"
	"strings"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type Proxy interface {
	Create(request apiContracts.CreateProxyRequest) (apiContracts.CreateProxyResponse, errorx.Error)
	Delete(request apiContracts.DeleteProxyRequest) (apiContracts.DeleteProxyResponse, errorx.Error)
}

func NewProxy(apiClient Client) Proxy {
	return &proxy{
		apiClient: apiClient,
	}
}

type proxy struct {
	apiClient Client
}

func (thisRef proxy) Create(request apiContracts.CreateProxyRequest) (apiContracts.CreateProxyResponse, errorx.Error) {
	body, err := json.Marshal(request)
	if err != nil {
		return apiContracts.CreateProxyResponse{}, apiContracts.ErrAPI_ProxyCreate_CantPrepRequest
	}

	raw, err := thisRef.apiClient.Post("/device/connect", body)
	if err != nil {
		return apiContracts.CreateProxyResponse{}, apiContracts.ErrAPI_ProxyCreate_CantSendRequest
	}

	var response apiContracts.CreateProxyResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return apiContracts.CreateProxyResponse{}, apiContracts.ErrAPI_ProxyCreate_CantReadResponse
	}

	if response.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		if strings.Contains(response.Reason, apiContracts.API_ERROR_CODE_REASON_SERVICE_NOT_FOUND_FOR_UID) {
			return apiContracts.CreateProxyResponse{}, apiContracts.ErrAPI_ProxyCreate_NoServiceFound
		}
		if !(len(strings.TrimSpace(response.Reason)) <= 0) {
			return apiContracts.CreateProxyResponse{}, errorx.New(apiContracts.ErrAPI_ProxyCreate_Generic, response.Reason)
		}
		return apiContracts.CreateProxyResponse{}, apiContracts.ErrAPI_ProxyCreate_Unknown
	}

	return response, nil
}

func (thisRef proxy) Delete(request apiContracts.DeleteProxyRequest) (apiContracts.DeleteProxyResponse, errorx.Error) {
	body, err := json.Marshal(request)
	if err != nil {
		return apiContracts.DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantPrepRequest
	}

	raw, err := thisRef.apiClient.Post("/device/connect/stop", body)
	if err != nil {
		return apiContracts.DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantSendRequest
	}

	var response apiContracts.DeleteProxyResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return apiContracts.DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantReadResponse
	}

	if response.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		return apiContracts.DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_Unknown
	}

	return response, nil
}
