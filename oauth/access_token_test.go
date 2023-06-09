package oauth

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestNewAccessToken(t *testing.T) {
	var serverURL = "http://localhost:8200/api/v1/oauth2/token"
	var redirectURI = "http://localhost:8200/"
	var username = "A9GzcBd1Qt3jbv4YBPOEHb1xKXDyNBBIFiK"
	var password = "aOh4131NG3odJWF3o1c2TWncpwy2qgoCuNxcS5uw"
	var code = "OGNJNTIXNGMTZWFKNC0ZZDJKLWJMMWUTYMY3NTYXOTU4ZTC2"
	token := NewAccessToken(serverURL, redirectURI, username, password, code, AccessTokenWithGrantType("authorization_code"), AccessTokenWithContentType("application/json"), AccessTokenWithResponseHandler(func(resp *http.Response) ([]byte, error) {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}))
	data, err := token.DoRequest()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}
