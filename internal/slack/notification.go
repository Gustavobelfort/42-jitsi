package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

// SendNotification sends a notification to multiple users containing a link to a meet.jit.si server
func (client *ThatClient) SendNotification(scaleTeamID int, logins []string) error {

	logrus.WithField("scale_team_id", scaleTeamID).Info("getting scale team users' emails")
	userEmails, err := client.getUserEmails(logins)
	roomName := formatRoomName(scaleTeamID, logins)
	ctxfields := logrus.Fields{
		"scale_team_id": scaleTeamID,
		"room_name":     roomName,
	}
	if err != nil {
		return logging.WithLog(err, logrus.ErrorLevel, ctxfields)
	}

	logrus.WithFields(ctxfields).Info("posting message to slack_that")
	if err := client.postMessage(
		PostMessageUserEmailsOption(userEmails),
		PostMessageLinkOption(roomName),
	); err != nil {
		return logging.WithLog(err, logrus.ErrorLevel, ctxfields)
	}

	return nil
}

func (client *ThatClient) getUserEmails(logins []string) ([]string, error) {
	var userEmails []string
	for _, login := range logins {
		email, err := client.Intra.GetUserEmail(context.Background(), login)
		if err != nil {
			return nil, err
		}
		userEmails = append(userEmails, email)
	}
	return userEmails, nil
}

func formatRoomName(scaleTeamID int, logins []string) string {
	jitsiServer := "https://meet.jit.si/"
	usernames := strings.Join(logins, "-")

	return fmt.Sprintf("%s%d-", jitsiServer, scaleTeamID) + usernames
}
