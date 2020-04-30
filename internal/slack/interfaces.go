package slack

import (
	"bytes"

	"github.com/gustavobelfort/42-jitsi/internal/db"
)

// PostMessageParameters is the structure used to create the PostMessage request's body.
type PostMessageParameters struct {
	Channel    string   `json:"channel"`
	Workspace  string   `json:"workspace,omitempty"`
	UserEmails []string `json:"user_emails,omitempty"`

	AsUser    bool   `json:"as_user"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`

	Text        string `json:"text"`
	Parse       string `json:"parse,omitempty"`
	Markdown    bool   `json:"mrkdwn"`
	LinkNames   bool   `json:"link_names"`
	UnfurlLinks bool   `json:"unfurl_links"`
	UnfurlMedia bool   `json:"unfurl_media"`

	Attachments map[string]interface{} `json:"attachements,omitempty"`

	buffer *bytes.Buffer
}

// PostMessageOptions are functions that will edit the PostParameters structure used to create the request's body.
type PostMessageOptions func(parameters *PostMessageParameters)

// SlackThatClient will allow you to make prepared request to a slack_that server.
type SlackThat interface {
	PostMessage(options ...PostMessageOptions) (map[string]interface{}, error)
	SendNotification(scaleTeam db.ScaleTeam) error
	GetHealth() (map[string]interface{}, error)
}
