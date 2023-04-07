package oauth2_client

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestNewOauthUserInfo(t *testing.T) {
	var serverURL = "http://localhost:8200/api/v1/oauth2/userinfo"
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODA4NzUyNzQsImlhdCI6MTY4MDg2ODA3NCwidXVpZCI6InVzZXItMGY2N2EwMmU4LTU3OTUtNDAzNC05Njg2LWM4YWIzNTEwNWU2MiJ9.bpMUDhwsW6pXbwYpuByObA2iVb9b-NPXpg-DqBB5S94"
	user := NewOauthUserInfo(serverURL, token, OauthUserInfoWithResponseHandler(func(resp *http.Response) ([]byte, error) {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, err
	}))
	data, err := user.DoRequest()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data.([]byte)))
}
