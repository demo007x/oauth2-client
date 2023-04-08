package types

import (
	"io"
	"log"
	"net/http"
)

var (
	// DefaultAuthorizeResponseType default authorize response type is code
	DefaultAuthorizeResponseType = "code"
	// DefaultAccessTokenGrantType default access token grant type is authorization_code
	DefaultAccessTokenGrantType = "authorization_code"
)

// OauthResponseHandler oauth2 server response handler
type OauthResponseHandler func(resp *http.Response) ([]byte, error)

func DefaultOauthResponseHandler(resp *http.Response) ([]byte, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return io.ReadAll(resp.Body)
}
