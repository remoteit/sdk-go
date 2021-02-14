package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

var cachedAuthResponse apiContracts.Authentication
var cachedAuthResponseCreateTime = time.Time{}
var cachedAuthResponseMutex sync.Mutex

const cachedAuthResponseExpireDuration = 1 * time.Hour

type Client interface {
	CanConnect(onlineCheckEndpoint string, onlineCheckEndpointReply string) bool

	LoginWithPassword(username string, password string) (apiContracts.Authentication, errorx.Error)

	LoginWithAuthHash(username string, authHash string) (apiContracts.Authentication, errorx.Error)
	LoginWithAuthHashIgnoreCache(username string, authHash string) (apiContracts.Authentication, errorx.Error)
	ExpireAuthHash()

	Post(endpointURL string, payload []byte) ([]byte, errorx.Error)
	Get(endpointURL string) ([]byte, errorx.Error)
}

func NewClient(apiURL string, apiKey string, apiTimeout time.Duration, userAgent string) Client {
	return &client{
		apiURL:     apiURL,
		apiKey:     apiKey,
		apiTimeout: apiTimeout,
		userAgent:  userAgent,
	}
}

type client struct {
	apiURL     string
	apiKey     string
	apiTimeout time.Duration
	userAgent  string
}

func (thisRef client) CanConnect(onlineCheckEndpoint string, onlineCheckEndpointReply string) bool {
	response, err := http.Get(onlineCheckEndpoint)
	if err != nil {
		return false
	}

	responseAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	response.Body.Close()

	responseAsString := strings.TrimSpace(string(responseAsBytes))

	return responseAsString == onlineCheckEndpointReply
}

func (thisRef client) LoginWithPassword(username string, password string) (apiContracts.Authentication, errorx.Error) {
	type requestT struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	requestPayload, err := json.Marshal(requestT{password, username})
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_CantPrepPasswordSignin
	}

	responsePayload, err := thisRef.Post("/user/login", requestPayload)
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_CantSendPasswordSignin
	}

	var response loginResponse
	err = json.Unmarshal(responsePayload, &response)
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_CantReadPasswordSignin
	}

	if response.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		errCodes := []string{
			apiContracts.API_ERROR_CODE_REASON_MFA_1,
			apiContracts.API_ERROR_CODE_REASON_MFA_2,
			apiContracts.API_ERROR_CODE_REASON_MFA_3,
		}

		for _, v := range errCodes {
			if v == response.Code {
				return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_MFA_ENABLED
			}
		}

		if strings.Contains(response.Reason, apiContracts.API_ERROR_CODE_REASON_USER_OR_PASSWORD_INVALID) {
			return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_PasswordInvalid
		}
		if strings.Contains(response.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_USER) {
			return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_NoSuchUser
		}
		if !(len(strings.TrimSpace(response.Reason)) <= 0) {
			return apiContracts.Authentication{}, errorx.New(apiContracts.ErrAPI_Auth_Generic, response.Reason)
		}
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_Unknown
	}

	if response.ServiceAuthHash == "" {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_NoAuthHash
	}

	return apiContracts.Authentication{
		AuthHash: response.ServiceAuthHash,
		User:     username,
		UserID:   response.GUID,
		Token:    response.Token,
	}, nil
}

func (thisRef client) LoginWithAuthHash(username string, authHash string) (apiContracts.Authentication, errorx.Error) {

	cachedAuthResponseMutex.Lock()
	defer cachedAuthResponseMutex.Unlock()

	if !cachedAuthResponseCreateTime.IsZero() &&
		time.Since(cachedAuthResponseCreateTime) < cachedAuthResponseExpireDuration {
		return cachedAuthResponse, nil
	}

	return thisRef.LoginWithAuthHashIgnoreCache(username, authHash)
}

func (thisRef client) LoginWithAuthHashIgnoreCache(username string, authHash string) (apiContracts.Authentication, errorx.Error) {
	type requestBody struct {
		AuthHash string `json:"authhash"`
		Username string `json:"username"`
	}

	creds := requestBody{authHash, username}
	body, err := json.Marshal(creds)
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_AuthHashCantPrepRequest
	}

	// Send the API request
	raw, err := thisRef.Post("/user/login/authhash", body)
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_AuthHashCantSendRequest
	}

	var response loginResponse
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_AuthHashCantReadResult
	}

	if response.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		// [0102] The username or password are invalid
		if strings.Contains(response.Reason, apiContracts.API_ERROR_CODE_REASON_USER_OR_PASSWORD_INVALID) {
			return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_AuthHashInvalid
		}
		if strings.Contains(response.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_USER) {
			return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_NoSuchUser
		}
		if !(len(strings.TrimSpace(response.Reason)) <= 0) {
			return apiContracts.Authentication{}, errorx.New(apiContracts.ErrAPI_Auth_Generic, response.Reason)
		}
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_Unknown
	}

	if response.Token == "" {
		return apiContracts.Authentication{}, apiContracts.ErrAPI_Auth_NoToken
	}

	cachedAuthResponseMutex.Lock()
	defer cachedAuthResponseMutex.Unlock()

	cachedAuthResponse = apiContracts.Authentication{
		AuthHash: response.ServiceAuthHash,
		User:     username,
		UserID:   response.GUID,
		Token:    response.Token,
	}
	cachedAuthResponseCreateTime = time.Now()
	return cachedAuthResponse, nil
}

func (thisRef client) ExpireAuthHash() {
	cachedAuthResponseCreateTime = time.Time{}
}

func (thisRef client) Post(endpointURL string, payload []byte) ([]byte, errorx.Error) {
	return thisRef.prepAndDoHTTPRequest("POST", endpointURL, payload)
}

func (thisRef client) Get(endpointURL string) ([]byte, errorx.Error) {
	return thisRef.prepAndDoHTTPRequest("GET", endpointURL, nil)
}

func (thisRef client) prepAndDoHTTPRequest(method string, endpointURL string, payload []byte) ([]byte, errorx.Error) {
	headers := map[string]string{
		"User-Agent": thisRef.userAgent,
		"apikey":     thisRef.apiKey,
	}

	cachedAuthResponseMutex.Lock()
	if cachedAuthResponse.Token != "" {
		headers["token"] = cachedAuthResponse.Token
	}
	cachedAuthResponseMutex.Unlock()

	_, data, err := doHTTPRequest(method, headers, thisRef.apiURL+endpointURL, payload, thisRef.apiTimeout)
	if err != nil {
		return nil, apiContracts.ErrAPI_Client_Error
	}

	return data, nil
}

type loginResponse struct {
	ServiceAuthHash string `json:"service_authhash"`
	GUID            string `json:"guid"`
	Status          string `json:"status"`
	Reason          string `json:"reason"`
	Code            string `json:"code"`
	Token           string `json:"token"`
}
