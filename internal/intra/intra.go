package intra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/intra/oauth"
)

var (
	baseURL = "https://api.intra.42.fr/"
)

// intraClient to make request to 42's API.
type intraClient struct {
	oauthClient *oauth.Client
}

// NewClient returns a new client made to make request to 42's API.
func NewClient(clientID, clientSecret string, httpClient *http.Client) (Client, error) {
	client, err := oauth.NewClient(baseURL, clientID, clientSecret, httpClient)
	if err != nil {
		return nil, err
	}

	return &intraClient{
		oauthClient: client,
	}, nil
}

// Request makes a request to the API through the configured client.
func (c *intraClient) request(ctx context.Context, method, endpoint string, params oauth.Params, data map[string]interface{}, v interface{}) error {
	resp, err := c.oauthClient.Request(ctx, method, endpoint, params, data)
	if err != nil {
		return err
	}
	if resp.StatusCode == 429 {
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			// If can't parse "Retry-After", return the response
			// atoi's error won't be useful.
			return &HTTPError{Response: resp}
		}
		time.Sleep(time.Second * time.Duration(retryAfter))
		return c.request(ctx, method, endpoint, params, data, v)
	}

	if err := validateResponse(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// GetUserEmail returns a user's email with 42's API.
func (c *intraClient) GetUserEmail(ctx context.Context, login string) (string, error) {
	endpoint := fmt.Sprintf("/v2/users/%s", login)

	user := make(map[string]interface{})

	if err := c.request(ctx, http.MethodGet, endpoint, nil, nil, &user); err != nil {
		return "", err
	}

	return user["email"].(string), nil
}

// GetTeamMembers returns the members of a team with 42's API.
func (c *intraClient) GetTeamMembers(ctx context.Context, teamID int) ([]string, error) {
	endpoint := fmt.Sprintf("/v2/teams/%d/users", teamID)

	users := make([]map[string]interface{}, 0)
	if err := c.request(ctx, http.MethodGet, endpoint, nil, nil, &users); err != nil {
		return nil, err
	}

	logins := make([]string, len(users))
	for i, user := range users {
		logins[i] = user["login"].(string)
	}
	return logins, nil
}
