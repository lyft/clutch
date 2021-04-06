package slackbot

import (
	"context"
	"errors"
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
	hSignature = "X-Slack-Signature"
	hTimestamp = "X-Slack-Request-Timestamp"
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
	m.logger.Info("received event from Slack API")

	// TODO: (sperry) for this flow to work, need to add custom header matchers b/c by default grpc drops
	// "X-Slack-Signature" and "X-Slack-Request-Timestamp" due to not meeting the preconditions here
	// https://github.com/grpc-ecosystem/grpc-gateway/blob/519e7f45f48e3e378e3341d076579f3ce52aa717/runtime/mux.go#L63-L74
	err := verifySlackRequest(ctx, req, m.signingSecret)
	if err != nil {
		m.logger.Error("verify slack request error", log.ErrorField(err))
		return nil, status.New(codes.Unauthenticated, err.Error()).Err()
	}

	// recieved the first time when configuring an Events API Request URL, https://api.slack.com/events/url_verification
	// the Slack API will send us the challenge string in the request and we must respond back with the same challenge
	if req.Type == slackevents.URLVerification {
		challenge := req.Challenge
		// this should never happen
		if challenge == "" {
			msg := "slack API provided an empty challenge string"
			m.logger.Error(msg)
			return nil, status.New(codes.InvalidArgument, msg).Err()
		}
		return &slackbotv1.EventResponse{Challenge: challenge}, nil
	}


	return &slackbotv1.EventResponse{}, nil
}

// retrieves the signature and timestamp from the headers which will be used to verify the request was from Slack
func getSlackRequestHeaders(ctx context.Context) (http.Header, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no headers present on request")
	}

	signature := md.Get("x-slack-signature")
	if len(signature) == 0 {
		return nil, errors.New("x-slack-signature header not present on request")
	}

	timestamp := md.Get("x-slack-request-timestamp")
	if len(timestamp) == 0 {
		return nil, errors.New("x-slack-request-timestamp header not present on request")
	}

	return http.Header{
		hSignature: signature,
		hTimestamp: timestamp,
	}, nil
}

// verifies the request came from Slack
// https://api.slack.com/authentication/verifying-requests-from-slack#verifying-requests-from-slack-using-signing-secrets__a-recipe-for-security__step-by-step-walk-through-for-validating-a-request
func verifySlackRequest(ctx context.Context, req *slackbotv1.EventRequest, signingSecret string) error {
	// returns the signature and timestamp from the requuest headers
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

