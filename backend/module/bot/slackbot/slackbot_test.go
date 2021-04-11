package slackbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

)

func TestModule(t *testing.T) {
	service.Registry["clutch.service.bot"] = botmock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	cfg, _ := anypb.New(&slackbotconfigv1.Config{BotToken: "foo", VerifcationToken: "bar"})

	m, err := New(cfg, log, scope)
	assert.NoError(t, err)

	r := moduletest.NewRegisterChecker()
	assert.NoError(t, m.Register(r))
	assert.NoError(t, r.HasAPI("clutch.bot.slackbot.v1.SlackBotAPI"))
	assert.True(t, r.JSONRegistered())
}

func TestVerifySlackRequest(t *testing.T) {
	token := "foo"

	err := verifySlackRequest(token, "bar")
	assert.Error(t, err)

	err = verifySlackRequest(token, "foo")
	assert.NoError(t, err)
}

