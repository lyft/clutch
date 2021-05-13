package bot

import (
	"fmt"
	"regexp"
	"strings"

	botv1 "github.com/lyft/clutch/backend/api/config/service/bot/v1"
	"github.com/lyft/clutch/backend/service/bot/slack"
)

const (
	HelpIntro = "Hi! I respond to the following commands:\n"
	// TODO: (sperry) this should be auto generated based on what we allow users to do via the bot
	HelpDetails = "help: Print this message\nhello: I will greet you with Hello, World!"
)

var (
	spaceRegex = regexp.MustCompile(`\s+`)
)

func sanitize(botProvider botv1.Bot, text string) (string, error) {
	var sanitized string

	switch botProvider {
	case botv1.Bot_SLACK:
		sanitized = slack.TrimSlackUserId(text)
	default:
		return "", fmt.Errorf("bot type '%s' not implemented", botProvider)
	}

	sanitized = TrimRedundantSpaces(sanitized)
	sanitized = strings.ToLower(sanitized)

	return sanitized, nil
}

func defaultHelp() string {
	// TODO: (sperry) handle pretty formatting, we'll want the majority of the bot response's to be in code blocks
	return fmt.Sprintf("%s\n%s", HelpIntro, HelpDetails)
}

// replaces trailing/leading and duplicate whitespace
func TrimRedundantSpaces(text string) string {
	text = spaceRegex.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
