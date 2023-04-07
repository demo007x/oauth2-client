package utils

import (
	"context"
	"io"
	"net/http"
	nurl "net/url"
)

func DoRequest(url, method string, header map[string]string) (*http.Response, error) {
	req, err := buildRequest(context.Background(), method, url, header)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// buildRequest build http request params
func buildRequest(ctx context.Context, method, url string, header map[string]string) (*http.Request, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	if len(header) != 0 {
		for key, val := range header {
			req.Header.Set(key, val)
		}
	}

	return req, nil
}
