package oauth

import (
	"github.com/demo007x/oauth2-client/errorx"
	"github.com/demo007x/oauth2-client/types"
	"github.com/demo007x/oauth2-client/utils"
	"net/http"
	"net/url"
	"strings"
)

type (
	AccessTokenOption func(ac *AccessToken)
	AccessToken       struct {
		ServerURL   string
		Key         string
		Secret      string
		Code        string
		GrantType   string
		RedirectURI string
		ContentType string
		// Internal field
		handler types.OauthResponseHandler
		sup     *url.URL
		values  url.Values
		header  map[string]string
		err     error
	}
)

// AccessTokenWithGrantType
// Config AccessToken with grant type
func AccessTokenWithGrantType(grantType string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.GrantType = grantType
	}
}

// AccessTokenWithRedirectURI
// Config AccessToken with redirect uri
func AccessTokenWithRedirectURI(redirectURI string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.RedirectURI = redirectURI
	}
}

// AccessTokenWithContentType set content type
func AccessTokenWithContentType(contentType string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.ContentType = contentType
		if strings.TrimSpace(contentType) == "" {
			ac.ContentType = "application/json"
		}
		ac.header["Content-Type"] = contentType
	}
}

func AccessTokenWithServerURL(serverURL string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.ServerURL = serverURL
	}
}

func AccessTokenWithKeyAndSecret(key, secret string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.Key, ac.Secret = key, secret
	}
}

func AccessTokenWithCode(code string) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.Code = code
	}
}

// AccessTokenWithResponseHandler
// Custom access token handle. Response from server with call AccessTokenRespHandler
func AccessTokenWithResponseHandler(handler types.OauthResponseHandler) AccessTokenOption {
	return func(ac *AccessToken) {
		ac.handler = handler
	}
}

// set server uri
func (ac *AccessToken) setServerURI() *AccessToken {
	if ac.err == nil {
		ac.sup, ac.err = url.Parse(ac.ServerURL)
		if ac.err == nil {
			ac.values = ac.sup.Query()
		}
	}
	return ac
}

// set key and secret
func (ac *AccessToken) setKeyAndSecret() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.Key) == "" {
			ac.err = errorx.ClientKeyError
			return ac
		}

		if strings.TrimSpace(ac.Secret) == "" {
			ac.err = errorx.SecretKeyError
			return ac
		}
		// set authorization
		ac.header["Authorization"] = utils.GenerateBaseAuthorization(ac.Key, ac.Secret)
	}
	return ac
}

// set grant type. if empty set default
func (ac *AccessToken) setGrantType() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.GrantType) != "" {
			ac.GrantType = types.DefaultAccessTokenGrantType
			ac.values.Set("grant_type", ac.GrantType)
		}
	}
	return ac
}

// set code
func (ac *AccessToken) setCode() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.Code) == "" {
			ac.err = errorx.CodeEmptyError
			return ac
		}
		ac.values.Set("code", ac.Code)
	}
	return ac
}

// set redirect uri
func (ac *AccessToken) setRedirectURI() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.RedirectURI) != "" {
			ac.values.Set("redirect_url", ac.RedirectURI)
		}
	}
	return ac
}

// DoRequest request access token from oauth server
func (ac *AccessToken) DoRequest() ([]byte, error) {
	if err := ac.setServerURI().
		setKeyAndSecret().
		setGrantType().
		setCode().
		setRedirectURI().
		err; err != nil {
		return nil, ac.err
	}

	ac.sup.RawQuery = ac.values.Encode()
	var requestHost = ac.sup.String()

	resp, err := utils.DoRequest(requestHost, http.MethodPost, ac.header)
	if err != nil {
		return nil, err
	}

	if ac.handler == nil {
		return types.DefaultOauthResponseHandler(resp)
	}

	return ac.handler(resp)
}

// NewAccessToken return AccessToken implement
func NewAccessToken(serverURL, redirectURI, key, secret, code string, opts ...AccessTokenOption) *AccessToken {
	var AccessToken = &AccessToken{}
	AccessToken.header = make(map[string]string)
	opts = append(opts, AccessTokenWithCode(code), AccessTokenWithKeyAndSecret(key, secret), AccessTokenWithServerURL(serverURL), AccessTokenWithRedirectURI(redirectURI))
	for _, opt := range opts {
		opt(AccessToken)
	}
	return AccessToken
}
