# Golang OAuth 2.0 Client

<p align="center">
<img align="center" width="150px" src="https://www.oauth.com/wp-content/themes/oauthdotcom/images/oauth_logo@2x.png" />
</p>
<h3 align=center>Oauth2 Client Package For Golang</h3>

English | [简体中文](README-CN.md) | [Oauth2 Flow](oauth-flow.md)

## OAuth 2.0 Protocol Flow
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

## The implementation and features of oauth client

- Provide Generate Authorization Code Url
- Provide Get Access Token By Authorization Code From Oauth2 Server
- Provide Get User Info By Access Token
- Provide Refresh Token By Access Token 
- Provide Inject Token Ability

## Installation

Run the following command under your project:

`go get -u github.com/demo007x/oauth2-client`

## Quick Start

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
If you like or are using this project to learn or start your solution, please give it a star. Thanks!

## Buy me a coffee

<a href="https://www.buymeacoffee.com/demo007x" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
