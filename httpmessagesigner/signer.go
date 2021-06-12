package httpmessagesigner

import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"hash"
	"net/http"
	"strings"
	"time"
)

func NewSigner(headers []string, key string, secret string, hasherAlgo func() hash.Hash) Signer {
	return &signer{
		headers:    headers,
		key:        key,
		secret:     secret,
		hasherAlgo: hasherAlgo,
	}
}

type signer struct {
	headers    []string
	key        string
	secret     string
	hasherAlgo func() hash.Hash
}

func (thisRef signer) Sign(req *http.Request) {
	// update some headers before sign
	req.Header.Del(requestDateHeader)
	req.Header.Set(requestDateHeader, time.Now().Format(time.RFC1123))

	// sign
	headersToSign := strings.Builder{}
	for _, header := range thisRef.headers {
		header = strings.TrimSpace(strings.ToLower(header))

		headerLine := ""
		if header == requestTargetHeader {
			var url string = ""
			if req.URL != nil {
				url = req.URL.RequestURI()
			}

			headerLine = fmt.Sprintf("%s: %s %s", requestTargetHeader, strings.ToLower(req.Method), url)

		} else {
			headerLine = fmt.Sprintf("%s: %s", header, req.Header.Get(header))
		}

		headersToSign.WriteString(headerLine + "\n")
	}

	hasher := hmac.New(thisRef.hasherAlgo, []byte(thisRef.secret))
	hasher.Write([]byte(headersToSign.String()))
	signedHeadersAsBase64 := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// build signed line
	signedLine := fmt.Sprintf(
		`key="%s",signature="%s"`, // `key="%s",algorithm="%s",signature="%s"`
		thisRef.key,
		signedHeadersAsBase64,
	)

	// set
	req.Header.Del(requestSignatureHeader)
	req.Header.Add(requestSignatureHeader, signedLine)

	req.Header.Del(requestAuthorizationHeader)
	req.Header.Add(requestAuthorizationHeader, requestAuthScheme+signedLine)
}
