package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// Client is a OAuth2 client made simplify authenticated requests to an OAuth2 API.
type Client struct {
	client *http.Client
	ctx    context.Context

	baseURL *url.URL
}

// NewClient returns a new OAuth2 client configured with the given base url, client id and client secret.
func NewClient(baseURL, clientID, clientSecret string, httpClient *http.Client) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	return &Client{
		client: oauth2.NewClient(ctx, newTokenSource(parsedURL, clientID, clientSecret)),

		baseURL: parsedURL,
	}, nil
}

func (*Client) prepareBody(method string, data map[string]interface{}) (io.Reader, error) {
	if method == http.MethodGet || data == nil {
		return nil, nil
	}
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (c *Client) prepareURL(endpoint string, params Params) string {
	requestURL := joinURL(c.baseURL, endpoint)
	if params != nil {
		requestURL.RawQuery = params.Encode()
	}
	return requestURL.String()
}

func (c *Client) prepareRequest(ctx context.Context, method, endpoint string, params Params, data map[string]interface{}) (*http.Request, error) {
	body, err := c.prepareBody(method, data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.prepareURL(endpoint, params), body)
	if err != nil {
		return nil, err
	}
	return req.WithContext(ctx), nil
}

// Request makes a request to the API through the configured client.
func (c *Client) Request(ctx context.Context, method, endpoint string, params Params, data map[string]interface{}) (*http.Response, error) {
	request, err := c.prepareRequest(ctx, method, endpoint, params, data)
	if err != nil {
		return nil, err
	}

	return c.client.Do(request)
}
