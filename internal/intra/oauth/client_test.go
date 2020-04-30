package oauth

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type ClientSuite struct {
	suite.Suite

	client *Client
	mock   *ServerMock

	clientID     string
	clientSecret string
	token        string
}

func (s *ClientSuite) SetupSuite() {
	var err error

	s.mock = NewServerMock()

	s.clientID = "id"
	s.clientSecret = "secret"
	s.token = "Bearer token"

	s.client, err = NewClient(s.mock.Server.URL, s.clientID, s.clientSecret, s.mock.Server.Client())
	s.Require().NoError(err)
	s.Require().NotNil(s.client)
}

func (s *ClientSuite) expectOauth() {
	expectedData := gin.H{
		"grant_type":    "client_credentials",
		"client_id":     s.clientID,
		"client_secret": s.clientSecret,
	}
	expectedResponse := gin.H{
		"access_token": "token",
		"token_type":   "Bearer",
		"created_at":   time.Now().Unix(),
		"expires_in":   7200,
	}

	s.mock.On("/oauth/token", http.MethodPost, url.Values{}, expectedData, "").Return(200, expectedResponse, gin.H{}).Once()
}

func (s *ClientSuite) SetupTest() {
	s.mock.ExpectedCalls = []*mock.Call{}
	s.mock.Calls = []mock.Call{}
}

func (s *ClientSuite) Test00_Get() {
	s.expectOauth() // First Test should try to authenticate

	s.mock.On("/get", http.MethodGet, url.Values{}, gin.H(nil), s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodGet, "/get", nil, nil)

	s.mock.On("/get_ignore_body", http.MethodGet, url.Values{}, gin.H(nil), s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodGet, "/get_ignore_body", nil, gin.H{"ignore": "me"})

	expectedQuery := url.Values{
		"key":  {"value"},
		"key2": {"value2"},
	}

	s.mock.On("/get_query", http.MethodGet, expectedQuery, gin.H(nil), s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodGet, "/get_query", Params(expectedQuery), nil)
}

func (s *ClientSuite) Test01_Post() {
	s.mock.On("/post", http.MethodPost, url.Values{}, gin.H(nil), s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodPost, "/post", nil, nil)

	expectedBody := gin.H{
		"dont": "ignore_me",
	}

	s.mock.On("/post_body", http.MethodPost, url.Values{}, expectedBody, s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodPost, "/post_body", nil, expectedBody)

	expectedQuery := url.Values{
		"key":  {"value"},
		"key2": {"value2"},
	}

	s.mock.On("/post_query", http.MethodPost, expectedQuery, gin.H(nil), s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodPost, "/post_query", Params(expectedQuery), nil)

	s.mock.On("/post_query_and_body", http.MethodPost, expectedQuery, expectedBody, s.token).Return(200, gin.H{}, gin.H{}).Once()
	s.client.Request(context.Background(), http.MethodPost, "/post_query_and_body", Params(expectedQuery), expectedBody)
}

func (s *ClientSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *ClientSuite) TearDownSuite() {
	s.mock.Server.Close()
}
