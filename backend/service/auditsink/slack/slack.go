package slack

// <!-- START clutchdoc -->
// description: Posts events to a configured Slack workspace and channel.
// <!-- END clutchdoc -->

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/slack-go/slack"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
	"github.com/lyft/clutch/backend/gateway/log"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/auditsink"
)

const Name = "clutch.service.auditsink.slack"

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &configv1.SlackConfig{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}

	s := &svc{
		logger:   logger,
		filter:   config.Filter,
		override: config.Override,
		scope:    scope,
		slack:    slack.New(config.Token),
		channel:  config.Channel,
	}
	return s, nil
}

type svc struct {
	logger *zap.Logger
	scope  tally.Scope

	filter   *auditconfigv1.Filter
	override *auditconfigv1.Override

	slack   *slack.Client
	channel string
}

func (s *svc) Write(event *auditv1.Event) error {
	if !auditsink.Filter(s.filter, event) {
		return nil
	}

	switch event.GetEventType().(type) {
	case *auditv1.Event_Event:
		return s.writeRequestEvent(event.GetEvent())
	default:
		return nil
	}
}

func (s *svc) writeRequestEvent(event *auditv1.RequestEvent) error {
	// Get user ID for pretty message printing.
	user, err := s.slack.GetUserByEmail(event.Username)

	var username string
	if err != nil {
		s.logger.Warn(
			"failure to get user information from slack",
			zap.String("username", event.Username),
			zap.Error(err),
		)
		username = event.Username
	} else {
		username = fmt.Sprintf("<@%s>", user.ID)
	}

	// TODO: clean up this if/else flow
	var messageText string
	defaultText := formatText(username, event)
	customTemplate, ok := auditsink.GetCustomSlackAudit(s.override, event)
	if !ok {
		// custom slack audit wasn't specified
		messageText = defaultText
	} else {
		customText, err := FormatCustomText(customTemplate, event)
		if err != nil {
			// use the default text as the fallback
			messageText = defaultText
			s.logger.Error("Unable to parse custom slack audit template", log.ErrorField(err))
		} else {
			// append the default and custom message to the final message text
			messageText = fmt.Sprintf("%s\n%s", defaultText, customText)
		}
	}

	// Post
	if _, _, err := s.slack.PostMessage(s.channel, slack.MsgOptionText(messageText, false)); err != nil {
		return err
	}
	return nil
}

func formatText(username string, event *auditv1.RequestEvent) string {
	// Make text: `user` performed `/path/to/action` via `service` using Clutch on resources: ...
	const messageFormat = "`%s` performed `%s` via `%s` using Clutch on resource(s):"
	messageText := fmt.Sprintf(
		messageFormat,
		username,
		event.MethodName,
		event.ServiceName,
	)

	// - resourceName (`resourceTypeUrl`)
	const resourceFormat = "\n- %s (`%s`)"
	for _, resource := range event.Resources {
		messageText += fmt.Sprintf(resourceFormat, resource.Id, resource.TypeUrl)
	}

	return messageText
}

// FormatCustomText performs variable substitution in the custom slack template using the audit event metadata
func FormatCustomText(message string, event *auditv1.RequestEvent) (string, error) {
	tmpl, err := template.New("customText").Parse(message)
	if err != nil {
		return "", err
	}

	// TODO: handle custom slack templates that request both req and resp metadata
	reqMetadata := event.RequestMetadata.Body
	// UnmarshalNew returns the message of the specified type based on the TypeURL
	m, err := reqMetadata.UnmarshalNew()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m); err != nil {
		return "", err
	}

	return buf.String(), nil
}
