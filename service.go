package api

import (
	"encoding/json"
	"strings"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

// ServiceRegister -
type ServiceRegister struct {
	Helper ServiceRegistrationHelper
}

// ServiceRegistrationInfo -
type ServiceRegistrationInfo struct {
	Name             string
	ServiceType      string
	ServiceTypeAsInt int
	HardwareID       string
}

/*
NewService registers a new service with a given configuration.
*/
func (s *ServiceRegister) NewService(info ServiceRegistrationInfo, projectKey string, projectSecret string) (*apiContracts.Service, errorx.Error) {
	empty := new(apiContracts.Service)

	uid, err := s.Helper.GenerateUID(projectKey, projectSecret)
	if err != nil {
		return empty, err
	}

	err = s.Helper.CreateService(uid, info.ServiceType)
	if err != nil {
		return empty, err
	}

	hardwareID := uid
	if info.HardwareID != "" {
		hardwareID = info.HardwareID
	}
	secret, err := s.Helper.RegisterService(info.Name, uid, hardwareID, info.ServiceType, info.ServiceTypeAsInt)
	if err != nil {
		return empty, err
	}

	overload := 0
	if info.ServiceTypeAsInt == apiContracts.BulkServiceID {
		overload = apiContracts.MultiPortServiceID
	}

	return &apiContracts.Service{
		HardwareID: hardwareID,
		Overload:   overload,
		Secret:     secret,
		Type:       info.ServiceTypeAsInt,
		UID:        uid,
	}, nil
}

/*
RemoveService removes a given service by its UID.
*/
func (c *Client) RemoveService(uid string) errorx.Error {
	// Construct the request data to send to the API
	type request struct {
		UID string `json:"deviceaddress"`
	}
	data := request{uid}
	body, err := json.Marshal(data)
	if err != nil {
		return apiContracts.ErrAPI_Service_CantPrepRequest
	}

	// Make the API request
	raw, errx := c.Post("/device/delete", body)
	if errx != nil {
		return errx
	}

	// Parse the JSON response into a usable struct.
	type response struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	var resp response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return apiContracts.ErrAPI_Service_CantReadResponse
	}

	if resp.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_SERVICE_NOT_FOUND_FOR_UID) {
			return apiContracts.ErrAPI_Service_NoServiceFound
		}
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return errorx.New(apiContracts.ErrAPI_Service_Generic, resp.Reason)
		}
		return apiContracts.ErrAPI_Service_Unknown
	}

	return nil
}
