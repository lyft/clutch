package slack

// <!-- START clutchdoc -->
// description: Posts events to a configured Slack workspace and channel.
// <!-- END clutchdoc -->

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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

const (
	Name        = "clutch.service.auditsink.slack"
	defaultNone = "None"
)

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	config := &configv1.SlackConfig{}
	if err := ptypes.UnmarshalAny(cfg, config); err != nil {
		return nil, err
	}

	s := &svc{
		logger:    logger,
		filter:    config.Filter,
		overrides: config.Overrides,
		scope:     scope,
		slack:     slack.New(config.Token),
		channel:   config.Channel,
	}
	return s, nil
}

// The struct is used in a Go template and the fields need to be exported
type auditTemplateData struct {
	Request  map[string]interface{}
	Response map[string]interface{}
}

// TODO: (sperry) expand on helper funcs
// helper functions to format values in the Go template
var funcMap = template.FuncMap{
	// for inputs that are type slice/map, returns a formatted slack list
	// currently only iterates 1-level deep
	"slackList": slackList,
}

type svc struct {
	logger *zap.Logger
	scope  tally.Scope

	filter    *auditconfigv1.Filter
	overrides []*configv1.CustomMessage

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

// FormatCustomText applies the audit event metadata to the custom slack message
func FormatCustomText(message string, event *auditv1.RequestEvent) (string, error) {
	tmpl, err := template.New("customText").Funcs(funcMap).Parse(message)
	if err != nil {
		return "", err
	}

	data, err := getAuditTemplateData(event)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// When a value is nil, the Go Template sets the value as "<no value>".
	// Ex of when we hit this scenario: https://github.com/lyft/clutch/blob/main/api/k8s/v1/k8s.proto#L860
	sanitized := strings.ReplaceAll(buf.String(), "<no value>", defaultNone)

	return sanitized, nil
}

// returns the API request/response details saved in an audit event
func getAuditTemplateData(event *auditv1.RequestEvent) (*auditTemplateData, error) {
	// proto -> json
	reqJSON, err := protojson.Marshal(event.RequestMetadata.Body)
	if err != nil {
		return nil, err
	}
	respJSON, err := protojson.Marshal(event.ResponseMetadata.Body)
	if err != nil {
		return nil, err
	}

	var requestMetadata map[string]interface{}
	var responseMetadata map[string]interface{}
	// json -> map[string]interface{}
	err = json.Unmarshal(reqJSON, &requestMetadata)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respJSON, &responseMetadata)
	if err != nil {
		return nil, err
	}

	return &auditTemplateData{
		Request:  requestMetadata,
		Response: responseMetadata,
	}, nil
}

// for inputs that are type slice/map, returns a formatted slack list
func slackList(data interface{}) string {
	if data == nil {
		return defaultNone
	}

	var b strings.Builder
	const sliceItemFormat = "\n- %v"
	const mapItemFormat = "\n- %v: %v"

	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		if s.Len() == 0 {
			return defaultNone
		}
		// example:
		// - foo
		// - bar
		for i := 0; i < s.Len(); i++ {
			fmt.Fprintf(&b, sliceItemFormat, s.Index(i).Interface())
		}
		return b.String()
	case reflect.Map:
		m := reflect.ValueOf(data)
		if m.Len() == 0 {
			return defaultNone
		}
		// example:
		// - foo: value
		// - bar: value
		for _, k := range m.MapKeys() {
			v := m.MapIndex(k)
			fmt.Fprintf(&b, mapItemFormat, k.Interface(), v.Interface())
		}
		return b.String()
	}
	// We don't return an error so that the custom slack message is still returned
	// in this way, the user will know they didn't pass the correct input type to slackList
	return "ERR_INPUT_NOT_SLICE_OR_MAP"
}
