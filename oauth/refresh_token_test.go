package oauth

import (
	"io"
	"net/http"
	"testing"
)

func TestNewOauthRefreshToken(t *testing.T) {
	var serverURL = "http://localhost:8200/api/v1/oauth2/refresh_token"
	var grantType = "refresh_token"
	var key = "A9GzcBd1Qt3jbv4YBPOEHb1xKXDyNBBIFiK"
	var secret = "aOh4131NG3odJWF3o1c2TWncpwy2qgoCuNxcS5uw"
	var refreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODM0NjAwNzQsImlhdCI6MTY4MDg2ODA3NCwidXVpZCI6InVzZXItMGY2N2EwMmU4LTU3OTUtNDAzNC05Njg2LWM4YWIzNTEwNWU2MiJ9.X5lE52dzMT9_LqpzVKsiZW_D3tSUixZ42mwNz-UL--Q"
	token := NewRefreshToken(serverURL, key, secret, refreshToken, RefreshTokenWithGrantType(grantType), RefreshTokenWithResponseHandler(func(resp *http.Response) ([]byte, error) {
		defer func() {
			resp.Body.Close()
		}()
		return io.ReadAll(resp.Body)
	}))
	data, err := token.DoRequest()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}
