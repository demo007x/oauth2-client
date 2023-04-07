package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/utils"
	"net/http"
	"net/url"
	"strings"
)

type (
	OauthRefreshTokenOption          func(token *OauthRefreshToken)
	OauthRefreshTokenResponseHandler func(resp *http.Response) ([]byte, error)
	OauthRefreshToken                struct {
		ServerURL    string
		Key          string
		Secret       string
		RefreshToken string
		GrantType    string
		// internal field
		respHandler OauthRefreshTokenResponseHandler
		sup         *url.URL
		header      map[string]string
		values      url.Values
		err         error
	}
)

func defaultOauthRefreshTokenResponseHandler(resp *http.Response) ([]byte, error) {
	//todo please implement me
	panic("")
}

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

// verifyServerURI
// todo 统一处理 Oauth 服务的校验
func (ort *OauthRefreshToken) verifyServerURI() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.ServerURL) == "" {
			ort.err = errorx.ServerURLError
			return ort
		}
		ort.sup, ort.err = url.Parse(ort.ServerURL)
		if ort.err == nil {
			ort.values = ort.sup.Query()
		}
	}

	return ort
}

func (ort *OauthRefreshToken) verifyRefreshToken() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.RefreshToken) == "" {
			ort.err = errorx.RefreshTokenNotEmpty
			return ort
		}
		ort.values.Set("refresh_token", ort.RefreshToken)
	}

	return ort
}

func (ort *OauthRefreshToken) verifyGrantType() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.GrantType) == "" {
			ort.GrantType = "refresh_token"
		}
		ort.values.Set("grant_type", ort.GrantType)
	}

	return ort
}

func (ort *OauthRefreshToken) verifyKeyAndSecret() *OauthRefreshToken {
	if ort.err == nil {
		if strings.TrimSpace(ort.Key) == "" || strings.TrimSpace(ort.Secret) == "" {
			ort.err = errorx.SecretKeyError
			return ort
		}
		ort.header["Authorization"] = utils.GenerateBaseAuthorization(ort.Key, ort.Secret)
	}

	return ort
}

// SetOauthRefreshTokenResponseHandler
// set OauthRefreshToken request server response handler
func (ort *OauthRefreshToken) SetOauthRefreshTokenResponseHandler(handler OauthRefreshTokenResponseHandler) *OauthRefreshToken {
	if ort.err == nil {
		ort.respHandler = handler
	}

	return ort
}

func (ort *OauthRefreshToken) DoRequest() ([]byte, error) {
	if err := ort.verifyServerURI().
		verifyRefreshToken().
		verifyKeyAndSecret().
		err; err != nil {
		return nil, err
	}
	ort.sup.RawQuery = ort.values.Encode()
	var requestHost = ort.sup.String()
	resp, err := utils.DoRequest(requestHost, http.MethodPost, ort.header)
	if err != nil {
		return nil, err
	}
	if ort.respHandler == nil {
		return defaultOauthRefreshTokenResponseHandler(resp)
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
