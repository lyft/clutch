package slackbot

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	slackbotv1 "github.com/lyft/clutch/backend/api/bot/slackbot/v1"
	slackbotconfigv1 "github.com/lyft/clutch/backend/api/config/module/bot/slackbot/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/bot"
)

const (
	Name      = "clutch.module.bot.slackbot"
	dmMessage = "im"
)

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &slackbotconfigv1.Config{}

	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	svc, ok := service.Registry["clutch.service.bot"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	bot, ok := svc.(bot.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	m := &mod{
		slack:             slack.New(config.BotToken),
		verificationToken: config.VerificationToken,
		bot:               bot,
		logger:            logger,
		scope:             scope,
	}

	return m, nil
}

type mod struct {
	slack             *slack.Client
	verificationToken string
	bot               bot.Service
	logger            *zap.Logger
	scope             tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	slackbotv1.RegisterSlackBotAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(slackbotv1.RegisterSlackBotAPIHandler)
}

// One request, one event. https://api.slack.com/apis/connections/events-api#the-events-api__receiving-events
func (m *mod) Event(ctx context.Context, req *slackbotv1.EventRequest) (*slackbotv1.EventResponse, error) {
	err := verifySlackRequest(m.verificationToken, req.Token)
	if err != nil {
		m.logger.Error("verify slack request error", log.ErrorField(err))
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	m.logger.Info("received event from Slack API")

	if req.Type == slackevents.URLVerification {
		return m.handleURLVerificationEvent(req.Challenge)
	}

	// at this point in the flow, we can ack the Slack API with a 2xx response and process the event seperately
	// TODO: (sperry) use a worker goroutine / save event to DB
	go m.handleEvent(req)

	return &slackbotv1.EventResponse{}, nil
}

// TODO: (sperry) Slack currently still supports this verification type but they plan to deprecate this inf the future &
// we want to check the signing secret instead so this is a short-term solution (need to change API design first)
func verifySlackRequest(token, requestToken string) error {
	t := slackevents.TokenComparator{VerificationToken: token}
	ok := t.Verify(requestToken)
	if !ok {
		return errors.New("invalid verification token")
	}
	return nil
}

// event type recieved the first time when configuring Slack API to send requests to our endpoint, https://api.slack.com/events/url_verification
func (m *mod) handleURLVerificationEvent(challenge string) (*slackbotv1.EventResponse, error) {
	// this should never happen
	if challenge == "" {
		msg := "slack API provided an empty challenge string"
		m.logger.Error(msg)
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return &slackbotv1.EventResponse{Challenge: challenge}, nil
}

// this outer event is 1 of 3 types: url_verification, event_callback, and app_rate_limited
// this should be called via `go` in order to avoid blocking main execution
func (m *mod) handleEvent(req *slackbotv1.EventRequest) {
	switch req.Type {
	case slackevents.CallbackEvent:
		err := m.handleCallBackEvent(req.Event)
		if err != nil {
			m.logger.Error("handle callback event error", log.ErrorField(err))
		}
	case slackevents.AppRateLimited:
		// received if endpoint got more than 30,000 request events in 60 minutes, https://api.slack.com/apis/connections/events-api#the-events-api__responding-to-events__rate-limiting
		m.logger.Error("app's event subscriptions are being rate limited")
	default:
		m.logger.Error("received unexpected event type", zap.String("event type", req.Type))
	}
}

// for the inner event, we currently support 2 types: app_mention (messages that mention the bot directly) and message (specifically DMs with the bot)
// full list of Slack event types: https://api.slack.com/events
func (m *mod) handleCallBackEvent(event *slackbotv1.Event) error {
	switch event.Type {
	case slackevents.AppMention:
		return m.handleAppMentionEvent(event)
	case slackevents.Message:
		return m.handleMessageEvent(event)
	default:
		return fmt.Errorf("received unsuported event type: %s", event.Type)
	}
}

// event type for messages that mention the bot directly
func (m *mod) handleAppMentionEvent(event *slackbotv1.Event) error {
	match := m.bot.MatchCommand(event.Text)
	err := sendBotReply(m.slack, event.Channel, match)
	if err != nil {
		return err
	}

	return nil
}

// event type for DMs with the bot
func (m *mod) handleMessageEvent(event *slackbotv1.Event) error {
	if event.ChannelType != dmMessage {
		return fmt.Errorf("unexpected channel type: %s", event.ChannelType)
	}
	// it's standard behavior of the Slack Events API that the bot will receive all message events, including from it's own posts
	// so we need to filter out these events by checking that the BotId field is empty; otherwise by reacting to them we will create an infinte loop
	if event.BotId == "" {
		match := m.bot.MatchCommand(event.Text)
		err := sendBotReply(m.slack, event.Channel, match)
		if err != nil {
			return err
		}
	}
	return nil
}
