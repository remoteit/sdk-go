package api

import (
	"encoding/json"
	"fmt"
	"strings"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

// ServiceRegistrationHelper -
type ServiceRegistrationHelper interface {
	GenerateUID(projectKey string, projectSecret string) (uid string, err errorx.Error)
	CreateService(uid string, serviceType string) errorx.Error
	RegisterService(name string, uid string, hardwareID string, serviceType string, serviceTypeAsInt int) (secret string, err errorx.Error)
}

// ServiceRegistrationHelp -
type ServiceRegistrationHelp struct {
	Client ClientCommunicator
}

func (s *ServiceRegistrationHelp) GenerateUID(projectKey string, projectSecret string) (string, errorx.Error) {
	// Send the API request
	raw, errx := s.Client.Get(fmt.Sprintf("/device/address/%s/%s", projectKey, projectSecret))
	if errx != nil {
		return "", errx
	}

	// Parse the JSON response into a usable struct.
	type response struct {
		UID    string `json:"deviceaddress"`
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	var resp response
	err := json.Unmarshal(raw, &resp)
	if err != nil {
		return "", apiContracts.ErrAPI_Service_CantPrepRequest
	}

	// Handle situations where the API returns an error and pass the reason
	// back to the consumer.
	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return "", errorx.New(apiContracts.ErrAPI_Helpers_Generic, resp.Reason)
		}
		return "", apiContracts.ErrAPI_Helpers_CantCreateServiceID
	}

	// For some reason, the API returned no UID in the request. This is unlikely
	// to happen but we handle it anyways.
	if resp.UID == "" {
		return "", apiContracts.ErrAPI_Helpers_NoUIDReturned
	}

	return resp.UID, nil
}

/*
CreateService creates a blank service to be configured by the API. This is the second
part of "traditional" device registration.
*/
func (s *ServiceRegistrationHelp) CreateService(uid string, serviceType string) errorx.Error {
	// Construct the request data to send to the API
	type requestBody struct {
		UID         string `json:"deviceaddress"`
		ServiceType string `json:"devicetype"`
	}
	rawBody := requestBody{uid, serviceType}
	body, err := json.Marshal(rawBody)
	if err != nil {
		return apiContracts.ErrAPI_Service_CantPrepRequest
	}

	// Attempt to create the device via the API
	raw, err := s.Client.Post("/device/create", body)
	if err != nil {
		return apiContracts.ErrAPI_Service_CantSendRequest
	}

	// Parse the JSON response into a usable struct.
	type response struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	var resp response
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return apiContracts.ErrAPI_Service_CantReadResponse
	}

	// Handle a variety of possible error conditions
	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_API_TOKEN) {
			return apiContracts.ErrAPI_Helpers_NoToken
		}
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return errorx.New(apiContracts.ErrAPI_Helpers_Generic, resp.Reason)
		}
		return apiContracts.ErrAPI_Helpers_Unknown
	}

	return nil
}

/*
RegisterService associates a service with a UID to create a service. This is the third
part of "traditional" device registration.
*/
func (s *ServiceRegistrationHelp) RegisterService(
	name string,
	uid string,
	hardwareID string,
	serviceType string,
	serviceTypeAsInt int,
) (string, errorx.Error) {
	// Construct the request data to send to the API
	type requestBody struct {
		UID         string `json:"deviceaddress"`
		ServiceType string `json:"devicetype"`
		Name        string `json:"devicealias"`
		HardwareID  string `json:"hardwareid"`
		SkipSecret  string `json:"skipsecret"`
		SkipEmail   string `json:"skipemail"`
	}
	skipEmail := "true"
	if serviceTypeAsInt == apiContracts.BulkServiceID {
		skipEmail = "false"
	}
	rawBody := requestBody{
		UID:         uid,
		ServiceType: serviceType,
		Name:        name,
		HardwareID:  hardwareID,
		SkipSecret:  "true",
		SkipEmail:   skipEmail,
	}
	body, err := json.Marshal(rawBody)
	if err != nil {
		return "", apiContracts.ErrAPI_Service_CantPrepRequest
	}

	// Attempt to create the device via the API
	raw, errx := s.Client.Post("/device/register", body)
	if errx != nil {
		return "", errx
	}

	// Parse the JSON response into a usable struct.
	type registeredDevice struct {
		UID          string `json:"deviceaddress"`
		Name         string `json:"name"`
		Created      string `json:"created"`
		Owner        string `json:"owner"`
		Enabled      string `json:"enabled"`
		Alerted      string `json:"alerted"`
		Title        string `json:"title"`
		Type         string `json:"devicetype"`
		Manufacturer string `json:"manufacturer"`
		Region       string `json:"region"`
		State        string `json:"devicestate"`
		Secret       string `json:"secret"`
	}
	type response struct {
		Status string           `json:"status"`
		Reason string           `json:"reason"`
		Secret string           `json:"secret"`
		Device registeredDevice `json:"device"`
	}

	var resp response
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return "", apiContracts.ErrAPI_Service_CantReadResponse
	}

	// fmt.Printf("RESPONSE\n%s\n", string(raw))

	// Handle a variety of possible error conditions
	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_DEVICE_NOT_FOUND) {
			return "", apiContracts.ErrAPI_Helpers_ServiceUIDNotFound
		}
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_DUPLICATE_NAME) {
			return "", errorx.New(apiContracts.ErrAPI_Helpers_Generic, fmt.Sprintf(`a device in your account already has the name "%s"`, name))
		}
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_BAD_DEVICE_ADDRESS) {
			return "", apiContracts.ErrAPI_Helpers_ServiceUIDMissing
		}
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_API_TOKEN) {
			return "", apiContracts.ErrAPI_Helpers_NoToken
		}
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return "", errorx.New(apiContracts.ErrAPI_Helpers_Generic, resp.Reason)
		}
		return "", apiContracts.ErrAPI_Helpers_Unknown
	}

	return strings.Replace(resp.Secret, ":", "", -1), nil
}
