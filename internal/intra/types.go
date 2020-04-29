package intra

import (
	"net/http"
	"net/url"
)

type Client interface {
	GetToken() error
	GetUserEmail(login string) (string, error)
}

type intraClient struct {
	Token      string
	baseURL    *url.URL
	httpClient *http.Client
}
