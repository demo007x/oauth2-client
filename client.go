package oauth2_client

import (
	"github.com/anziguoer/oauth2-client/errorx"
	"net/url"
	"strings"
)

type (
	Oauth2Client struct {
		ServerURL    string
		ResponseType string
		Key          string
		Secret       string
		RedirectURI  string
		State        string
		Scope        string
		suParser     *url.URL
		values       url.Values
		err          error
	}

	// WithOption config option with oauth client field
	WithOption func(client *Oauth2Client)
)

// WithResponseType set client response type
func WithResponseType(responseType string) WithOption {
	return func(client *Oauth2Client) {
		client.ResponseType = responseType
	}
}

// WithServerURl set oauth client request server url
func WithServerURl(serverURl string) WithOption {
	return func(client *Oauth2Client) {
		client.ServerURL = serverURl
	}
}

// WithRedirectURI set oauth client redirectURI field
// oauth server authorization with code redirect to redirectURI
func WithRedirectURI(redirectURI string) WithOption {
	return func(client *Oauth2Client) {
		client.RedirectURI = redirectURI
	}
}

func WithScope(scope string) WithOption {
	return func(client *Oauth2Client) {
		client.Scope = scope
	}
}

// WithState set oauth client request with state field
// oauth server authorization with state redirect to redirectURI
func WithState(state string) WithOption {
	return func(client *Oauth2Client) {
		client.State = state
	}
}

// parse server uri and set client suParser field
func (client *Oauth2Client) setServerURI() *Oauth2Client {
	if client.err == nil {
		if strings.TrimSpace(client.ServerURL) == "" {
			client.err = errorx.ServerURLError
			return client
		}

		parseURL, err := url.Parse(client.ServerURL)
		if err != nil {
			client.err = err
			return client
		}
		client.suParser = parseURL
		client.values = parseURL.Query()
	}
	return client
}

func (client *Oauth2Client) setRedirect() *Oauth2Client {
	if client.err == nil {
		if strings.TrimSpace(client.RedirectURI) == "" {
			return client
		}
		client.values.Set("redirect_uri", client.RedirectURI)
	}
	return client
}

func (client *Oauth2Client) setResponseType() *Oauth2Client {
	if client.err == nil {
		responseType := client.ResponseType
		if strings.TrimSpace(responseType) == "" {
			responseType = "code"
		}
		client.values.Set("response_type", responseType)
	}
	return client
}

func (client *Oauth2Client) setScope() *Oauth2Client {
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
func (client *Oauth2Client) setState() *Oauth2Client {
	if client.err == nil && strings.TrimSpace(client.State) != "" {
		client.values.Set("state", client.State)
	}
	return client
}

func (client *Oauth2Client) setClientID() *Oauth2Client {
	if client.err == nil {
		if strings.TrimSpace(client.Key) == "" {
			client.err = errorx.ClientKeyError
			return client
		}
		client.values.Set("client_id", client.Key)
	}
	return client
}

func (client *Oauth2Client) AuthorizeURL() (string, error) {
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
	c.suParser.RawQuery = client.values.Encode()
	return c.suParser.String(), nil
}

func NewOauth2Client(clientID string, opts ...WithOption) *Oauth2Client {
	var client = &Oauth2Client{Key: clientID}
	for _, opt := range opts {
		opt(client)
	}
	return client
}
