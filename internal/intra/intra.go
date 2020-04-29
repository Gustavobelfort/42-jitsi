package intra

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/config"
)

// New creates and returns http client to handle exchanges with 42 API
func New() (Client, error) {
	parsedURL, err := url.Parse("https://api.intra.42.fr")
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(5 * time.Second)
	baseClient := http.Client{
		Timeout: timeout,
	}

	intraClient := &intraClient{
		httpClient: &baseClient,
		baseURL:    parsedURL,
		Token:      "",
	}

	return intraClient, nil
}
func (client *intraClient) getURL(endpoint string) string {
	urlCopy := &url.URL{}
	*urlCopy = *client.baseURL

	urlCopy.Path = path.Join(urlCopy.Path, endpoint)
	return urlCopy.String()
}

func (client *intraClient) request(method, endpoint string, reader io.Reader, v interface{}) error {
	request, err := http.NewRequest(method, client.getURL(endpoint), reader)
	if err != nil {
		return err
	}

	if client.Token != "" {
		formatBearer := fmt.Sprintf("Bearer %s", client.Token)
		request.Header.Set("Authorization", formatBearer)
	}

	resp, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}

	if err := validateResponse(resp); err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}

// GetToken encapsulate getting a valid Token from 42 API
func (client *intraClient) GetToken() error {
	values := url.Values{}
	values.Set("client_id", config.Conf.Intra.AppID)
	values.Set("client_secret", config.Conf.Intra.AppSecret)
	values.Set("grant_type", "client_credentials")
	values.Set("scope", "public projects")

	params := strings.NewReader(values.Encode())

	data := make(map[string]interface{})
	if err := client.request(http.MethodPost, "/oauth/token", params, &data); err != nil {
		return err
	}

	if value, ok := data["access_token"]; ok {
		client.Token = value.(string)
	}
	return nil
}

func (client *intraClient) GetUserEmail(login string) (string, error) {
	endpoint := fmt.Sprintf("/v2/users/%s", login)
	data := make(map[string]interface{})

	if err := client.request(http.MethodGet, endpoint, nil, &data); err != nil {
		return "", err
	}

	if value, ok := data["email"]; ok {
		return value.(string), nil
	}
	return "", nil
}
