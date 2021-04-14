package botmock

import (
	"github.com/lyft/clutch/backend/service/bot"
)

type svc struct{}

func (s svc) MatchCommand(string) (reply string) {
	return ""
}

func New() bot.Service {
	return &svc{}
}
