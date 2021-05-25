package slackbot

import (
	"context"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/types/known/anypb"

	slackbotv1 "github.com/lyft/clutch/backend/api/bot/slackbot/v1"
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
	svc := botmock.New()
	return &mod{logger: log, bot: svc}
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

func TestHandleEvent(t *testing.T) {
	// testing log count
	testCases := []struct {
		req   *slackbotv1.EventRequest
		count int
	}{
		{
			req: &slackbotv1.EventRequest{
				Type:  "event_callback",
				Event: &slackbotv1.Event{Type: "foo"},
			},
			count: 1,
		},
		{req: &slackbotv1.EventRequest{Type: "app_rate_limited"}, count: 1},
		{req: &slackbotv1.EventRequest{Type: "foo"}, count: 1},
	}

	for _, test := range testCases {
		core, recorded := observer.New(zapcore.DebugLevel)
		log := zap.New(core)
		m := &mod{logger: log}

		m.handleEvent(context.Background(), test.req)
		assert.Equal(t, test.count, recorded.Len())
	}
}

func TestHandleCallBackEvent(t *testing.T) {
	m := modMock(t)

	// invalid event type
	event := &slackbotv1.Event{Type: "foo"}
	err := m.handleCallBackEvent(context.Background(), event)
	assert.Error(t, err)
}

func TestAppMentionEvent(t *testing.T) {
	m := modMock(t)

	// setting up slack mocks
	slacktest.NewTestServer()
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	m.slack = client

	event := &slackbotv1.Event{Type: "app_mention"}
	err := m.handleAppMentionEvent(context.Background(), event)
	assert.NoError(t, err)
}

func TestHandleMessageEvent(t *testing.T) {
	m := modMock(t)

	// invalid event type
	event := &slackbotv1.Event{Type: "foo"}
	err := m.handleMessageEvent(context.Background(), event)
	assert.Error(t, err)

	// setting up slack mocks
	slacktest.NewTestServer()
	s := slacktest.NewTestServer()
	go s.Start()
	defer s.Stop()

	client := slack.New("ABCDEFG", slack.OptionAPIURL(s.GetAPIURL()))
	m.slack = client

	event = &slackbotv1.Event{Type: "message", ChannelType: "im"}
	err = m.handleMessageEvent(context.Background(), event)
	assert.NoError(t, err)
}
