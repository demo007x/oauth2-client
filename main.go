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
