package slack

import (
	"regexp"
	"strings"
)

var (
	// example: <@UVWXYZ123>
	slackUserIdRegex = regexp.MustCompile(`^<@(\w+)>`)
)

// TrimSlackUserId removes the user id of the bot from the command string. Commands issued in a channel (instead of in a DM) will start
// with the user ID of the bot. Example: '<@UVWXYZ123> describe pod'
func TrimSlackUserId(text string) string {
	trimmed := slackUserIdRegex.ReplaceAllString(text, "bot_user_id")
	// removes an extra whitespace so that there isn't a leading whitespace in the beginning of the new string
	return strings.Replace(trimmed, "bot_user_id ", "", 1)
}
