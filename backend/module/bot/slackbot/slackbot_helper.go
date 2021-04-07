package slackbot

import (
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

var (
	// example <@UVWXYZ123>
	slackUserRegex = regexp.MustCompile(`^<@(\w+)>`)
)

// Commands issued in a channel (instead of in a DM) will start with the user ID of the bot. Example: '<@UVWXYZ123> describe pod'.
// This should be removed before matching the command to a bot reply.
func TrimUserId(text string) string {
	trimmed := slackUserRegex.ReplaceAllString(text, "bot_user_id")
	return strings.Replace(trimmed, "bot_user_id", "", 1)
}

func sendBotReply(client *slack.Client, channel, message string) error {
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}
