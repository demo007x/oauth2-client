package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/utils"
	"net/http"
	"net/url"
	"strings"
)

type (
	AccessTokenWithOption func(ac *OauthAccessToken)
	OauthAccessToken      struct {
		ServerURL   string
		Key         string
		Secret      string
		Code        string
		GrantType   string
		RedirectURI string
		ContentType string
		// Internal field
		handler OauthResponseHandler
		sup     *url.URL
		values  url.Values
		header  map[string]string
		err     error
	}
)

// AccessTokenWithGrantType
// Config OauthAccessToken with grant type
func AccessTokenWithGrantType(grantType string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.GrantType = grantType
	}
}

// AccessTokenWithRedirectURI
// Config OauthAccessToken with redirect uri
func AccessTokenWithRedirectURI(redirectURI string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.RedirectURI = redirectURI
	}
}

// AccessTokenWithContentType set content type
func AccessTokenWithContentType(contentType string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.ContentType = contentType
		if strings.TrimSpace(contentType) == "" {
			ac.ContentType = "application/json"
		}
		ac.header["Content-Type"] = contentType
	}
}

func AccessTokenWithServerURL(serverURL string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.ServerURL = serverURL
	}
}

func AccessTokenWithKeyAndSecret(key, secret string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.Key, ac.Secret = key, secret
	}
}

func AccessTokenWithCode(code string) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.Code = code
	}
}

// AccessTokenWithResponseHandler
// Custom access token handle. Response from server with call AccessTokenRespHandler
func AccessTokenWithResponseHandler(handler OauthResponseHandler) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.handler = handler
	}
}

// set server uri
func (ac *OauthAccessToken) setServerURI() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.ServerURL) == "" {
			ac.err = errorx.ServerURLError
			return ac
		}
		ac.sup, ac.err = url.Parse(ac.ServerURL)
		if ac.err == nil {
			ac.values = ac.sup.Query()
		}
	}
	return ac
}

// set key and secret
func (ac *OauthAccessToken) setKeyAndSecret() *OauthAccessToken {
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
func (ac *OauthAccessToken) setGrantType() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.GrantType) != "" {
			ac.GrantType = DefaultAccessTokenGrantType
			ac.values.Set("grant_type", ac.GrantType)
		}
	}
	return ac
}

// set code
func (ac *OauthAccessToken) setCode() *OauthAccessToken {
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
func (ac *OauthAccessToken) setRedirectURI() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.RedirectURI) != "" {
			ac.values.Set("redirect_url", ac.RedirectURI)
		}
	}
	return ac
}

// DoRequest request access token from oauth server
func (ac *OauthAccessToken) DoRequest() ([]byte, error) {
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
		return defaultOauthResponseHandler(resp)
	}
	// handler response
	return ac.handler(resp)
}

// NewOauthAccessToken return OauthAccessToken implement
func NewOauthAccessToken(serverURL, redirectURI, key, secret, code string, opts ...AccessTokenWithOption) *OauthAccessToken {
	var OauthAccessToken = &OauthAccessToken{}
	OauthAccessToken.header = make(map[string]string)
	opts = append(opts, AccessTokenWithCode(code), AccessTokenWithKeyAndSecret(key, secret), AccessTokenWithServerURL(serverURL), AccessTokenWithRedirectURI(redirectURI))
	for _, opt := range opts {
		opt(OauthAccessToken)
	}
	return OauthAccessToken
}
