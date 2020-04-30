package slack_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/slack"
)

func TestSendNotification(t *testing.T) {

	mockScaleTeam := &db.ScaleTeamModel{
		ID:       123,
		BeginAt:  time.Now(),
		Notified: false,
		Users: []db.UserModel{
			db.UserModel{
				ID:          1,
				ScaleTeamID: 123,
				Login:       "gbelfort",
			},
		},
	}

	parsedURL, _ := url.Parse("http://localhost:8080")
	baseClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	mockSlackClient := slack.SlackThatClient{
		HttpClient: &baseClient,
		BaseURL:    parsedURL,
	}

	mockSlackClient.SendNotification(mockScaleTeam)
}
