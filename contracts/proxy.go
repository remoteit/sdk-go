package contracts

type CreateProxyResponse struct {
	Status     string              `json:"status"`     // "true"
	Reason     string              `json:"reason"`     // "..."
	Connection ProxyConnectionInfo `json:"connection"` //

	//
	// Ignore for now
	//
	// Wait         bool   `json:"wait,omitempty"`         // tells the API to not return until the connection is made or times out
	// ConnectionID string `json:"connectionid,omitempty"` // "118EB6E7-1BA5-0AA0-B80E-7D15B9059260"
}

//
// SAMPLE
//
// {
//     "status": "true",
//     "connection": {
//         "deviceaddress": "80:00:00:00:01:00:40:C5",
//         "status": "running",
//         "requested": "1\/15\/2021T6:20 PM",
//         "proxy": "http:\/\/proxy40.rt3.io:34168",
//         "proxyserver": "proxy40.rt3.io",
//         "proxyport": "34168",
//         "expirationsec": "28800",
//         "connectionid": "48844999-bbf4-48cf-9a3a-353b7385737a",
//         "proxyServerPort": 34168,
//         "proxyExpirationSec": 28800,
//         "sessionID": "E2A92E32911A247214B962839DD03CE1A0CE2C7D",
//         "initiator": "nicolae@remote.it",
//         "targetUID": "80:00:00:00:01:00:40:C5",
//         "clientID": "WeavedRESTAPI6yThKw9",
//         "filteredIP": "latching",
//         "proxyURL": "proxy40.rt3.io:34168",
//         "reverseProxy": false,
//         "p2pConnected": true,
//         "serviceConnected": true,
//         "peerReqEP": "18.184.71.109:62292",
//         "peerEP": "18.184.71.109:62292",
//         "latchedIP": "0.0.0.0",
//         "initiatorUID": "f0:0f:85:c8:7b:3f:a9:e9",
//         "lifeLeft": 86400,
//         "idleLeft": 900,
//         "requestedAt": "2021-01-15T23:20:00+00:00"
//     },
//     "wait": true,
//     "connectionid": "48844999-bbf4-48cf-9a3a-353b7385737a"
// }
//
type ProxyConnectionInfo struct {
	ConnectionID string `json:"connectionid,omitempty"` // "118EB6E7-1BA5-0AA0-B80E-7D15B9059260"

	SessionID    string `json:"sessionID,omitempty"`    // "E2A92E32911A247214B962839DD03CE1A0CE2C7D"
	InitiatorUID string `json:"initiatorUID,omitempty"` // "f0:0f:85:c8:7b:3f:a9:e9" - proxy UID
	TargetUID    string `json:"targetUID,omitempty"`    // "80:00:00:00:01:00:40:C5" - target UID

	ProxyServer string `json:"proxyserver,omitempty"` // "proxy40.rt3.io"
	ProxyPort   string `json:"proxyport,omitempty"`   // "34168"

	Proxy           string `json:"proxy,omitempty"`           // "http://proxy40.rt3.io:34168"
	ProxyServerPort int    `json:"proxyServerPort,omitempty"` // 34168
	ProxyURL        string `json:"proxyURL,omitempty"`        // "proxy40.rt3.io:34168"
	ReverseProxy    bool   `json:"reverseProxy,omitempty"`    // false

	//
	// Ignore for now
	//
	// Initiator          string `json:"initiator,omitempty"`          // "nicolae@remote.it"
	// FilteredIP         string `json:"filteredIP,omitempty"`         // "latching"
	// P2PConnected       bool   `json:"p2pConnected,omitempty"`       // true
	// ServiceConnected   bool   `json:"serviceConnected,omitempty"`   // true
	// DeviceAddress      string `json:"deviceaddress,omitempty"`      // "80:00:00:00:01:00:40:C5" - target UID
	// Status             string `json:"status,omitempty"`             // "running"
	// ProxyExpirationSec int    `json:"proxyExpirationSec,omitempty"` // 28800 - ???????
	// ClientID           string `json:"clientID,omitempty"`           // "WeavedRESTAPI6yThKw9" - app-id string
	// PeerReqEP          bool   `json:"peerReqEP,omitempty"`          // "18.184.71.109:62292" - ??????
	// PeerEP             bool   `json:"peerEP,omitempty"`             // "18.184.71.109:62292" - ??????
	// LatchedIP          bool   `json:"latchedIP,omitempty"`          // "0.0.0.0" - ??????
	// LifeLeft           bool   `json:"lifeLeft,omitempty"`           // 86400 - ??????
	// IdleLeft           bool   `json:"idleLeft,omitempty"`           // 900 - ??????
	// Requested          string `json:"requested,omitempty"`          // "2/12/2020T10:18 AM"
	// RequestedAt        string `json:"requestedAt,omitempty"`        // "2021-01-15T23:20:00+00:00"
}

type DeleteProxyResponse struct {
	Status string `json:"status"` // "true"
}

type DeleteProxyRequest struct {
	DeviceAddress string `json:"deviceaddress,omitempty"` // "80:00:00:00:01:00:40:C4"
	ConnectionID  string `json:"connectionid,omitempty"`  // "118EB6E7-1BA5-0AA0-B80E-7D15B9059260"
}

//
// SAMPLE
// curl 'https://api.remote.it/apv/v27/device/connect/'
// -H 'developerKey: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
// -H 'token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
// -H 'Accept-Language: en-US,en;q=0.9'
// --data-binary '{"deviceaddress":"80:00:00:00:01:00:40:C4","devicetype":28,"hostip":"255.255.255.255","wait":"true","isolate":"domain=app.remote.it"}'
//
type CreateProxyRequest struct {
	DeviceAddress string `json:"deviceaddress,omitempty"` // "80:00:00:00:01:00:40:C4"
	DeviceType    int    `json:"devicetype,omitempty"`    // 28
	HostIP        string `json:"hostip,omitempty"`        // "255.255.255.255"
	Wait          string `json:"wait,omitempty"`          // "true"
	Isolate       string `json:"isolate,omitempty"`       // "domain=app.remote.it"
	Concurrent    bool   `json:"concurrent,omitempty"`    // FIXME: will this be remvoed in the future ?
	ProxyType     string `json:"proxyType,omitempty"`     // proxyType=port
}
