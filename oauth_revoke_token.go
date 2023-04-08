package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/utils"
	"net/http"
	"net/url"
	"strings"
)

type (
	OauthRevokeTokenOption func(token *OauthRevokeToken)

	OauthRevokeToken struct {
		ServerURL     string
		Key           string
		Secret        string
		AccessToken   string
		TokenTypeHint string

		// internal field
		u       *url.URL
		values  url.Values
		header  map[string]string
		handler OauthResponseHandler
		err     error
	}
)

func OauthRevokeTokenWithServerURL(serverURL string) OauthRevokeTokenOption {
	return func(token *OauthRevokeToken) {
		token.ServerURL = serverURL
	}
}

func OauthRevokeTokenWithKeyAndSecret(key, secret string) OauthRevokeTokenOption {
	return func(token *OauthRevokeToken) {
		token.Key = key
		token.Secret = secret
	}
}

func OauthRevokeTokenWithAccessToken(accessToken string) OauthRevokeTokenOption {
	return func(token *OauthRevokeToken) {
		token.AccessToken = accessToken
	}
}

func OauthRevokeTokenWithTokenTypeHint(tokenTypeHint string) OauthRevokeTokenOption {
	return func(token *OauthRevokeToken) {
		token.TokenTypeHint = tokenTypeHint
	}
}

func OauthRevokeTokenWithResponseHandler(handler OauthResponseHandler) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.respHandler = handler
	}
}

func OauthRevokeTokenWithContentType(contentType string) OauthRevokeTokenOption {
	return func(token *OauthRevokeToken) {
		token.header["Content-Type"] = contentType
		if strings.TrimSpace(contentType) == "" {
			token.header["Content-Type"] = "application/json"
		}
	}
}

func (ort *OauthRevokeToken) setServerURL() *OauthRevokeToken {
	if ort.err == nil {
		ort.u, ort.err = url.Parse(ort.ServerURL)
		if ort.err != nil {
			ort.values = ort.u.Query()
		}
	}
	return ort
}

func (ort *OauthRevokeToken) setKeyAndSecret() *OauthRevokeToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.Key) == "" {
			ort.err = errorx.ClientKeyError
			return ort
		}
		if strings.TrimSpace(ort.Secret) == "" {
			ort.err = errorx.SecretKeyError
			return ort
		}
		// set header
		ort.header["Authorization"] = utils.GenerateBaseAuthorization(ort.Key, ort.Secret)
	}
	return ort
}

func (ort *OauthRevokeToken) setAccessToken() *OauthRevokeToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.Key) == "" {
			ort.err = errorx.ClientKeyError
			return ort
		}
		ort.values.Set("token", ort.AccessToken)
	}
	return ort
}

func (ort *OauthRevokeToken) setTokenTypeHint() *OauthRevokeToken {
	if ort.err == nil {
		ort.values.Set("token_type_hint", ort.TokenTypeHint)
		if strings.TrimSpace(ort.TokenTypeHint) == "" {
			ort.values.Set("token_type_hint", "access_token")
		}
	}
	return ort
}

func (ort *OauthRevokeToken) DoRequest() ([]byte, error) {
	if err := ort.setServerURL().
		setKeyAndSecret().
		setAccessToken().
		setTokenTypeHint().
		err; err != nil {
		return nil, err
	}
	ort.u.RawQuery = ort.values.Encode()
	var requestURL = ort.u.String()
	resp, err := utils.DoRequest(requestURL, http.MethodPost, ort.header)
	if err != nil {
		return nil, err
	}
	if ort.handler == nil {
		return defaultOauthResponseHandler(resp)
	}
	return ort.handler(resp)
}

func NewOauthRevokeToken(serverURL, key, secret, accessToken string, opts ...OauthRevokeTokenOption) *OauthRevokeToken {
	var token = &OauthRevokeToken{
		header: make(map[string]string),
		values: url.Values{},
	}
	opts = append(opts, OauthRevokeTokenWithServerURL(serverURL), OauthRevokeTokenWithKeyAndSecret(key, secret), OauthRevokeTokenWithAccessToken(accessToken), OauthRevokeTokenWithContentType("application/json"))
	for _, opt := range opts {
		opt(token)
	}
	return token
}
