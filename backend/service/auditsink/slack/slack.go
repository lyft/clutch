package slack

// <!-- START clutchdoc -->
// description: Posts events to a configured Slack workspace and channel.
// <!-- END clutchdoc -->

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"text/template"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/slack-go/slack"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"

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
	Request  map[string]interface{}
	Response map[string]interface{}
}

// helper functions to format values in the Go template
var funcMap = template.FuncMap{
	// for inputs that are type slice/map, returns a formatted slack list OR `N/A` if the slice/map is empty
	// currently only iterates 1-level deep
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
	tmpl, err := template.New("customText").Funcs(funcMap).Parse(message)
	if err != nil {
		return "", err
	}

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

// returns the API request/response details saved in an audit event
func getAuditMetadata(event *auditv1.RequestEvent) (*auditEventMetadata, error) {
	var requestMetadata map[string]interface{}
	var responseMetadata map[string]interface{}

	// proto -> json
	reqJSON, err := protojson.Marshal(event.RequestMetadata.Body)
	if err != nil {
		return nil, err
	}
	respJSON, err := protojson.Marshal(event.ResponseMetadata.Body)
	if err != nil {
		return nil, err
	}

	// json -> map[string]interface{}
	err = json.Unmarshal(reqJSON, &requestMetadata)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respJSON, &responseMetadata)
	if err != nil {
		return nil, err
	}

	return &auditEventMetadata{
		Request:  requestMetadata,
		Response: responseMetadata,
	}, nil
}

// for inputs that are type slice/map, returns a formatted slack list OR `N/A` if the slice/map is empty
func slackList(data interface{}) (string, error) {
	var text string
	const listItemFormat = "\n- %v"
	const mapItemFormat = "\n- %v: %v"

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		if s.Len() == 0 {
			return "`N/A`", nil
		}
		// example:
		// - foo
		// - bar
		for i := 0; i < s.Len(); i++ {
			text += fmt.Sprintf(listItemFormat, s.Index(i).Interface())
		}
		return text, nil
	case reflect.Map:
		m := reflect.ValueOf(data)
		if m.Len() == 0 {
			return "`N/A`", nil
		}
		// example:
		// foo: value
		// bar: value
		for _, k := range m.MapKeys() {
			v := m.MapIndex(k)
			text += fmt.Sprintf(mapItemFormat, k.Interface(), v.Interface())
		}
		return text, nil
	}
	return "", errors.New("input type is not Slice or Map")
}
