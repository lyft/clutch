package slackbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	slackbotconfigv1 "github.com/lyft/clutch/backend/api/config/module/bot/slackbot/v1"
	"github.com/lyft/clutch/backend/mock/service/botmock"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.bot"] = botmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	cfg, _ := anypb.New(&slackbotconfigv1.Config{BotToken: "foo", VerificationToken: "bar"})

	m, err := New(cfg, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.bot.slackbot.v1.SlackBotAPI"))
	assert.True(t, r.JSONRegistered())
}

func modMock(t *testing.T) *mod {
	log := zaptest.NewLogger(t)
	return &mod{
		logger: log,
	}
}

func TestVerifySlackRequest(t *testing.T) {
	token := "foo"

	err := verifySlackRequest(token, "bar")
	assert.Error(t, err)

	err = verifySlackRequest(token, "foo")
	assert.NoError(t, err)
}

func TestHandleURLVerificationEvent(t *testing.T) {
	challenge := "foo"
	m := modMock(t)

	resp, err := m.handleURLVerificationEvent(challenge)
	assert.NoError(t, err)
	assert.Equal(t, resp.Challenge, challenge)

	resp, err = m.handleURLVerificationEvent("")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
