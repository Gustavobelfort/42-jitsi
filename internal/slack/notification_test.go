package slack

import (
	"net/http"
	"testing"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
	"github.com/stretchr/testify/suite"
)

func TestSlackClient(t *testing.T) {
	suite.Run(t, new(SlackClientSuite))
}

type SlackClientSuite struct {
	suite.Suite

	client SlackThat
}

func (s *SlackClientSuite) SetupSuite() {
	s.Require().Implements((*SlackThat)(nil), &ThatClient{})
	config.Initiate()
	intra, err := intra.NewClient(config.Conf.Intra.AppID, config.Conf.Intra.AppSecret, http.DefaultClient)
	s.client, err = New(intra)
	s.Require().NoError(err)
}

func (s *SlackClientSuite) Test00_SendNotification() {
	logins := []string{"gus"}
	err := s.client.SendNotification(1, logins)
	s.Require().NoError(err)
}
