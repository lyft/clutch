package slackbot

import (
	"fmt"
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
	spaceRegex       = regexp.MustCompile(`\s+`)
)

// TrimUserId removes the user id of the bot from the command string
// Commands issued in a channel (instead of in a DM) will start with the user ID of the bot. Example: '<@UVWXYZ123> describe pod'
// and should be removed before matching the command to a bot reply.
func TrimUserId(text string) string {
	trimmed := slackUserIdRegex.ReplaceAllString(text, "bot_user_id")
	return strings.Replace(trimmed, "bot_user_id", "", 1)
}

// TrimRedundantSpaces replaces whitespace substrings with a single space character
// this is a sanitization check before matching the command
func TrimRedundantSpaces(text string) string {
	return spaceRegex.ReplaceAllString(text, " ")
}

// DefaultHelp returns the default help message is the command isn't matched
func DefaultHelp() string {
	return fmt.Sprintf("%s/n%s", HelpIntro, HelpDetails)
}

func sendBotReply(client *slack.Client, channel, message string) error {
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	return nil
}
