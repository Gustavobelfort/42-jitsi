package slack

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/config"
)

type SlackThatClient struct {
	HttpClient *http.Client
	BaseURL    *url.URL
}

func (client *SlackThatClient) getURL(endpoint string) string {
	urlCopy := &url.URL{}
	*urlCopy = *client.BaseURL

	urlCopy.Path = path.Join(urlCopy.Path, endpoint)
	return urlCopy.String()
}

func (client *SlackThatClient) request(method string, endpoint string, reader io.Reader, v interface{}) error {
	request, err := http.NewRequest(method, client.getURL(endpoint), reader)
	if err != nil {
		return err
	}
	resp, err := client.HttpClient.Do(request)
	if err != nil {
		return err
	}

	if err := validateResponse(resp); err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}

// PostMessage makes a POST request to the slack_that API to send a message to multiple users in
// a workspace.
func (client *SlackThatClient) PostMessage(options ...PostMessageOptions) (map[string]interface{}, error) {
	params := defaultNotificationParameters()
	for _, opt := range options {
		opt(params)
	}

	data := make(map[string]interface{})
	if err := client.request(http.MethodPost, "/", params, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// GetHealth makes a GET request to the slack_that API's health endpoint.
func (client *SlackThatClient) GetHealth() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if err := client.request(http.MethodGet, "/health", nil, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// Initiate a SlackThat client ready to make requests to the base_url passed.
func Initiate() error {

	parsedURL, err := url.Parse(config.Conf.SlackThat.URL)
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	baseClient := http.Client{
		Timeout: timeout,
	}

	Client = &SlackThatClient{
		HttpClient: &baseClient,
		BaseURL:    parsedURL,
	}

	return nil
}

var Client *SlackThatClient
