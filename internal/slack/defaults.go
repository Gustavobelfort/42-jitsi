package slack

import "github.com/gustavobelfort/42-jitsi/internal/config"

func defaultPostMessageParameters() *PostMessageParameters {
	return &PostMessageParameters{
		Text:      "This is the link for your evaluation that will take place in 15 minutes.",
		Username:  config.Conf.SlackThat.Username,
		Workspace: config.Conf.SlackThat.Workspace,
		Attachments: []Attachment{
			{
				Title:     "42 Evaluation",
				Pretext:   "Make sure to arrive on time and follow the remote correction guidelines !",
				TitleLink: "https://meet.jit.si/correction",
				Color:     "#36a64f",
			},
		},
	}

}
