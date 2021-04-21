package slackbot

import (
	"context"

	"github.com/slack-go/slack"
)

func sendBotReply(ctx context.Context, client *slack.Client, channel, message string) error {
	_, _, err := client.PostMessageContext(ctx, channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}
