package api

import (
	"encoding/json"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

// DeviceListAllRequest -
//
// SAMPLE
// curl 'https://api.remot3.it/apv/v27/device/list/all?cache=false'
// -H 'developerKey: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
// -H 'token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
//

// Device -
type Device struct {
	DeviceAddress string `json:"deviceaddress,omitempty"`
	DeviceType    string `json:"devicetype,omitempty"`
	DeviceAlias   string `json:"devicealias,omitempty"`
	OwnerUserName string `json:"ownerusername,omitempty"`
	Scripting     bool   `json:"scripting,omitempty"`
}

// DeviceListAllResponse -
type DeviceListAllResponse struct {
	Status  string   `json:"status,omitempty"`
	Reason  string   `json:"reason"`
	Devices []Device `json:"devices,omitempty"`
}

// DeviceListAll - lists all deivces for the user
func (c *Client) DeviceListAll() (DeviceListAllResponse, errorx.Error) {
	raw, err := c.Get("/device/list/all?cache=false")
	if err != nil {
		return DeviceListAllResponse{}, apiContracts.ErrAPI_DeviceList_CantSendRequest
	}

	var response DeviceListAllResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return DeviceListAllResponse{}, apiContracts.ErrAPI_DeviceList_CantReadResponse
	}

	if response.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		return DeviceListAllResponse{}, errorx.New(apiContracts.ErrAPI_DeviceList_Generic, response.Reason)
	}

	return response, nil
}
