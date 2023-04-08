package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/utils"
	"net/http"
	"net/url"
	"strings"
)

type (
	OauthRefreshTokenOption func(token *OauthRefreshToken)
	OauthRefreshToken       struct {
		ServerURL    string
		Key          string
		Secret       string
		RefreshToken string
		GrantType    string
		ContentType  string
		// internal field
		respHandler OauthResponseHandler
		u           *url.URL
		header      map[string]string
		values      url.Values
		err         error
	}
)

func OauthRefreshTokenWithGrantType(grantType string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.GrantType = grantType
	}
}

func OauthRefreshTokenWithRefreshToken(refreshToken string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.RefreshToken = refreshToken
	}
}

func OauthRefreshTokenWithKey(key string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.Key = key
	}
}

func OauthRefreshTokenWithSecret(secret string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.Secret = secret
	}
}

func OauthRefreshTokenWithServerURL(serverURL string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.ServerURL = serverURL
	}
}

func RefreshTokenWithResponseHandler(handler OauthResponseHandler) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.respHandler = handler
	}
}

func OauthRefreshTokenWithContentType(contentType string) OauthRefreshTokenOption {
	return func(token *OauthRefreshToken) {
		token.ContentType = contentType
		if strings.TrimSpace(contentType) == "" {
			token.ContentType = "application/json"
		}
		token.header["Content-Type"] = contentType
	}
}

// setServerURI
// todo 统一处理 Oauth 服务的校验
func (ort *OauthRefreshToken) setServerURI() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.ServerURL) == "" {
			ort.err = errorx.ServerURLError
			return ort
		}
		ort.u, ort.err = url.Parse(ort.ServerURL)
		if ort.err == nil {
			ort.values = ort.u.Query()
		}
	}

	return ort
}

func (ort *OauthRefreshToken) setRefreshToken() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.RefreshToken) == "" {
			ort.err = errorx.RefreshTokenNotEmpty
			return ort
		}
		ort.values.Set("refresh_token", ort.RefreshToken)
	}

	return ort
}

func (ort *OauthRefreshToken) setGrantType() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.GrantType) == "" {
			ort.GrantType = "refresh_token"
		}
		ort.values.Set("grant_type", ort.GrantType)
	}

	return ort
}

func (ort *OauthRefreshToken) setKeyAndSecret() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.Key) == "" || strings.TrimSpace(ort.Secret) == "" {
			ort.err = errorx.SecretKeyError
			return ort
		}
		ort.header["Authorization"] = utils.GenerateBaseAuthorization(ort.Key, ort.Secret)
	}

	return ort
}

func (ort *OauthRefreshToken) DoRequest() ([]byte, error) {
	if err := ort.setServerURI().
		setRefreshToken().
		setKeyAndSecret().
		err; err != nil {
		return nil, err
	}
	ort.u.RawQuery = ort.values.Encode()
	var requestHost = ort.u.String()
	resp, err := utils.DoRequest(requestHost, http.MethodPost, ort.header)
	if err != nil {
		return nil, err
	}
	if ort.respHandler == nil {
		return defaultOauthResponseHandler(resp)
	}
	return ort.respHandler(resp)
}

func NewOauthRefreshToken(serverURL, key, secret, refreshToken string, opts ...OauthRefreshTokenOption) *OauthRefreshToken {
	var token = &OauthRefreshToken{}
	token.header = make(map[string]string)
	opts = append(opts, OauthRefreshTokenWithRefreshToken(refreshToken), OauthRefreshTokenWithKey(key), OauthRefreshTokenWithSecret(secret), OauthRefreshTokenWithServerURL(serverURL))
	for _, opt := range opts {
		opt(token)
	}

	return token
}
