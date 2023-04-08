# Golang OAuth 2.0 Client
<p align="center">
<img align="center" width="150px" src="https://www.oauth.com/wp-content/themes/oauthdotcom/images/oauth_logo@2x.png" />
</p>
<h3 align=center>Golang 实现的 OAuth2.0 客户端</h3>

[English](README.md) | 简体中文 | [Oauth2 流程介绍](oauth-flow-cn.md)

## OAuth 2.0协议流程
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+

## OAuth 客户端的实现及其特点

- 提供生成授权码URL
- 通过 OAuth2 服务器的授权码提供获取 AccessToken
- 通过AccessToken提供获取用户信息
- 通过AccessToken刷新 AccessToken 功能
- 注销 AccessToken 功能

## 安装

在你的项目中使用下面的命令安装:

`go get -u github.com/demo007x/oauth2-client`

## 快速开始

以下示例提供了一个github授权的示例代码:

```go
package main

import (
	"github.com/demo007x/oauth2-client/oauth"
	"log"
	"net/http"
	"net/url"
)

// This Is GitHub.com Oauth Restfull Demo
var (
	clientID    = "567bcc7f346c8ce22e1893cee0f43a3a" // change youself clientID
	secret      = "a4a2d532e29a262a8fc67bc5e4db01be"
	serverURL   = "https://github.com/login/oauth/authorize"
	redirectURL = "http://127.0.0.1:8080/oauth/callback"
	scope       = "user read:user"
	state       = "xxxx"
)

func handler(w http.ResponseWriter, r *http.Request) {
	githubClient := oauth.NewOauth2Client(serverURL, clientID, oauth.WithRedirectURI(redirectURL), oauth.WithState(state), oauth.WithScope(scope))
	authURL, err := githubClient.AuthorizeURL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, authURL, http.StatusFound)
	return
}

func callback(w http.ResponseWriter, r *http.Request) {
	var serverURL = "https://github.com/login/oauth/access_token"
	u, _ := url.ParseRequestURI(r.RequestURI)
	var code = u.Query().Get("code")
	log.Println("code = ", code)
	// get access token by code
	accessToken := oauth.NewAccessToken(serverURL, clientID, secret, code, oauth.AccessTokenWithContentType("application/json"))
	data, err := accessToken.DoRequest()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	getUserinfo(w, string(data))
}

func getUserinfo(w http.ResponseWriter, requestURI string) {
	//access_token=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&scope=user&token_type=bearer
	values, _ := url.ParseQuery(requestURI)
	var accessToken = values.Get("access_token")
	var serverURL = "https://api.github.com/user"
	user := oauth.NewUserInfo(serverURL, accessToken)
	data, err := user.DoRequest()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth/callback", callback)

	http.ListenAndServe(":8080", nil)
}
```

## Give a Star! ⭐
如果你喜欢或正在使用这个项目来学习或开始你的解决方案，请给它一颗星。谢谢！

## Buy me a coffee

<a href="https://www.buymeacoffee.com/demo007x" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>

## 问题讨论
<img src="https://user-images.githubusercontent.com/6418340/230716042-135d28f0-9912-4ba4-8adf-a8f14eb76b05.png" alt="discard with Me" style="height: 150px !important;width: 150px !important;" >


