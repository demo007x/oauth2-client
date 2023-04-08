package oauth

import (
	"testing"
)

func TestNewOauth2Client(t *testing.T) {
	var serverURL = "http://127.0.0.1:8200/auth/authorize"
	var redirectURI = ""
	var client = NewOauth2Client(
		serverURL,
		"2wLCawQ1fFhmsj0ADIQIquCLiGR6qSLA",
		WithResponseType("code"),
		WithRedirectURI(redirectURI),
		WithState("xxxxx"),
	)

	resp, err := client.AuthorizeURL()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
