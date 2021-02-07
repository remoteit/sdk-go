package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

// X ->	?????????????																		-> SendDeviceInfo
// X -> ProjectListURL				="/bulk/registration/device/friendly/configuration"		-> ListServiceIDs
// X -> ProvisionConfig				="/bulk/registration/configuration"						-> GetServiceConfig
// X -> BulkRegisterURL				="/bulk/registration/register"							-> RegisterService
// - -> ComponentVersionURL			="/device/component/version"
// ? -> ProjectEnablementGet		="/device/enablement"
// - -> ProvisionGet				="/project/provisioning"
// - -> ProvisionDownloadDirect		="/project/provisioning/download"

func SendDeviceInfo(registrationKey string, hardwareID string, cpuID string, macAddress string, version string,
	apiURL string, apiKey string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string,
	platformOSName string) errorx.Error {

	apiClient := NewClient(apiURL, apiKey, apiTimeout, productVersion, os, osVersion)

	var url = "/bulk/registration/device/information/"

	type deviceInfoRequest struct {
		BulkIdentificationCode string `json:"BulkIdentificationCode"`
		HardwareID             string `json:"HardwareId"`
		MACAddress             string `json:"MACAddress"`
		CPUId                  string `json:"CPUId"`
		OSLabel                string `json:"OSLabel"`
		R3Package              string `json:"R3Package"`
		TCPServiceList         string `json:"TCPServiceList"`
	}

	var req = deviceInfoRequest{
		BulkIdentificationCode: registrationKey,
		HardwareID:             hardwareID,
		MACAddress:             macAddress,
		CPUId:                  cpuID,
		OSLabel:                fmt.Sprintf("%v", platformOSName),
		R3Package:              version,
		TCPServiceList:         "",
		// "DeviceSecret": "",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return apiContracts.ErrAutoreg_CantPrepRequest
	}

	raw, err := apiClient.Post(url, body)
	if err != nil {
		return apiContracts.ErrAutoreg_CantSendRequest
	}

	type deviceInfoResponse struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	var resp deviceInfoResponse
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return apiContracts.ErrAutoreg_CantReadResponse
	}

	if resp.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_NO_MATCHING_BULK_PROJECT) {
			return apiContracts.ErrAutoreg_NoMatchingRegInfo
		}

		return errorx.New(apiContracts.ErrAutoReg_Generic, resp.Reason)
	}

	return nil
}

// GetProductTemplate -
func GetProductTemplate(registrationKey string, deviceUniquID string,
	apiURL string, apiKey string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string) ([]string, bool, errorx.Error) {

	apiClient := NewClient(apiURL, apiKey, apiTimeout, productVersion, os, osVersion)

	var url = fmt.Sprintf("/bulk/registration/device/friendly/configuration/%s/%s/", registrationKey, deviceUniquID)

	raw, errx := apiClient.Get(url)
	if errx != nil {
		return nil, false, errx
	}

	type projectsResponse struct {
		Projects   string `json:"projects"`
		Status     string `json:"status"`
		Reason     string `json:"reason"`
		Type       string `json:"type"`
		Registered string `json:"registered"`
		HardwareID string `json:"HardwareId"`
	}

	var resp projectsResponse
	err := json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, false, apiContracts.ErrAutoreg_CantPrepRequest
	}

	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_RESET {
		return nil, true, nil
	}

	return strings.Split(resp.Projects, ","), false, nil
}

// GetServiceConfigFromTemplateID -
func GetServiceConfigFromTemplateID(serviceID string, hardwareID string,
	apiURL string, apiKey string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string) (*ServiceConfigResponse, errorx.Error) {

	apiClient := NewClient(apiURL, apiKey, apiTimeout, productVersion, os, osVersion)

	var url = fmt.Sprintf("/bulk/registration/configuration/%s/%s/", serviceID, hardwareID)

	raw, errx := apiClient.Get(url)
	if errx != nil {
		return nil, errx
	}

	type serviceConfigResponse struct {
		ContentIP   string `json:"content_ip"`
		ContentPort string `json:"content_port"`
		ContentType string `json:"content_type"`
		Enabled     string `json:"enabled"`
		ProjectID   string `json:"project_id"`
		Status      string `json:"status"`
		Timestamp   string `json:"timestamp"`
	}

	var resp serviceConfigResponse
	err := json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, apiContracts.ErrAutoreg_CantPrepRequest
	}

	port, err := strconv.Atoi(resp.ContentPort)
	if err != nil {
		port = 65535
	}

	var serviceType int
	t, err := strconv.Atoi(resp.ContentType)
	if err != nil {
		serviceType = apiContracts.BulkServiceID
	} else {
		serviceType = t
	}

	if err != nil {
		return nil, apiContracts.ErrAutoreg_CantConvertPort
	}

	return &ServiceConfigResponse{
		Hostname: resp.ContentIP,
		Port:     port,
		Type:     serviceType,
		Disabled: resp.Enabled != "1",
	}, nil
}

// ServiceConfigResponse -
type ServiceConfigResponse struct {
	Hostname string
	Port     int
	Type     int
	Disabled bool
}

// RegisterService -
func RegisterService(serviceID string, uniqueDeviceID string, registrationKey string,
	apiURL string, apiKey string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string) (*ServiceCredentials, bool, errorx.Error) {

	apiClient := NewClient(apiURL, apiKey, apiTimeout, productVersion, os, osVersion)

	var url = "/bulk/registration/register"

	type deviceRegistrationRequest struct {
		RegistrationKey string `json:"registration_key"`
		HardwareID      string `json:"hardware_id"`
		ProjectID       string `json:"project_id"`
	}

	var req = deviceRegistrationRequest{
		RegistrationKey: registrationKey,
		HardwareID:      uniqueDeviceID,
		ProjectID:       serviceID,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, false, apiContracts.ErrAutoreg_CantPrepRequest
	}

	raw, errx := apiClient.Post(url, body)
	if errx != nil {
		return nil, false, errx
	}

	type deviceRegistrationResponse struct {
		Status       string `json:"status"`
		Reason       string `json:"reason"`
		Registration string `json:"registration"`
		UID          string `json:"uid"`
		Secret       string `json:"secret"`
	}

	var resp deviceRegistrationResponse
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, false, apiContracts.ErrAutoreg_CantReadResponse
	}

	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		return nil, false, errorx.New(apiContracts.ErrAPI_AutoReg_Generic, resp.Reason)
	}
	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_PENDING {
		return nil, true, nil
	}

	return &ServiceCredentials{
		UID:    resp.UID,
		Secret: resp.Secret,
	}, false, nil
}

// ServiceCredentials -
type ServiceCredentials struct {
	UID    string `json:"uid"`
	Secret string `json:"secret"`
}
