package oauth2_client

import "net/http"

var (
	// DefaultAuthorizeResponseType default authorize response type is code
	DefaultAuthorizeResponseType = "code"
	// DefaultAccessTokenGrantType default access token grant type is authorization_code
	DefaultAccessTokenGrantType = "authorization_code"
)

// OauthResponseHandler oauth2 server response handler
type OauthResponseHandler func(resp *http.Response) ([]byte, error)
