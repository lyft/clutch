package slackbot

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	slackbotv1 "github.com/lyft/clutch/backend/api/bot/slackbot/v1"
	slackbotconfigv1 "github.com/lyft/clutch/backend/api/config/module/bot/slackbot/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name = "clutch.module.bot.slackbot"
	// request headers used to verify a request came from Slack
	hSignature = "x-Slack-Signature"
	hTimestamp = "x-Slack-Request-Timestamp"
)

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &slackbotconfigv1.Config{}
	if cfg != nil {
		if err := cfg.UnmarshalTo(config); err != nil {
			return nil, err
		}
	}

	m := &mod{
		slack:         slack.New(config.BotToken),
		logger:        logger,
		scope:         scope,
		signingSecret: config.SigningSecret,
	}

	return m, nil
}

type mod struct {
	slack         *slack.Client
	signingSecret string
	logger        *zap.Logger
	scope         tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	slackbotv1.RegisterSlackBotAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(slackbotv1.RegisterSlackBotAPIHandler)
}

// One request, one event. https://api.slack.com/apis/connections/events-api#the-events-api__receiving-events
func (m *mod) Event(ctx context.Context, req *slackbotv1.EventRequest) (*slackbotv1.EventResponse, error) {
	err := verifySlackRequest(ctx, req, m.signingSecret)
	if err != nil {
		m.logger.Error("verify slack request error", log.ErrorField(err))
		return nil, status.New(codes.Unauthenticated, err.Error()).Err()
	}

	m.logger.Info("received event from Slack API")

	if req.Type == slackevents.URLVerification {
		return m.handleURLVerificationEvent(req.Challenge)
	}

	return &slackbotv1.EventResponse{}, nil
}

// retrieves the signature and timestamp from the headers which will be used to verify the request was from Slack
func getSlackRequestHeaders(ctx context.Context) (http.Header, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no headers present on request")
	}

	signature := md.Get(hSignature)
	if len(signature) == 0 {
		return nil, fmt.Errorf("%s header not present on request", hSignature)
	}

	timestamp := md.Get(hTimestamp)
	if len(timestamp) == 0 {
		return nil, fmt.Errorf("%s header not present on request", hTimestamp)
	}

	return http.Header{
		"X-Slack-Signature":         signature,
		"X-Slack-Request-Timestamp": timestamp,
	}, nil
}

// verifies the request came from Slack
// https://api.slack.com/authentication/verifying-requests-from-slack#verifying-requests-from-slack-using-signing-secrets__a-recipe-for-security__step-by-step-walk-through-for-validating-a-request
func verifySlackRequest(ctx context.Context, req *slackbotv1.EventRequest, signingSecret string) error {
	headers, err := getSlackRequestHeaders(ctx)
	if err != nil {
		return err
	}

	// returns a SecretsVerifier object in exchange for an http.Header object and signing secret
	sv, err := slack.NewSecretsVerifier(headers, signingSecret)
	if err != nil {
		return err
	}

	body, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	// the request body is used to form a basestring
	_, err = sv.Write(body)
	if err != nil {
		return err
	}

	// compares the signature sent from Slack with the actual computed hash to judge validity
	err = sv.Ensure()
	if err != nil {
		return err
	}

	// we're good! request came from Slack
	return nil
}

// event type recieved the first time when configuring Slack API to send requests to our endpoint, https://api.slack.com/events/url_verification
// the Slack API will send us the challenge string in the request and we must respond back with the same challenge
func (m *mod) handleURLVerificationEvent(challenge string) (*slackbotv1.EventResponse, error) {
	// this should never happen
	if challenge == "" {
		msg := "slack API provided an empty challenge string"
		m.logger.Error(msg)
		return nil, status.New(codes.InvalidArgument, msg).Err()
	}
	return &slackbotv1.EventResponse{Challenge: challenge}, nil
}

// the request has two layers - the "outer event" and "inner event". At this step, we are examining the outer event
// this can be 1 of three types: "url_verification" (handled seperately above), "event_callback", and "app_rate_limited"
func (m *mod) handleEvent(req *slackbotv1.EventRequest) error {
	switch req.Type {
	case slackevents.CallbackEvent:
		return m.handleCallBackEvent(req.Event)
	case slackevents.AppRateLimited:
		// this type of event is sent if our endpoint has received more than 30,000 request events in 60 minutes
		// https://api.slack.com/apis/connections/events-api#the-events-api__responding-to-events__rate-limiting
		// TODO: (sperry) do we want to handle this differently
		return errors.New("app's event subscriptions are being rate limited")
	default:
		// in the case that the Slack API adds a new "outer event" type
		return fmt.Errorf("received unexpected event type: %s", req.Type)
	}
}

// the request has two layers - the "outer event" and "inner event". At this step, we are examining the "inner event"
// the 2 event types we currently support are "app_mention" (messages that mention the bot directly) and "message" (specifically DMs with the bot)
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
