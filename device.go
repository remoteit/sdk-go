package api

import (
	"encoding/json"
	"fmt"
	"strings"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type Device interface {
	Unregister(uid string) errorx.Error
	Transfer(uid string, destinationAccount string) errorx.Error
	ListAll() (apiContracts.DeviceListAllResponse, errorx.Error)
}

func NewDevice(apiClient Client) Device {
	return &device{
		apiClient: apiClient,
	}
}

type device struct {
	apiClient Client
}

func (thisRef device) Unregister(uid string) errorx.Error {
	type request struct{}
	data := request{}
	body, err := json.Marshal(data)
	if err != nil {
		return apiContracts.ErrAPI_Device_CantPrepRequest
	}

	raw, err := thisRef.apiClient.Post(fmt.Sprintf("/developer/device/delete/registered/%s", uid), body)
	if err != nil {
		return apiContracts.ErrAPI_Device_CantSendRequest
	}

	type response struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	var resp response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return apiContracts.ErrAPI_Device_CantReadResponse
	}

	if resp.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_SERVICE_NOT_FOUND_FOR_UID) {
			return apiContracts.ErrAPI_Device_NoServiceFound
		}

		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return errorx.New(apiContracts.ErrAPI_Device_Generic, resp.Reason)
		}

		return apiContracts.ErrAPI_Device_Unknown
	}

	return nil
}

func (thisRef device) Transfer(uid string, destinationAccount string) errorx.Error {
	type request struct {
		User string `json:"user"`
	}
	data := request{
		User: destinationAccount,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return apiContracts.ErrAPI_Device_CantPrepRequest
	}

	raw, err := thisRef.apiClient.Post(fmt.Sprintf("/developer/devices/transfer/%s", uid), body)
	if err != nil {
		return apiContracts.ErrAPI_Device_CantSendRequest
	}

	type response struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	var resp response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return apiContracts.ErrAPI_Device_CantReadResponse
	}

	if resp.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return errorx.New(apiContracts.ErrAPI_Device_Generic, resp.Reason)
		}

		return apiContracts.ErrAPI_Device_Unknown
	}

	return nil
}

func (thisRef device) ListAll() (apiContracts.DeviceListAllResponse, errorx.Error) {
	raw, err := thisRef.apiClient.Get("/device/list/all?cache=false")
	if err != nil {
		return apiContracts.DeviceListAllResponse{}, apiContracts.ErrAPI_DeviceList_CantSendRequest
	}

	var response apiContracts.DeviceListAllResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return apiContracts.DeviceListAllResponse{}, apiContracts.ErrAPI_DeviceList_CantReadResponse
	}

	if response.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		return apiContracts.DeviceListAllResponse{}, errorx.New(apiContracts.ErrAPI_DeviceList_Generic, response.Reason)
	}

	return response, nil
}
