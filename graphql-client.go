package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

type GraphQLClient struct {
	endpointURL    string
	client         *http.Client
	productVersion string
	os             string
	osVersion      string
	skipQuery      bool
}

func NewGraphQLClient(
	gqlURL string, apiTimeout time.Duration,
	productVersion string, os string, osVersion string,
	skipQuery bool,
) *GraphQLClient {
	return &GraphQLClient{
		endpointURL: gqlURL,
		client: &http.Client{
			Timeout: apiTimeout,
		},
		productVersion: productVersion,
		os:             os,
		osVersion:      osVersion,
		skipQuery:      skipQuery,
	}
}

func (thisRef GraphQLClient) RunGraphQL(query string, authToken string) ([]byte, errorx.Error) {
	return thisRef.runGraphQLHelper(query, authToken)
}

func (thisRef GraphQLClient) runGraphQLHelper(query string, token string) ([]byte, errorx.Error) {

	if thisRef.skipQuery {
		return []byte{}, nil
	}

	type gqlReuqest struct {
		Query string `json:"query"`
	}

	request := gqlReuqest{Query: query}
	requestAsBytes, err := json.Marshal(request)
	if err != nil {
		return []byte{}, apiContracts.ErrAPI_GQL_CantPrepRequest
	}

	trimmedRequest := string(requestAsBytes)
	trimmedRequest = strings.ReplaceAll(trimmedRequest, "\t", "")
	trimmedRequest = strings.ReplaceAll(trimmedRequest, "\n", "")

	requestAsBytes = []byte(requestAsBytes)

	req, err := http.NewRequest("POST", thisRef.endpointURL, bytes.NewBuffer(requestAsBytes))
	if err != nil {
		return nil, apiContracts.ErrAPI_GQL_CantPrepRequest
	}

	req.Header.Set("User-Agent", fmt.Sprintf("cli-%s-%s-%s", thisRef.productVersion, thisRef.os, thisRef.osVersion))
	req.Header.Add("token", token)

	resp, err := thisRef.client.Do(req)
	if err != nil {
		return nil, apiContracts.ErrAPI_GQL_CantSendRequest
	}

	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	return raw, nil
}

var cachedApplicationTypes *[]apiContracts.ApplicationType
var cachedApplicationTypesCreateTime = time.Time{}
var cachedApplicationTypesMutex sync.Mutex

const cachedApplicationTypesExpireDuration = 10 * time.Hour

func (thisRef GraphQLClient) GetApplicationTypes(authToken string) (*[]apiContracts.ApplicationType, errorx.Error) {
	// 0. check for cached
	cachedApplicationTypesMutex.Lock()
	defer cachedApplicationTypesMutex.Unlock()

	if !cachedApplicationTypesCreateTime.IsZero() &&
		time.Since(cachedApplicationTypesCreateTime) < cachedApplicationTypesExpireDuration {
		return cachedApplicationTypes, nil
	}

	// 1. run
	raw, err := thisRef.RunGraphQL(`{
		applicationTypes {
			id
			name
			description
			port
			proxy
			protocol
		}
	}`, authToken)
	if err != nil {
		return &[]apiContracts.ApplicationType{}, err
	}

	// 2. read
	rawAsString := string(raw)
	if strings.TrimSpace(rawAsString) == "Unauthorized" {
		return &[]apiContracts.ApplicationType{}, apiContracts.ErrAPI_GQL_NotAuthorized
	}

	type gqlReply struct {
		Data struct {
			ApplicationTypes []apiContracts.ApplicationType `json:"applicationTypes"`
		} `json:"data"`
	}

	var response gqlReply
	if err := json.Unmarshal(raw, &response); err != nil {
		return &[]apiContracts.ApplicationType{}, apiContracts.ErrAPI_GQL_CantReadResponse
	}

	// 3. update cached
	cachedApplicationTypes = &response.Data.ApplicationTypes
	cachedApplicationTypesCreateTime = time.Now()

	return cachedApplicationTypes, nil
}

func (thisRef GraphQLClient) GetApplicationType(serviceID string, authToken string) (int, errorx.Error) {
	serviceID = strings.TrimSpace(serviceID)
	if len(serviceID) == 0 {
		return apiContracts.InvalidApplicationType, nil
	}

	// 1. run
	raw, err := thisRef.RunGraphQL(fmt.Sprintf(`{
		login {
			service(id: "%s") {
			  application
			}
		}
	}`, serviceID), authToken)
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

func (thisRef GraphQLClient) GetDeviceAndServiceNames(deviceID string, authToken string) (apiContracts.DefinedDevice, errorx.Error) {
	deviceID = strings.TrimSpace(deviceID)
	if len(deviceID) == 0 {
		return apiContracts.DefinedDevice{}, nil
	}

	// 1. run
	raw, err := thisRef.RunGraphQL(fmt.Sprintf(`{
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
	}`, deviceID), authToken)
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

func (thisRef GraphQLClient) GetServiceNamesByIDs(serviceIDs []string, authToken string) ([]apiContracts.DefinedService, errorx.Error) {
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

	raw, err := thisRef.RunGraphQL(fmt.Sprintf(`{
		login {
			service(id: [%s]) {
				id
				name
			}
		}
	}`, strings.Join(updatedServiceIDs, ", ")), authToken)
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
