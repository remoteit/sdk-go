package api

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

var cachedAuthResponse *apiContracts.AuthResponse
var cachedAuthResponseCreateTime = time.Time{}
var cachedAuthResponseMutex sync.Mutex

const cachedAuthResponseExpireDuration = 1 * time.Hour

func ExpireAuthHash() {
	cachedAuthResponseCreateTime = time.Time{}
}

/*
PasswordSignin handles password based logins.
*/
func (c *Client) PasswordSignin(username string, password string) (*apiContracts.Auth, errorx.Error) {
	// Construct the request data to send to the API
	type requestBody struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}
	creds := requestBody{password, username}
	body, err := json.Marshal(creds)
	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_CantPrepPasswordSignin
	}

	// Send the API request
	raw, err := c.Post("/user/login", body)
	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_CantSendPasswordSignin
	}

	// Parse the JSON response into a usable struct.
	errCodes := []string{
		apiContracts.API_ERROR_CODE_REASON_MFA_1,
		apiContracts.API_ERROR_CODE_REASON_MFA_2,
		apiContracts.API_ERROR_CODE_REASON_MFA_3,
	}

	var resp apiContracts.AuthResponse
	err = json.Unmarshal(raw, &resp)

	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_CantReadPasswordSignin
	}

	// Handle incoming API error messages in a structured way.
	// For common errors, we create a specific error type so it
	// can be checked by consuming code.
	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		for _, v := range errCodes {
			if v == resp.Code {
				return nil, apiContracts.ErrAPI_Auth_MFA_ENABLED
			}
		}

		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_USER_OR_PASSWORD_INVALID) {
			return nil, apiContracts.ErrAPI_Auth_PasswordInvalid
		}
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_USER) {
			return nil, apiContracts.ErrAPI_Auth_NoSuchUser
		}
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return nil, errorx.New(apiContracts.ErrAPI_Auth_Generic, resp.Reason)
		}
		return nil, apiContracts.ErrAPI_Auth_Unknown
	}

	// Although unlikely the API won't return an auth hash, we should
	// handle that case anyways.
	if resp.AuthHash == "" {
		return nil, apiContracts.ErrAPI_Auth_NoAuthHash
	}

	// Everything was successful so return the credentials
	return &apiContracts.Auth{
		AuthHash: resp.AuthHash,
		Username: username,
		Guid:     resp.Guid,
	}, nil
}

func (c *Client) AuthHashSigninNoCache(username string, authHash string) (*apiContracts.AuthResponse, errorx.Error) {
	type requestBody struct {
		AuthHash string `json:"authhash"`
		Username string `json:"username"`
	}

	creds := requestBody{authHash, username}
	body, err := json.Marshal(creds)
	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_AuthHashCantPrepRequest
	}

	// Send the API request
	raw, err := c.Post("/user/login/authhash", body)
	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_AuthHashCantSendRequest
	}

	var resp apiContracts.AuthResponse
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, apiContracts.ErrAPI_Auth_AuthHashCantReadResult
	}

	if resp.Status == apiContracts.API_ERROR_CODE_STATUS_FALSE {
		// [0102] The username or password are invalid
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_USER_OR_PASSWORD_INVALID) {
			return nil, apiContracts.ErrAPI_Auth_AuthHashInvalid
		}
		if strings.Contains(resp.Reason, apiContracts.API_ERROR_CODE_REASON_MISSING_USER) {
			return nil, apiContracts.ErrAPI_Auth_NoSuchUser
		}
		if !(len(strings.TrimSpace(resp.Reason)) <= 0) {
			return nil, errorx.New(apiContracts.ErrAPI_Auth_Generic, resp.Reason)
		}
		return nil, apiContracts.ErrAPI_Auth_Unknown
	}

	if resp.Token == "" {
		return nil, apiContracts.ErrAPI_Auth_NoToken
	}

	// Everything was successful so return the credentials
	cachedAuthResponse = &resp
	cachedAuthResponseCreateTime = time.Now()

	c.SetAuthToken(resp.Token)
	return &resp, nil
}

func (c *Client) AuthHashSignin(username string, authHash string) (*apiContracts.AuthResponse, errorx.Error) {

	cachedAuthResponseMutex.Lock()
	defer cachedAuthResponseMutex.Unlock()

	if !cachedAuthResponseCreateTime.IsZero() &&
		time.Since(cachedAuthResponseCreateTime) < cachedAuthResponseExpireDuration {

		c.SetAuthToken(cachedAuthResponse.Token)
		return cachedAuthResponse, nil
	}

	return c.AuthHashSigninNoCache(username, authHash)
}
