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
	RevokeTokenOption func(token *RevokeToken)

	RevokeToken struct {
		ServerURL     string
		ClientID      string
		Secret        string
		AccessToken   string
		TokenTypeHint string

		// internal field
		u       *url.URL
		values  url.Values
		header  map[string]string
		handler types.OauthResponseHandler
		err     error
	}
)

func RevokeTokenWithServerURL(serverURL string) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.ServerURL = serverURL
	}
}

func RevokeTokenWithKeyAndSecret(clientID, secret string) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.ClientID = clientID
		token.Secret = secret
	}
}

func RevokeTokenWithAccessToken(accessToken string) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.AccessToken = accessToken
	}
}

func RevokeTokenWithTokenTypeHint(tokenTypeHint string) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.TokenTypeHint = tokenTypeHint
	}
}

func RevokeTokenWithResponseHandler(handler types.OauthResponseHandler) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.handler = handler
	}
}

func RevokeTokenWithContentType(contentType string) RevokeTokenOption {
	return func(token *RevokeToken) {
		token.header["Content-Type"] = contentType
		if strings.TrimSpace(contentType) == "" {
			token.header["Content-Type"] = "application/json"
		}
	}
}

func (ort *RevokeToken) setServerURL() *RevokeToken {
	if ort.err == nil {
		ort.u, ort.err = url.Parse(ort.ServerURL)
		if ort.err == nil {
			ort.values = ort.u.Query()
		}
	}
	return ort
}

func (ort *RevokeToken) setKeyAndSecret() *RevokeToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.ClientID) == "" {
			ort.err = errorx.ClientKeyError
			return ort
		}
		if strings.TrimSpace(ort.Secret) == "" {
			ort.err = errorx.SecretKeyError
			return ort
		}
		// set header
		ort.header["Authorization"] = utils.GenerateBaseAuthorization(ort.ClientID, ort.Secret)
	}
	return ort
}

func (ort *RevokeToken) setAccessToken() *RevokeToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.ClientID) == "" {
			ort.err = errorx.ClientKeyError
			return ort
		}
		ort.values.Set("token", ort.AccessToken)
	}
	return ort
}

func (ort *RevokeToken) setTokenTypeHint() *RevokeToken {
	if ort.err == nil {
		ort.values.Set("token_type_hint", ort.TokenTypeHint)
		if strings.TrimSpace(ort.TokenTypeHint) == "" {
			ort.values.Set("token_type_hint", "access_token")
		}
	}
	return ort
}

func (ort *RevokeToken) DoRequest() ([]byte, error) {
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
		return types.DefaultOauthResponseHandler(resp)
	}
	return ort.handler(resp)
}

func NewOauthRevokeToken(serverURL, key, secret, accessToken string, opts ...RevokeTokenOption) *RevokeToken {
	var token = &RevokeToken{
		header: make(map[string]string),
		values: url.Values{},
	}
	opts = append(opts, RevokeTokenWithServerURL(serverURL), RevokeTokenWithKeyAndSecret(key, secret), RevokeTokenWithAccessToken(accessToken), RevokeTokenWithContentType("application/json"))
	for _, opt := range opts {
		opt(token)
	}
	return token
}
