package slackbot

import (
	"regexp"
	"strings"
)

var (
	// example <@UVWXYZ123>
	slackUserRegex = regexp.MustCompile(`^<@(\w+)>`)
)

// Commands issues in a channel (instead of in a DM) will start with the user ID of the bot. Example: '<@UVWXYZ123> describe pod'.
// This should be removed before matching the command to a bot reply.
func TrimUserId(text string) string {
	trimmed := slackUserRegex.ReplaceAllString(text, "bot_user_id")
	return strings.Replace(trimmed, "bot_user_id", "", 1)
}
