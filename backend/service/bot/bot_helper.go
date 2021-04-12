package bot

import (
	"regexp"
	"strings"
)

var (
	spaceRegex = regexp.MustCompile(`\s+`)
)

// TrimRedundantSpaces replaces trailing/leading and duplicate whitespace. A sanitization check before matching the command
func TrimRedundantSpaces(text string) string {
	text = spaceRegex.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
