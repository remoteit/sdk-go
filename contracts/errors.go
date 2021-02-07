package contracts

import errorx "github.com/remoteit/systemkit-errorx"

var (
	MFA_IS_ENABLED = "Sorry, your account has Two-Factor authentication and/or Google Auth configured. This version of the CLI does not support these features. Please either upgrade to a newer version of the CLI that supports these features or disable them in your account to use this version of the CLI. Please contact support@remote.it if you have any questions"

	ErrAutoReg_Generic           = 300
	ErrAutoreg_BICInvalid        = errorx.New(301, "AutoReg - RegistrationKey (Bulk Identification Code) is invalid")
	ErrAutoreg_MaxAttempts       = errorx.New(302, "AutoReg - AutoRegisterIfNeeded: max nr of attemps was reached and they all failed")
	ErrAutoreg_BICEmpty          = errorx.New(303, "AutoReg - Registration key (Bulk Identification Code) is empty")
	ErrAutoreg_CantPrepRequest   = errorx.New(304, "AutoReg - Can't prep autoreg details")
	ErrAutoreg_CantSendRequest   = errorx.New(305, "AutoReg - Can't send autoreg details")
	ErrAutoreg_CantReadResponse  = errorx.New(306, "AutoReg - Can't read autoreg response")
	ErrAutoreg_NoMatchingRegInfo = errorx.New(307, "AutoReg - No registration exists that matches the key")
	ErrAutoreg_CantConvertPort   = errorx.New(308, "AutoReg - Can't convert port to integer")
	ErrAutoreg_CantInstallAgent  = errorx.New(309, "AutoReg - Can't install agent")
	ErrAutoreg_CantStartAgent    = errorx.New(310, "AutoReg - Can't start agent")

	ErrAPI_Auth_Generic                 = 500
	ErrAPI_Auth_Unknown                 = errorx.New(501, "Auth - Unknown error occurred")
	ErrAPI_Auth_MFA_ENABLED             = errorx.New(502, MFA_IS_ENABLED)
	ErrAPI_Auth_NoAuthHash              = errorx.New(503, "Auth - No auth hash returned")
	ErrAPI_Auth_NoToken                 = errorx.New(504, "Auth - No token returned")
	ErrAPI_Auth_CantPrepPasswordSignin  = errorx.New(505, "Auth - Can't prpe password signin")
	ErrAPI_Auth_CantSendPasswordSignin  = errorx.New(506, "Auth - Can't send password signin")
	ErrAPI_Auth_CantReadPasswordSignin  = errorx.New(507, "Auth - Can't read password signin reply")
	ErrAPI_Auth_PasswordInvalid         = errorx.New(508, "Auth - Password is invalid")
	ErrAPI_Auth_NoSuchUser              = errorx.New(509, "Auth - No such user")
	ErrAPI_Auth_AuthHashCantPrepRequest = errorx.New(510, "Auth - AuthHash can't prep request")
	ErrAPI_Auth_AuthHashCantSendRequest = errorx.New(511, "Auth - AuthHash can't send request")
	ErrAPI_Auth_AuthHashCantReadResult  = errorx.New(512, "Auth - AuthHash can't read result")
	ErrAPI_Auth_AuthHashInvalid         = errorx.New(513, "Auth - AuthHash invalid")

	ErrAPI_AutoReg_Generic = 600

	ErrAPI_DeviceList_Generic          = 700
	ErrAPI_DeviceList_CantSendRequest  = errorx.New(701, "Device list - Can't send request")
	ErrAPI_DeviceList_CantReadResponse = errorx.New(702, "Device list - Can't read response")

	ErrAPI_Device_Generic          = 800
	ErrAPI_Device_Unknown          = errorx.New(801, "Device - Unknown error occurred")
	ErrAPI_Device_NoServiceFound   = errorx.New(802, "Device - No service found matching that UID")
	ErrAPI_Device_CantPrepRequest  = errorx.New(803, "Device - Can't prepare request")
	ErrAPI_Device_CantSendRequest  = errorx.New(804, "Device - Can't send request")
	ErrAPI_Device_CantReadResponse = errorx.New(805, "Device - Can't read response")

	ErrAPI_ProxyCreate_Generic          = 900
	ErrAPI_ProxyCreate_Unknown          = errorx.New(901, "Create Proxy - Unknown error occurred")
	ErrAPI_ProxyCreate_NoServiceFound   = errorx.New(902, "Create Proxy - No service found matching that UID")
	ErrAPI_ProxyCreate_CantPrepRequest  = errorx.New(903, "Create Proxy - Can't prep request")
	ErrAPI_ProxyCreate_CantSendRequest  = errorx.New(904, "Create Proxy - Can't send request")
	ErrAPI_ProxyCreate_CantReadResponse = errorx.New(905, "Create Proxy - Can't read response")

	ErrAPI_ProxyDelete_Generic          = 1000
	ErrAPI_ProxyDelete_Unknown          = errorx.New(1001, "Delete Proxy - Unknown error occurred")
	ErrAPI_ProxyDelete_CantPrepRequest  = errorx.New(1002, "Delete Proxy - Can't prep request")
	ErrAPI_ProxyDelete_CantSendRequest  = errorx.New(1003, "Delete Proxy - Can't send request")
	ErrAPI_ProxyDelete_CantReadResponse = errorx.New(1004, "Delete Proxy - Can't read response")

	ErrAPI_RestoreClient_Generic           = 2000
	ErrAPI_RestoreClient_Unknown           = errorx.New(2001, "Restore Client - Unknown error occurred")
	ErrAPI_RestoreClient_DeviceActive      = errorx.New(2002, "Restore Client - The device state is active")
	ErrAPI_RestoreClient_TokenNotSpecified = errorx.New(2003, "Restore Client - Token not specified or invalid")
	ErrAPI_RestoreClient_DeviceNotExists   = errorx.New(2004, "Restore Client - Device does not exist or is not owned by the user")
	ErrAPI_RestoreClient_CantPrepRequest   = errorx.New(2005, "Restore Client - Can't prep request")
	ErrAPI_RestoreClient_CantSendRequest   = errorx.New(2006, "Restore Client - Can't send request")
	ErrAPI_RestoreClient_CantReadResponse  = errorx.New(2007, "Restore Client - Can't read response")

	ErrAPI_Service_Generic          = 3000
	ErrAPI_Service_Unknown          = errorx.New(3001, "Service - Unknown error occurred")
	ErrAPI_Service_NoServiceFound   = errorx.New(3002, "Service - No service found matching that UID")
	ErrAPI_Service_CantPrepRequest  = errorx.New(3003, "Service - Can't prep request")
	ErrAPI_Service_CantSendRequest  = errorx.New(3004, "Service - Can't send request")
	ErrAPI_Service_CantReadResponse = errorx.New(3005, "Service - Can't read response")

	ErrAPI_Helpers_Generic             = 4000
	ErrAPI_Helpers_Unknown             = errorx.New(4001, "API Helpers - Unknown error occurred creating service")
	ErrAPI_Helpers_CantCreateServiceID = errorx.New(4002, "API Helpers - Unknown error occurred creating a service UID")
	ErrAPI_Helpers_NoUIDReturned       = errorx.New(4003, "API Helpers - No UID was returned by the API")
	ErrAPI_Helpers_NoToken             = errorx.New(4004, "API Helpers - Missing authentication token")
	ErrAPI_Helpers_ServiceUIDNotFound  = errorx.New(4005, "API Helpers - Service matching UID not found")
	ErrAPI_Helpers_ServiceUIDMissing   = errorx.New(4006, "API Helpers - Service UID invalid or missing")
	ErrAPI_Helpers_DeviceExists        = errorx.New(4007, "API Helpers - Service UID invalid or missing")

	ErrAPI_Client_Generic           = 5000
	ErrAPI_Client_CantCreateRequest = errorx.New(5001, "API Client - Error creating request")
	ErrAPI_Client_CantSend          = errorx.New(5002, "API Client - Error sending request")
	ErrAPI_Client_CantRead          = errorx.New(5003, "API Client - Error reading request")

	ErrAPI_GQL_Generic          = 6000
	ErrAPI_GQL_CantPrepRequest  = errorx.New(6001, "GQL Client - Error creating request")
	ErrAPI_GQL_CantSendRequest  = errorx.New(6002, "GQL Client - Error sending request")
	ErrAPI_GQL_CantReadResponse = errorx.New(6003, "GQL Client - Error read response")
	ErrAPI_GQL_NotAuthorized    = errorx.New(6004, "GQL Client - Not authorized")
)
