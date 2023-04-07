package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"github.com/anziguoer/oauth2-client/utils"
	"io"
	"net/http"
	"net/url"
)

type (
	WithOauthUserInfoOption      func(info *OauthUserInfo)
	OauthUserInfoResponseHandler func(resp *http.Response) ([]byte, error)
	OauthUserInfo                struct {
		AccessToken string
		ServerURL   string

		// internal field
		respHandler OauthUserInfoResponseHandler
		header      map[string]string
		err         error
	}

	OauthUserInfoResp struct {
		Id       string `json:"id"`
		UserName string `json:"userName"`
		Mobile   string `json:"mobile"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}
)

func defaultUserInfoHandler(resp *http.Response) ([]byte, error) {
	defer func() {
		resp.Body.Close()
	}()
	return io.ReadAll(resp.Body)
}

func OauthUserInfoWithAccessToken(token string) WithOauthUserInfoOption {
	return func(info *OauthUserInfo) {
		info.AccessToken = token
	}
}

func OauthUserInfoWithServerURL(serverURL string) WithOauthUserInfoOption {
	return func(info *OauthUserInfo) {
		info.ServerURL = serverURL
	}
}

func OauthUserInfoWithResponseHandler(handler OauthUserInfoResponseHandler) WithOauthUserInfoOption {
	return func(info *OauthUserInfo) {
		info.respHandler = handler
	}
}

// verifyServerURL verify server url invalid
// todo 统一url的验证函数
func (info *OauthUserInfo) verifyServerURL() *OauthUserInfo {
	if info.err == nil {
		_, err := url.Parse(info.ServerURL)
		info.err = err
	}
	return info
}

func (info *OauthUserInfo) verifyToken() *OauthUserInfo {
	if info.err == nil {
		info.header["Authorization"] = utils.GenerateBearAuthorization(info.AccessToken)
	}
	return info
}

// DoRequest request oauth server get user info
func (info *OauthUserInfo) DoRequest() (interface{}, error) {
	if err := info.verifyServerURL().verifyToken().err; err != nil {
		return nil, err
	}
	resp, err := utils.DoRequest(info.ServerURL, http.MethodPost, info.header)
	if err != nil {
		return nil, errorx.RequestServerURLError
	}
	if info.header == nil {
		return defaultUserInfoHandler(resp)
	}
	return info.respHandler(resp)
}

func NewOauthUserInfo(serverURL, accessToken string, opts ...WithOauthUserInfoOption) *OauthUserInfo {
	var info = &OauthUserInfo{
		ServerURL:   serverURL,
		AccessToken: accessToken,
	}
	info.header = make(map[string]string, 1)
	for _, opt := range opts {
		opt(info)
	}
	return info
}
