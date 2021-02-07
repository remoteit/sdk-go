package api

import (
	"encoding/json"

	apiContracts "github.com/remoteit/sdk-go/contracts"
	errorx "github.com/remoteit/systemkit-errorx"
)

// DeleteProxyRequest -
type DeleteProxyRequest struct {
	DeviceAddress string `json:"deviceaddress,omitempty"` // "80:00:00:00:01:00:40:C4"
	ConnectionID  string `json:"connectionid,omitempty"`  // "118EB6E7-1BA5-0AA0-B80E-7D15B9059260"
}

// NewDeleteProxyRequest -
func NewDeleteProxyRequest(deviceAddress string, connectionID string) DeleteProxyRequest {
	return DeleteProxyRequest{
		DeviceAddress: deviceAddress,
		ConnectionID:  connectionID,
	}
}

// DeleteProxyResponse -
type DeleteProxyResponse struct {
	Status string `json:"status"` // "true"
}

// DeleteProxy - deletes a REMOTEIT proxy
func (c *Client) DeleteProxy(request DeleteProxyRequest) (DeleteProxyResponse, errorx.Error) {
	body, err := json.Marshal(request)
	if err != nil {
		return DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantPrepRequest
	}

	raw, err := c.Post("/device/connect/stop", body)
	if err != nil {
		return DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantSendRequest
	}

	var response DeleteProxyResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_CantReadResponse
	}

	if response.Status != apiContracts.API_ERROR_CODE_STATUS_TRUE {
		return DeleteProxyResponse{}, apiContracts.ErrAPI_ProxyDelete_Unknown
	}

	return response, nil
}
