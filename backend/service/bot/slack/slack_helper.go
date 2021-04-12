package slack

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	HelpIntro = "Hi! I respond to the following commands:\n"
	// TODO: (sperry) this should be auto generated based on what we allow users to do via the bot
	HelpDetails = "```help: Print this message\nhello: I will greet you with Hello, World!```"
)

var (
	// example: <@UVWXYZ123>
	slackUserIdRegex = regexp.MustCompile(`^<@(\w+)>`)
)

// TrimSlackUserId removes the user id of the bot from the command string. Commands issued in a channel (instead of in a DM) will start
// with the user ID of the bot and should be removed before matching the command. Example: '<@UVWXYZ123> describe pod'
func TrimSlackUserId(text string) string {
	trimmed := slackUserIdRegex.ReplaceAllString(text, "bot_user_id")
	// removes an extra whitespace so that there isn't a leading whitespace in the beginning of the new string
	return strings.Replace(trimmed, "bot_user_id ", "", 1)
}

// DefaultHelp returns the default help message is the command isn't matched
func DefaultHelp() string {
	return fmt.Sprintf("%s\n%s", HelpIntro, HelpDetails)
}
