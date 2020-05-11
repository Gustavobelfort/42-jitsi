package slack

import (
	"context"
	"fmt"
	"strings"
)

// SendNotification sends a notification to multiple users containing a link to a meet.jit.si server
func (client *ThatClient) SendNotification(scaleTeamID int, logins []string) error {

	userEmails, err := client.getUserEmails(logins)

	roomName := formatRoomName(scaleTeamID, logins)
	if err != nil {
		return err
	}

	if err := client.postMessage(
		PostMessageUserEmailsOption(userEmails),
		PostMessageLinkOption(roomName),
	); err != nil {
		return err
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
