package oauth

import (
	"github.com/demo007x/oauth2-client/errorx"
	"github.com/demo007x/oauth2-client/types"
	"net/url"
	"strings"
)

type (
	Client struct {
		ServerURL    string
		ClientID     string
		Secret       string
		RedirectURI  string
		State        string
		Scope        string
		ResponseType string
		// internal filed
		u      *url.URL
		values url.Values
		err    error
	}

	// WithOption config option with oauth client field
	WithOption func(client *Client)
)

// WithResponseType set client response type
func WithResponseType(responseType string) WithOption {
	return func(client *Client) {
		client.ResponseType = responseType
	}
}

// WithServerURl set oauth client request server url
func withServerURl(serverURl string) WithOption {
	return func(client *Client) {
		client.ServerURL = serverURl
	}
}

// WithRedirectURI set oauth client redirectURI field
// oauth server authorization with code redirect to redirectURI
func WithRedirectURI(redirectURI string) WithOption {
	return func(client *Client) {
		client.RedirectURI = redirectURI
	}
}

func WithScope(scope string) WithOption {
	return func(client *Client) {
		client.Scope = scope
	}
}

// WithState set oauth client request with state field
// oauth server authorization with state redirect to redirectURI
func WithState(state string) WithOption {
	return func(client *Client) {
		client.State = state
	}
}

func withClientID(clientID string) WithOption {
	return func(client *Client) {
		client.ClientID = clientID
	}
}

// parse server uri and set client suParser field
func (client *Client) setServerURI() *Client {
	if client.err == nil {
		client.u, client.err = url.Parse(client.ServerURL)
		if client.err == nil {
			client.values = client.u.Query()
		}
	}
	return client
}

func (client *Client) setRedirect() *Client {
	if client.err == nil {
		if strings.TrimSpace(client.RedirectURI) == "" {
			return client
		}
		client.values.Set("redirect_uri", client.RedirectURI)
	}
	return client
}

func (client *Client) setResponseType() *Client {
	if client.err == nil {
		responseType := client.ResponseType
		if strings.TrimSpace(responseType) == "" {
			responseType = types.DefaultAuthorizeResponseType
		}
		client.values.Set("response_type", responseType)
	}
	return client
}

func (client *Client) setScope() *Client {
	if client.err == nil {
		scope := strings.TrimSpace(client.Scope)
		if len(scope) == 0 {
			scope = "get_user_info"
		}
		client.values.Set("scope", scope)
	}
	return client
}

// The application generates a random string and includes it in the request.
// It should then check that the same value is returned after the user authorizes the app. This is used to prevent CSRF attacks.
func (client *Client) setState() *Client {
	if client.err == nil && strings.TrimSpace(client.State) != "" {
		client.values.Set("state", client.State)
	}
	return client
}

func (client *Client) setClientID() *Client {
	if client.err == nil {
		if strings.TrimSpace(client.ClientID) == "" {
			client.err = errorx.ClientKeyError
			return client
		}
		client.values.Set("client_id", client.ClientID)
	}
	return client
}

func (client *Client) AuthorizeURL() (string, error) {
	c := client.
		setServerURI().
		setRedirect().
		setResponseType().
		setScope().
		setState().
		setClientID()
	if c.err != nil {
		return "", c.err
	}
	c.u.RawQuery = client.values.Encode()
	return c.u.String(), nil
}

func NewOauth2Client(serverURL, clientID string, opts ...WithOption) *Client {
	var client = &Client{}
	opts = append(opts, withClientID(clientID), withServerURl(serverURL))
	for _, opt := range opts {
		opt(client)
	}
	return client
}
