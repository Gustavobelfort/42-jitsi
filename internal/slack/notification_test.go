package slack

import (
	"context"
	"testing"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestSlackClient(t *testing.T) {
	suite.Run(t, new(SlackClientSuite))
}

type SlackClientSuite struct {
	suite.Suite

	mock   *ServerMock
	client SlackThat
}

type IntraMock struct {
	mock.Mock
}

func (m IntraMock) GetTeamMembers(ctx context.Context, teamID int) ([]string, error) {
	return nil, nil
}

func (m IntraMock) GetUserEmail(ctx context.Context, login string) (string, error) {
	return "nil", nil
}

func (s *SlackClientSuite) SetupSuite() {
	s.Require().Implements((*SlackThat)(nil), &ThatClient{})

	config.Conf.SlackThat.Workspace = "testWorkspace"
	s.mock = NewServerMock()

	client, err := New(IntraMock{}, s.mock.Server.URL)
	s.client = client.(*ThatClient)
	s.Require().NoError(err)
	s.Require().NotNil(s.client)

}

func (s *SlackClientSuite) Test00_SendNotification() {

	expectedLogin := []string{"xlogin"}
	expectedScaleTeamID := 1

	s.mock.On("SendNotification", expectedLogin, expectedScaleTeamID).Return(nil).Once()
	err := s.client.SendNotification(1, expectedLogin)
	s.NoError(err)
}
