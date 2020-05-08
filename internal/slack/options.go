package slack

// PostMessageUserEmailsOption changes the users emails to whom post the message.
func PostMessageUserEmailsOption(userEmails []string) PostMessageOptions {
	return func(parameters *PostMessageParameters) {
		parameters.UserEmails = userEmails
	}
}

// PostMessageAttachmentsOption changes the users emails to whom post the message.
func PostMessageLinkOption(link string) PostMessageOptions {
	return func(parameters *PostMessageParameters) {
		parameters.Attachments[0].TitleLink = link
	}
}
