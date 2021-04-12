package bot

import (
	"strings"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/bot/slack"
)

const Name = "clutch.service.bot"

func New(_ *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	s := &svc{
		logger: logger,
		scope:  scope,
	}
	return s, nil
}

type svc struct {
	logger *zap.Logger
	scope  tally.Scope
}

type Service interface {
	// Functions that perform provider-specific sanitization/formatting on a given command string
	// They all pass the command to the generic MatchCommand function
	MatchSlackCommand(command string) string

	// MatchCommand matches a command to a Clutch API call
	MatchCommand(command string) (string, bool)
}

// MatchSlackCommand takes in a user's slack message command and returns the bot's reply
func (s *svc) MatchSlackCommand(command string) string {
	command = slack.TrimSlackUserId(command)
	match, ok := s.MatchCommand(command)
	if !ok {
		return slack.DefaultHelp()
	}
	return match
}

// MatchCommand matches the given pattern in a command to a Clutch API call, returning the match if found. Otherwise ok is false.
func (s *svc) MatchCommand(command string) (string, bool) {
	command = strings.ToLower(command)
	command = TrimRedundantSpaces(command)
	switch command {
	case "hello":
		return "Hello, World!", true
	default:
		return "", false
	}
}
