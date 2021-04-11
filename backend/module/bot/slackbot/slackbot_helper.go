package slackbot

import (
	"github.com/slack-go/slack"
)

func sendBotReply(client *slack.Client, channel, message string) error {
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}
