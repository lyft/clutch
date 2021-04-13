package bot

import (
	"fmt"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	botv1 "github.com/lyft/clutch/backend/api/config/service/bot/v1"
	"github.com/lyft/clutch/backend/service"
)

const (
	Name = "clutch.service.bot"
)

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &botv1.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	t, err := getBotProvider(config)
	if err != nil {
		return nil, err
	}

	s := &svc{
		logger:      logger,
		scope:       scope,
		botProvider: t,
	}
	return s, nil
}

func getBotProvider(cfg *botv1.Config) (botv1.Bot, error) {
	switch cfg.BotProvider {
	case botv1.Bot_SLACK:
		return botv1.Bot_SLACK, nil
	default:
		return 0, fmt.Errorf("bot type '%v' not implemented", cfg.BotProvider)
	}
}

type svc struct {
	logger      *zap.Logger
	scope       tally.Scope
	botProvider botv1.Bot
}

type Service interface {
	MatchCommand(command string) string
}

// MatchCommand matches a command to a Clutch API call
func (s *svc) MatchCommand(command string) string {
	command, err := sanitize(s.botProvider, command)
	if err != nil {
		// return the error as the bot's reply instead of failing silently
		return err.Error()
	}

	// TODO: (sperry) a placeholder for now, a followup PR will handle matching the command to a Clutch API call
	switch command {
	case "hello":
		return "Hello, World!"
	default:
		return defaultHelp()
	}
}
