package api

import (
	"encoding/json"
	"fmt"
	"strings"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type Service interface {
	Create(uid string, serviceType string) errorx.Error
	Remove(uid string) errorx.Error
	GenerateUID(projectKey string, projectSecret string) (uid string, err errorx.Error)
	Register(name string, uid string, hardwareID string, serviceType string, serviceTypeAsInt int) (secret string, err errorx.Error)

	CreateFullService(info apiContracts.ServiceRegistrationInfo, projectKey string, projectSecret string) (*apiContracts.Service, errorx.Error)
}

func NewService(apiClient Client) Service {
	return &service{
		apiClient: apiClient,
	}
}

type service struct {
	apiClient Client
}

func (thisRef service) Create(uid string, serviceType string) errorx.Error {
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
	raw, err := thisRef.apiClient.Post("/device/create", body)
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

func (thisRef service) Remove(uid string) errorx.Error {
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
	raw, errx := thisRef.apiClient.Post("/device/delete", body)
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

func (thisRef service) GenerateUID(projectKey string, projectSecret string) (string, errorx.Error) {
	// Send the API request
	raw, errx := thisRef.apiClient.Get(fmt.Sprintf("/device/address/%s/%s", projectKey, projectSecret))
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

func (thisRef service) Register(name string, uid string, hardwareID string, serviceType string, serviceTypeAsInt int) (string, errorx.Error) {
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
	raw, errx := thisRef.apiClient.Post("/device/register", body)
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

func (thisRef service) CreateFullService(info apiContracts.ServiceRegistrationInfo, projectKey string, projectSecret string) (*apiContracts.Service, errorx.Error) {
	empty := new(apiContracts.Service)

	uid, err := thisRef.GenerateUID(projectKey, projectSecret)
	if err != nil {
		return empty, err
	}

	err = thisRef.Create(uid, info.ServiceType)
	if err != nil {
		return empty, err
	}

	hardwareID := uid
	if info.HardwareID != "" {
		hardwareID = info.HardwareID
	}
	secret, err := thisRef.Register(info.Name, uid, hardwareID, info.ServiceType, info.ServiceTypeAsInt)
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
