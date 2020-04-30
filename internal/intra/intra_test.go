package intra

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestIntraClient(t *testing.T) {
	suite.Run(t, new(IntraClientSuite))
}

type IntraClientSuite struct {
	suite.Suite

	client *intraClient
	mock   *ServerMock
}

func (s *IntraClientSuite) SetupSuite() {
	s.Require().Implements((*Client)(nil), &intraClient{})

	s.mock = NewServerMock()
	baseURL = s.mock.Server.URL
	client, err := NewClient("id", "secret", s.mock.Server.Client())
	s.Require().NoError(err)
	s.Require().NotNil(client)

	s.client = client.(*intraClient)
}

func (s *IntraClientSuite) SetupTest() {
	s.mock.ExpectedCalls = []*mock.Call{}
	s.mock.Calls = []mock.Call{}
}

func (s *IntraClientSuite) Test00_GetUserEmail() {
	expectedLogin := "xlogin"
	expectedEmail := "xlogin@student.42campus.org"
	expectedPayload := gin.H{
		"email": expectedEmail,
	}

	s.mock.On("GetUser", expectedLogin).Return(200, expectedPayload, gin.H{}, gin.H{}).Once()
	email, err := s.client.GetUserEmail(context.Background(), expectedLogin)
	s.NoError(err)
	s.Equal(expectedEmail, email)
}

func (s *IntraClientSuite) Test01_GetUserEmail_Error() {
	expectedLogin := "xlogin"

	s.mock.On("GetUser", expectedLogin).Return(404, gin.H{}, gin.H{}).Once()
	email, err := s.client.GetUserEmail(context.Background(), expectedLogin)
	s.Error(err)
	s.Zero(email)
}

func (s *IntraClientSuite) Test02_GetTeamMembers() {
	expectedID := 4242
	expectedLogins := []string{"xlogin", "ylogin", "zlogin"}
	expectedPayload := []gin.H{
		{"login": expectedLogins[0]},
		{"login": expectedLogins[1]},
		{"login": expectedLogins[2]},
	}

	s.mock.On("GetTeamUsers", strconv.Itoa(expectedID)).Return(200, expectedPayload, gin.H{}).Once()
	logins, err := s.client.GetTeamMembers(context.Background(), expectedID)
	s.NoError(err)
	s.Equal(expectedLogins, logins)

}

func (s *IntraClientSuite) Test03_GetTeamMembers_Error() {
	expectedID := 4242

	s.mock.On("GetTeamUsers", strconv.Itoa(expectedID)).Return(404, gin.H{}, gin.H{}).Once()
	logins, err := s.client.GetTeamMembers(context.Background(), expectedID)
	s.Error(err)
	s.Zero(logins)
}

func (s *IntraClientSuite) Test04_RateLimitHandling() {
	expectedLogin := "xlogin"
	expectedEmail := "xlogin@student.42campus.org"
	expectedPayload := gin.H{
		"email": expectedEmail,
	}

	s.mock.On("GetUser", expectedLogin).Return(429, gin.H{}, gin.H{"Retry-After": "1"}).Twice()
	s.mock.On("GetUser", expectedLogin).Return(200, expectedPayload, gin.H{}, gin.H{}).Once()
	email, err := s.client.GetUserEmail(context.Background(), expectedLogin)
	s.NoError(err)
	s.Equal(expectedEmail, email)
}

func (s *IntraClientSuite) Test05_RateLimitHandling_Error() {
	expectedLogin := "xlogin"

	s.mock.On("GetUser", expectedLogin).Return(429, gin.H{}, gin.H{}).Once()
	email, err := s.client.GetUserEmail(context.Background(), expectedLogin)
	s.Error(err)
	s.Zero(email)
}

func (s *IntraClientSuite) Test06_ContextCanceled() {
	expectedLogin := "xlogin"

	expectedContext, cancel := context.WithCancel(context.Background())
	cancel()

	email, err := s.client.GetUserEmail(expectedContext, expectedLogin)
	s.Error(err)
	s.True(errors.Is(err, context.Canceled))
	s.Zero(email)
}

func (s *IntraClientSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *IntraClientSuite) TearDownSuite() {
	s.mock.Server.Close()
	baseURL = "https://api.intra.42.fr/"
}

//func TestGetUserEmail(t *testing.T) {
//
//	config.Initiate()
//
//	client, err := New()
//	client.GetToken()
//
//	if err != nil {
//		t.Errorf("Failed to initialize the client with error: %d", err)
//	}
//
//	email, err := client.GetUserEmail("gbelfort")
//	if err != nil {
//		t.Errorf("Failed to Get the user email with error: %s", err.Error())
//	}
//
//	assert.Equal(t, email, "gbelfort@student.42.us.org")
//}
