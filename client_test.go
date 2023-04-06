package oauth2_client

import (
	"testing"
)

func TestNewOauth2Client(t *testing.T) {
	var serverURL = ""
	var redirectURI = ""
	var client = NewOauth2Client(
		"2wLCawQ1fFhmsj0ADIQIquCLiGR6qSLA",
		WithResponseType("code"),
		WithServerURl(serverURL),
		WithRedirectURI(redirectURI),
		WithState("xxxxx"),
	)

	resp, err := client.AuthorizeURL()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
		t.Log("this is a demo")
	}
}
