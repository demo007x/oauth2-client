package oauth2_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/types"
	"github.com/anziguoer/oauth2-client/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	AccessTokenWithOption  func(ac *AccessToken)
	AccessTokenRespHandler func(resp *http.Response) (interface{}, error)
	AccessToken            struct {
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

	// AccessTokenResp  request access token response
	AccessTokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
)

func defaultAccessTokenRespHandler(resp *http.Response) (interface{}, error) {
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("the server response code %d, data %s", resp.StatusCode, string(data)))
	}

	if err != nil {
		return nil, err
	}

	var atr AccessTokenResp
	if err := json.Unmarshal(data, &atr); err != nil {
		return nil, err
	}
	return atr, nil
}

func AccessTokenWithGrantType(grantType string) AccessTokenWithOption {
	return func(ac *AccessToken) {
		ac.GrantType = grantType
	}
}

// AccessTokenWithRedirectURI
// Config AccessToken with redirect uri
func AccessTokenWithRedirectURI(redirectURI string) AccessTokenWithOption {
	return func(ac *AccessToken) {
		ac.RedirectURI = redirectURI
	}
}

// AccessTokenWithContentType set content type
func AccessTokenWithContentType(contentType string) AccessTokenWithOption {
	return func(ac *AccessToken) {
		ac.ContentType = contentType
		ac.header["Content-Type"] = contentType
	}
}

// AccessTokenWithResponseHandler
// Custom access token handle. Response from server with call AccessTokenRespHandler
func AccessTokenWithResponseHandler(handler AccessTokenRespHandler) AccessTokenWithOption {
	return func(ac *AccessToken) {
		ac.respHandler = handler
	}
}

// verify server uri
func (ac *AccessToken) verifyServerURI() *AccessToken {
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
func (ac *AccessToken) verifyKeyAndSecret() *AccessToken {
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
		ac.header["Authorization"] = utils.GenerateBaseAuth(ac.Key, ac.Secret)
	}
	return ac
}

// verify grant type. if empty set default
func (ac *AccessToken) verifyGrantType() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.GrantType) == "" {
			ac.GrantType = types.DefaultAccessTokenGrantType
			ac.values.Set("grant_type", ac.GrantType)
		}
	}
	return ac
}

// verify code
func (ac *AccessToken) verifyCode() *AccessToken {
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
func (ac *AccessToken) verifyRedirectURI() *AccessToken {
	if ac.err == nil {
		if strings.TrimSpace(ac.RedirectURI) != "" {
			ac.values.Set("redirect_url", ac.RedirectURI)
		}
	}
	return ac
}

// DoRequest request access token from oauth server
func (ac *AccessToken) DoRequest() (interface{}, error) {
	if err := ac.verifyServerURI().
		verifyKeyAndSecret().
		verifyGrantType().
		verifyCode().
		verifyRedirectURI().
		err; err != nil {
		return ac, ac.err
	}

	ac.sup.RawQuery = ac.values.Encode()
	var requestHost = ac.sup.String()

	resp, err := utils.DoRequest(requestHost, http.MethodPost, ac.header)
	if err != nil {
		return ac, err
	}

	if ac.respHandler == nil {
		ac.respHandler = defaultAccessTokenRespHandler
	}
	// handler response
	return ac.respHandler(resp)
}

// NewAccessToken return accessToken implement
func NewAccessToken(serverURL, key, secret, code string, opts ...AccessTokenWithOption) *AccessToken {
	var accessToken = &AccessToken{
		ServerURL: serverURL,
		Key:       key,
		Secret:    secret,
		Code:      code,
	}
	for _, opt := range opts {
		opt(accessToken)
	}
	accessToken.header = make(map[string]string)
	return accessToken
}
