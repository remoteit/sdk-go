package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

var cachedApplicationTypes []apiContracts.ApplicationType
var cachedApplicationTypesCreateTime = time.Time{}
var cachedApplicationTypesMutex sync.Mutex

const cachedApplicationTypesExpireDuration = 10 * time.Hour

type GraphQLClient interface {
	GetApplicationTypes() ([]apiContracts.ApplicationType, errorx.Error)
	GetApplicationType(serviceID string) (int, errorx.Error)
	GetDeviceAndServiceNames(deviceID string) (apiContracts.DefinedDevice, errorx.Error)
	GetServiceNamesByIDs(serviceIDs []string) ([]apiContracts.DefinedService, errorx.Error)
}

func NewGraphQLClient(apiURL string, apiToken string, apiTimeout time.Duration, userAgent string) GraphQLClient {
	return &graphQLClient{
		apiURL:     apiURL,
		apiToken:   apiToken,
		apiTimeout: apiTimeout,
		userAgent:  userAgent,
	}
}

type graphQLClient struct {
	apiURL     string
	apiToken   string
	apiTimeout time.Duration
	userAgent  string
}

func (thisRef graphQLClient) GetApplicationTypes() ([]apiContracts.ApplicationType, errorx.Error) {
	// 0. check for cached
	cachedApplicationTypesMutex.Lock()
	defer cachedApplicationTypesMutex.Unlock()

	if !cachedApplicationTypesCreateTime.IsZero() &&
		time.Since(cachedApplicationTypesCreateTime) < cachedApplicationTypesExpireDuration {
		return cachedApplicationTypes, nil
	}

	// 1. run
	raw, err := thisRef.prepAndDoHTTPRequest(`{
		applicationTypes {
			id
			name
			description
			port
			proxy
			protocol
		}
	}`)
	if err != nil {
		return []apiContracts.ApplicationType{}, err
	}

	// 2. read
	rawAsString := string(raw)
	if strings.TrimSpace(rawAsString) == "Unauthorized" {
		return []apiContracts.ApplicationType{}, apiContracts.ErrAPI_GQL_NotAuthorized
	}

	type gqlReply struct {
		Data struct {
			ApplicationTypes []apiContracts.ApplicationType `json:"applicationTypes"`
		} `json:"data"`
	}

	var response gqlReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return []apiContracts.ApplicationType{}, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	// 3. update cached
	cachedApplicationTypes = response.Data.ApplicationTypes
	cachedApplicationTypesCreateTime = time.Now()

	return cachedApplicationTypes, nil
}

func (thisRef graphQLClient) GetApplicationType(serviceID string) (int, errorx.Error) {
	serviceID = strings.TrimSpace(serviceID)
	if len(serviceID) == 0 {
		return apiContracts.InvalidApplicationType, nil
	}

	// 1. run
	raw, err := thisRef.prepAndDoHTTPRequest(fmt.Sprintf(`{
		login {
			service(id: "%s") {
			  application
			}
		}
	}`, serviceID))
	if err != nil {
		return apiContracts.InvalidApplicationType, err
	}

	// 2. read
	type gqlReply struct {
		Data struct {
			Login struct {
				Service []struct {
					Application int `json:"application"`
				} `json:"service"`
			} `json:"login"`
		} `json:"data"`
	}

	var response gqlReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return apiContracts.InvalidApplicationType, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	// 3. return
	if len(response.Data.Login.Service) > 0 {
		return response.Data.Login.Service[0].Application, nil
	}

	return apiContracts.InvalidApplicationType, nil
}

func (thisRef graphQLClient) GetDeviceAndServiceNames(deviceID string) (apiContracts.DefinedDevice, errorx.Error) {
	deviceID = strings.TrimSpace(deviceID)
	if len(deviceID) == 0 {
		return apiContracts.DefinedDevice{}, nil
	}

	// 1. run
	raw, err := thisRef.prepAndDoHTTPRequest(fmt.Sprintf(`{
		login {
			device(id: "%s") {
				id
				name
				services {
					id
					name
				}
			}
		}
	}`, deviceID))
	if err != nil {
		return apiContracts.DefinedDevice{}, err
	}

	// 2. read
	type gqlReply struct {
		Data struct {
			Login struct {
				Device []apiContracts.DefinedDevice `json:"device"`
			} `json:"login"`
		} `json:"data"`
	}

	var response gqlReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return apiContracts.DefinedDevice{}, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	// 3. return
	if len(response.Data.Login.Device) > 0 {
		definedDevice := apiContracts.DefinedDevice{
			ID:       response.Data.Login.Device[0].ID,
			Name:     response.Data.Login.Device[0].Name,
			Services: []apiContracts.DefinedService{},
		}

		if len(response.Data.Login.Device[0].Services) > 0 {
			for _, service := range response.Data.Login.Device[0].Services {
				definedDevice.Services = append(definedDevice.Services, service)
			}
		}

		return definedDevice, nil
	}

	return apiContracts.DefinedDevice{}, nil
}

func (thisRef graphQLClient) GetServiceNamesByIDs(serviceIDs []string) ([]apiContracts.DefinedService, errorx.Error) {
	// 1. run
	updatedServiceIDs := []string{}
	for _, serviceID := range serviceIDs {
		trimmedServiceID := strings.TrimSpace(serviceID)
		if trimmedServiceID == "" {
			continue
		}

		updatedServiceIDs = append(updatedServiceIDs, "\""+trimmedServiceID+"\"")
	}

	if len(updatedServiceIDs) == 0 {
		return []apiContracts.DefinedService{}, nil
	}

	raw, err := thisRef.prepAndDoHTTPRequest(fmt.Sprintf(`{
		login {
			service(id: [%s]) {
				id
				name
			}
		}
	}`, strings.Join(updatedServiceIDs, ", ")))
	if err != nil {
		return []apiContracts.DefinedService{}, err
	}

	// 2. read
	type gqlReply struct {
		Data struct {
			Login struct {
				Services []apiContracts.DefinedService `json:"service"`
			} `json:"login"`
		} `json:"data"`
	}

	var response gqlReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return []apiContracts.DefinedService{}, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	// 3. return
	definedServices := []apiContracts.DefinedService{}
	for _, definedService := range response.Data.Login.Services {
		for _, serviceID := range serviceIDs {
			if definedService.ID == serviceID {
				definedServices = append(definedServices, definedService)
				break
			}
		}
	}

	return definedServices, nil
}

func (thisRef graphQLClient) prepAndDoHTTPRequest(query string) ([]byte, errorx.Error) {
	type gqlReuqest struct {
		Query string `json:"query"`
	}

	request := gqlReuqest{Query: query}
	payload, err := json.Marshal(request)
	if err != nil {
		return []byte{}, apiContracts.ErrAPI_GQL_CantPrepRequest
	}

	headers := map[string]string{
		"User-Agent": thisRef.userAgent,
		"token":      thisRef.apiToken,
	}

	_, data, err := doHTTPRequest(http.MethodPost, headers, thisRef.apiURL, payload, thisRef.apiTimeout)
	if err != nil {
		return nil, apiContracts.ErrAPI_GQL_Error
	}

	return data, nil
}
