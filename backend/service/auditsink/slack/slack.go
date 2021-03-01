package slack

// <!-- START clutchdoc -->
// description: Posts events to a configured Slack workspace and channel.
// <!-- END clutchdoc -->

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/slack-go/slack"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	auditconfigv1 "github.com/lyft/clutch/backend/api/config/service/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
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
		logger:  logger,
		filter:  config.Filter,
		scope:   scope,
		slack:   slack.New(config.Token),
		channel: config.Channel,
	}
	return s, nil
}

// The struct is used in a Go template and the fields need to be exported
type auditEventMetadata struct {
	Request  proto.Message
	Response proto.Message
}

// helper functions to format values in the Go template
var funcMap = template.FuncMap{
	// for values that are type List or Map, returns a formatted slack list
	"slackList": slackList,
}

type svc struct {
	logger *zap.Logger
	scope  tally.Scope

	filter *auditconfigv1.Filter

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

	messageText := formatText(username, event)

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

// FormatCustomText parses out the audit event metdata request for the custom slack message text
func FormatCustomText(message string, event *auditv1.RequestEvent) (string, error) {
	tmpl := template.Must(template.New("customText").Funcs(funcMap).Parse(message))

	data, err := getAuditMetadata(event)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// returns the API request/response details contained in a audit event
func getAuditMetadata(event *auditv1.RequestEvent) (*auditEventMetadata, error) {
	requestMetadata := event.RequestMetadata.Body
	// UnmarshalNew returns the message of the specified type based on the TypeURL
	requestProto, err := requestMetadata.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	responseMetdata := event.ResponseMetadata.Body
	responseProto, err := responseMetdata.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	return &auditEventMetadata{
		Request:  requestProto,
		Response: responseProto,
	}, nil
}

func slackList(data interface{}, name string) (string, error) {
	pb, ok := data.(proto.Message)
	if !ok {
		// not the type we can process
		return "", errors.New("could not use input as proto message")
	}
	m := pb.ProtoReflect()

	fd := m.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		// didn't find field by name
		return "", errors.New("could not find field by name")
	}

	v := m.Get(fd)

	if fd.IsList() {
		return resolveSlice(v.List()), nil
	}

	if fd.IsMap() {
		return resolveMap(v.Map()), nil
	}

	return "", errors.New("input is not type List or Map")
}

func resolveSlice(list protoreflect.List) string {
	var text string
	const listFormat = "\n%v"

	for i := 0; i < list.Len(); i++ {
		v := list.Get(i)
		text += fmt.Sprintf(listFormat, v.Interface())
	}
	return text
}

func resolveMap(mapValue protoreflect.Map) string {
	var text string
	const mapFormat = "\n%v: %v"

	mapValue.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		var value interface{}
		switch v.Interface().(type) {
		case protoreflect.Message:
			value = v.Message().Interface()
		default:
			value = v.Interface()
		}

		text += fmt.Sprintf(mapFormat, k.Interface(), value)
		return true
	})
	return text
}
