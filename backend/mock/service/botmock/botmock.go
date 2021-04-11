package botmock

import (
	"github.com/lyft/clutch/backend/service/bot"
)

type svc struct{}

func (s svc) MatchSlackCommand(string) string {
	return ""
}

func (s svc) MatchCommand(string) (string, bool) {
	return "", true
}

func New() bot.Service {
	return &svc{}
}
