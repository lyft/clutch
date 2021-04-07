package slackbot

import (
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

// the default bot reply if the command is not matched
const (
	HelpIntro = "Hi! I respond to the following commands:\n"
	// TODO: (sperry) this should be auto generated based on what we allow users to do via the bot
	HelpDetails = "```help: Print this message\nhello: I will greet you with Hello, World!```"
)

var (
	// example: <@UVWXYZ123>
	slackUserIdRegex = regexp.MustCompile(`^<@(\w+)>`)
)

// Commands issued in a channel (instead of in a DM) will start with the user ID of the bot. Example: '<@UVWXYZ123> describe pod'.
// This should be removed before matching the command to a bot reply.
func TrimUserId(text string) string {
	trimmed := slackUserIdRegex.ReplaceAllString(text, "bot_user_id")
	return strings.Replace(trimmed, "bot_user_id", "", 1)
}

func sendBotReply(client *slack.Client, channel, message string) error {
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}
