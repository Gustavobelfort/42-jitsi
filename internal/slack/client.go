package slack

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
)

type ThatClient struct {
	HTTPClient *http.Client
	Intra      intra.Client
	BaseURL    *url.URL
}

func (client *ThatClient) getURL(endpoint string) string {
	urlCopy := &url.URL{}
	*urlCopy = *client.BaseURL

	urlCopy.Path = path.Join(urlCopy.Path, endpoint)
	return urlCopy.String()
}

func (client *ThatClient) request(method string, endpoint string, reader io.Reader) error {
	request, err := http.NewRequest(method, client.getURL(endpoint), reader)
	if err != nil {
		return err
	}
	resp, err := client.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("unable to send request to the slack that client")
	}

	return nil
}

func (client *ThatClient) postMessage(options ...PostMessageOptions) error {
	params := defaultPostMessageParameters()
	for _, opt := range options {
		opt(params)
	}

	if err := client.request(http.MethodPost, "/", params); err != nil {
		return err
	}

	return nil
}

// GetHealth makes a GET request to the slack_that API's health endpoint.
func (client *ThatClient) GetHealth() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if err := client.request(http.MethodGet, "/health", nil); err != nil {
		return nil, err
	}
	return data, nil
}

// Initiate a SlackThat client ready to make requests to the base_url passed.
func New(intra intra.Client) (SlackThat, error) {

	parsedURL, err := url.Parse(config.Conf.SlackThat.URL)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(5 * time.Second)
	baseClient := http.Client{
		Timeout: timeout,
	}

	client := &ThatClient{
		BaseURL:    parsedURL,
		HTTPClient: &baseClient,
		Intra:      intra,
	}

	return client, nil
}
