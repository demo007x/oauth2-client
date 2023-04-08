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
	RefreshTokenOption func(token *RefreshToken)
	RefreshToken       struct {
		ServerURL    string
		ClientID     string
		Secret       string
		RefreshToken string
		GrantType    string
		ContentType  string
		// internal field
		respHandler types.OauthResponseHandler
		u           *url.URL
		header      map[string]string
		values      url.Values
		err         error
	}
)

func RefreshTokenWithResponseHandler(handler types.OauthResponseHandler) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.respHandler = handler
	}
}

func RefreshTokenWithGrantType(grantType string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.GrantType = grantType
	}
}

func RefreshTokenWithRefreshToken(refreshToken string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.RefreshToken = refreshToken
	}
}

func RefreshTokenWithKey(clientID string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.ClientID = clientID
	}
}

func RefreshTokenWithSecret(secret string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.Secret = secret
	}
}

func RefreshTokenWithServerURL(serverURL string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.ServerURL = serverURL
	}
}

func RefreshTokenWithContentType(contentType string) RefreshTokenOption {
	return func(token *RefreshToken) {
		token.ContentType = contentType
		if strings.TrimSpace(contentType) == "" {
			token.ContentType = "application/json"
		}
		token.header["Content-Type"] = contentType
	}
}

// setServerURI
// todo 统一处理 Oauth 服务的校验
func (ort *RefreshToken) setServerURI() *RefreshToken {
	if ort.err == nil {
		ort.u, ort.err = url.Parse(ort.ServerURL)
		if ort.err == nil {
			ort.values = ort.u.Query()
		}
	}

	return ort
}

func (ort *RefreshToken) setRefreshToken() *RefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.RefreshToken) == "" {
			ort.err = errorx.RefreshTokenNotEmpty
			return ort
		}
		ort.values.Set("refresh_token", ort.RefreshToken)
	}

	return ort
}

func (ort *RefreshToken) setGrantType() *RefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.GrantType) == "" {
			ort.GrantType = "refresh_token"
		}
		ort.values.Set("grant_type", ort.GrantType)
	}

	return ort
}

func (ort *RefreshToken) setKeyAndSecret() *RefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.ClientID) == "" || strings.TrimSpace(ort.Secret) == "" {
			ort.err = errorx.SecretKeyError
			return ort
		}
		ort.header["Authorization"] = utils.GenerateBaseAuthorization(ort.ClientID, ort.Secret)
	}

	return ort
}

func (ort *RefreshToken) DoRequest() ([]byte, error) {
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
		return types.DefaultOauthResponseHandler(resp)
	}
	return ort.respHandler(resp)
}

func NewRefreshToken(serverURL, key, secret, refreshToken string, opts ...RefreshTokenOption) *RefreshToken {
	var token = &RefreshToken{}
	token.header = make(map[string]string)
	opts = append(opts, RefreshTokenWithRefreshToken(refreshToken), RefreshTokenWithKey(key), RefreshTokenWithSecret(secret), RefreshTokenWithServerURL(serverURL))
	for _, opt := range opts {
		opt(token)
	}

	return token
}
