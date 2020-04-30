package slack

import (
	"fmt"
	"strings"

	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
)

func (client *SlackThatClient) SendNotification(scaleTeam db.ScaleTeam) error {

	scaleTeamID := scaleTeam.GetID()
	scaleTeamUsers, err := scaleTeam.GetUsers()
	if err != nil {
		return err
	}

	var logins []string
	for _, user := range scaleTeamUsers {
		logins = append(logins, user.GetLogin())
	}

	userEmails, err := intra.Client.GetUsersEmails(logins)
	if err != nil {
		return err
	}

	roomName := formatRoomName(scaleTeamID, logins)

	client.PostMessage(
		PostMessageUserEmailsOption(userEmails),
		PostMessageLinkOption(roomName),
	)

	return nil
}

func formatRoomName(scaleTeamID int, logins []string) string {
	jitsiServer := "https://meet.jit.si/"
	usernames := strings.Join(logins, "-")

	return fmt.Sprintf("%s%d-", jitsiServer, scaleTeamID) + usernames
}
