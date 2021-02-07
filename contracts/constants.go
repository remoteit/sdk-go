package contracts

import "time"

const (
	CONNECTION_IP_LATCHING      = "255.255.255.255" // connection allowed to first connected, refuse others
	PROXY_CREATE_WAIT           = "true"
	PROXY_CREATE_ISOLATE        = "domain=app.remote.it"  // FIXME:  for now this is like this, but in portal it is taken from browser URL
	PROXY_CREATE_CONCURRENT     = true                    // INFO: pass this JSON flag to API to enable concurrent proxies, FIXME: will this be remvoed in the future ?
	ONLINE_CHECK_ENDPOINT       = "https://api.remot3.it" //
	ONLINE_CHECK_ENDPOINT_REPLY = "api.remot3.it"
	API_TIMEOUT                 = 40 * time.Second

	// These are here for two reasons
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
