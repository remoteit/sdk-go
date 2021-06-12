package httpmessagesigner

import (
	"net/http"
)

const (
	requestDateHeader          = "date"
	requestSignatureHeader     = "Signature"
	requestAuthorizationHeader = "Authorization"
	requestTargetHeader        = "(request-target)"
	requestAuthScheme          = "Signature "
)

type Signer interface {
	Sign(req *http.Request)
}
