package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type apiToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`

	CreatedAt int64         `json:"created_at"`
	ExpiresIn time.Duration `json:"expires_in"`
}

func (t *apiToken) token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: t.AccessToken,
		TokenType:   t.TokenType,
		Expiry:      time.Unix(t.CreatedAt, 0).Add(t.ExpiresIn * time.Second),
	}
}

type apiTokenSource struct {
	clientID     string
	clientSecret string

	tokenEndpoint string
	token         *oauth2.Token

	m *sync.Mutex
}

func newTokenSource(baseURL *url.URL, clientID, clientSecret string) oauth2.TokenSource {
	return &apiTokenSource{
		clientID:     clientID,
		clientSecret: clientSecret,

		tokenEndpoint: joinURL(baseURL, "/oauth/token").String(),

		m: new(sync.Mutex),
	}
}

func (ts *apiTokenSource) prepareRequestBody() (io.Reader, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(ts); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (ts *apiTokenSource) requestToken() (*http.Response, error) {
	reader, err := ts.prepareRequestBody()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ts.tokenEndpoint, "application/json", reader)
	if err != nil {
		return nil, err
	}

	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
		}
	}

	return resp, nil
}

func (ts *apiTokenSource) treatResponseBody(body io.ReadCloser) error {
	t := new(apiToken)
	defer body.Close()
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(t); err != nil {
		return err
	}
	ts.token = t.token()
	return nil
}

// Token returns the current token if it is valid.
//
// Otherwise it requests a new token.
func (ts *apiTokenSource) Token() (*oauth2.Token, error) {
	ts.m.Lock()
	defer ts.m.Unlock()
	if ts.token.Valid() {
		return ts.token, nil
	}

	resp, err := ts.requestToken()
	if err != nil {
		return nil, err
	}

	if err := ts.treatResponseBody(resp.Body); err != nil {
		return nil, err
	}
	return ts.token, nil
}

// MarshalJSON returns a valid JSON corresponding to the body used to request a token.
func (ts *apiTokenSource) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     ts.clientID,
		"client_secret": ts.clientSecret,
	}
	return json.Marshal(toMarshal)
}
