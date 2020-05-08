package slack

import "bytes"

// PostMessageParameters is the structure used to create the PostMessage request's body.
type PostMessageParameters struct {
	Channel    string   `json:"channel"`
	Workspace  string   `json:"workspace,omitempty"`
	UserEmails []string `json:"user_emails,omitempty"`

	AsUser    bool   `json:"as_user"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`

	Text string `json:"text"`

	Attachments []Attachment `json:"attachments,omitempty"`

	buffer *bytes.Buffer
}

type Attachment struct {
	Color     string `json:"color,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
}

// PostMessageOptions are functions that will edit the PostParameters structure used to create the request's body.
type PostMessageOptions func(parameters *PostMessageParameters)

// SlackThat will allow you to make prepared request to a slack_that server.
type SlackThat interface {
	SendNotification(scaleTeamID int, logins []string) error
	GetHealth() (map[string]interface{}, error)
}
