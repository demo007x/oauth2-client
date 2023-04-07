package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/types"
	"github.com/anziguoer/oauth2-client/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	AccessTokenWithOption  func(ac *OauthAccessToken)
	AccessTokenRespHandler func(resp *http.Response) ([]byte, error)
	OauthAccessToken       struct {
		ServerURL   string
		Key         string
		Secret      string
		Code        string
		GrantType   string
		RedirectURI string
		ContentType string
		// Internal field
		respHandler AccessTokenRespHandler
		sup         *url.URL
		values      url.Values
		header      map[string]string
		err         error
	}
)

func defaultAccessTokenRespHandler(resp *http.Response) ([]byte, error) {
	defer func() {
		resp.Body.Close()
	}()
	return io.ReadAll(resp.Body)
}

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

// AccessTokenWithResponseHandler
// Custom access token handle. Response from server with call AccessTokenRespHandler
func AccessTokenWithResponseHandler(handler AccessTokenRespHandler) AccessTokenWithOption {
	return func(ac *OauthAccessToken) {
		ac.respHandler = handler
	}
}

// verify server uri
func (ac *OauthAccessToken) verifyServerURI() *OauthAccessToken {
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

// verify key and secret
func (ac *OauthAccessToken) verifyKeyAndSecret() *OauthAccessToken {
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

// verify grant type. if empty set default
func (ac *OauthAccessToken) verifyGrantType() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.GrantType) != "" {
			ac.GrantType = types.DefaultAccessTokenGrantType
			ac.values.Set("grant_type", ac.GrantType)
		}
	}
	return ac
}

// verify code
func (ac *OauthAccessToken) verifyCode() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.Code) == "" {
			ac.err = errorx.CodeEmptyError
			return ac
		}
		ac.values.Set("code", ac.Code)
	}
	return ac
}

// verify redirect uri
func (ac *OauthAccessToken) verifyRedirectURI() *OauthAccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.RedirectURI) != "" {
			ac.values.Set("redirect_url", ac.RedirectURI)
		}
	}
	return ac
}

// DoRequest request access token from oauth server
func (ac *OauthAccessToken) DoRequest() ([]byte, error) {
	if err := ac.verifyServerURI().
		verifyKeyAndSecret().
		verifyGrantType().
		verifyCode().
		verifyRedirectURI().
		err; err != nil {
		return nil, ac.err
	}

	ac.sup.RawQuery = ac.values.Encode()
	var requestHost = ac.sup.String()

	resp, err := utils.DoRequest(requestHost, http.MethodPost, ac.header)
	if err != nil {
		return nil, err
	}

	if ac.respHandler == nil {
		ac.respHandler = defaultAccessTokenRespHandler
	}
	// handler response
	return ac.respHandler(resp)
}

// NewOauthAccessToken return OauthAccessToken implement
func NewOauthAccessToken(serverURL, key, secret, code string, opts ...AccessTokenWithOption) *OauthAccessToken {
	var OauthAccessToken = &OauthAccessToken{
		ServerURL: serverURL,
		Key:       key,
		Secret:    secret,
		Code:      code,
	}
	OauthAccessToken.header = make(map[string]string)
	for _, opt := range opts {
		opt(OauthAccessToken)
	}
	return OauthAccessToken
}
