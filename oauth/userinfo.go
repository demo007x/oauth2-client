package oauth

import (
	"github.com/demo007x/oauth2-client/errorx"
	"github.com/demo007x/oauth2-client/types"
	"github.com/demo007x/oauth2-client/utils"
	"net/http"
	"net/url"
)

type (
	WithUserInfoOption func(info *UserInfo)
	UserInfo           struct {
		AccessToken string
		ServerURL   string

		// internal field
		handler types.OauthResponseHandler
		header  map[string]string
		err     error
	}
)

func userInfoWithAccessToken(token string) WithUserInfoOption {
	return func(info *UserInfo) {
		info.AccessToken = token
	}
}

func userInfoWithServerURL(serverURL string) WithUserInfoOption {
	return func(info *UserInfo) {
		info.ServerURL = serverURL
	}
}

func UserInfoWithResponseHandler(handler types.OauthResponseHandler) WithUserInfoOption {
	return func(info *UserInfo) {
		info.handler = handler
	}
}

// setServerURL set server url invalid
// todo 统一url的验证函数
func (info *UserInfo) setServerURL() *UserInfo {
	if info.err == nil {
		_, err := url.Parse(info.ServerURL)
		info.err = err
	}
	return info
}

func (info *UserInfo) setToken() *UserInfo {
	if info.err == nil {
		info.header["Authorization"] = utils.GenerateBearAuthorization(info.AccessToken)
	}
	return info
}

// DoRequest request oauth server get user info
func (info *UserInfo) DoRequest() ([]byte, error) {
	if err := info.setServerURL().setToken().err; err != nil {
		return nil, err
	}
	resp, err := utils.DoRequest(info.ServerURL, http.MethodPost, info.header)
	if err != nil {
		return nil, errorx.RequestServerURLError
	}
	if info.handler == nil {
		return types.DefaultOauthResponseHandler(resp)
	}
	return info.handler(resp)
}

func NewUserInfo(serverURL, accessToken string, opts ...WithUserInfoOption) *UserInfo {
	var info = &UserInfo{}
	opts = append(opts, userInfoWithServerURL(serverURL), userInfoWithAccessToken(accessToken))
	info.header = make(map[string]string, 1)
	for _, opt := range opts {
		opt(info)
	}
	return info
}
