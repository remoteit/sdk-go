package contracts

import "time"

const (
	DEFAULT_API_URL         = "https://api.remot3.it/apv/v27"
	DEFAULT_API_GRAPHQL_URL = "https://api.remote.it/graphql/v1"
	DEFAULT_API_RESTORE_URL = "https://install.remote.it/v1/restore"
	DEFAULT_API_TIMEOUT     = 40 * time.Second
	DEFAULT_API_USER_AGENT  = "remot3.it-sdk-go"

	DEFAULT_PROXY_CREATE_IP_LATCHING = "255.255.255.255"
	DEFAULT_PROXY_CREATE_WAIT        = "true"
	DEFAULT_PROXY_CREATE_ISOLATE     = "domain=app.remote.it"
	DEFAULT_PROXY_CREATE_CONCURRENT  = true

	DEFAULT_ONLINE_CHECK_ENDPOINT       = "https://api.remot3.it"
	DEFAULT_ONLINE_CHECK_ENDPOINT_REPLY = "api.remot3.it"

	// FIXME:
	// 1 - error code is a string with no standard
	// 2 - is part of the `resp.Reason` instead of `resp.Error`,
	//     `resp.Reason` must be a human readable form of the `resp.Error`
	API_ERROR_CODE_REASON_DEVICE_NOT_FOUND          = "[0806]"
	API_ERROR_CODE_REASON_DUPLICATE_NAME            = "[0807]"
	API_ERROR_CODE_REASON_SERVICE_NOT_FOUND_FOR_UID = "[0861]"
	API_ERROR_CODE_REASON_BAD_DEVICE_ADDRESS        = "bad device address"
	API_ERROR_CODE_REASON_MISSING_API_TOKEN         = "missing api token"
	API_ERROR_CODE_REASON_USER_OR_PASSWORD_INVALID  = "username or password are invalid"
	API_ERROR_CODE_REASON_MISSING_USER              = "missing user"
	API_ERROR_CODE_REASON_NO_MATCHING_BULK_PROJECT  = "no matching bulk project"
	API_ERROR_CODE_REASON_MFA_1                     = "SMS_MFA"
	API_ERROR_CODE_REASON_MFA_2                     = "SOFTWARE_TOKEN_MFA"
	API_ERROR_CODE_REASON_MFA_3                     = "MFA_SETUP"
	API_ERROR_CODE_STATUS_FALSE                     = "false"
	API_ERROR_CODE_STATUS_TRUE                      = "true"
	API_ERROR_CODE_STATUS_PENDING                   = "pending"
	API_ERROR_CODE_STATUS_RESET                     = "reset"
)
