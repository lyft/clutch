package botmock

import (
	"context"

	"github.com/lyft/clutch/backend/service/bot"
)

type svc struct{}

func (s svc) MatchCommand(ctx context.Context, command string) (reply string) {
	return ""
}

func New() bot.Service {
	return &svc{}
}
