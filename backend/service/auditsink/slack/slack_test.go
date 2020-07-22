package slack

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/service/auditsink/slack/v1"
	"github.com/lyft/clutch/backend/service/auditsink"
)

func TestNew(t *testing.T) {
	t.Parallel()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	_, err := New(nil, log, scope)
	assert.Error(t, err)

	cfg, _ := ptypes.MarshalAny(&configv1.SlackConfig{})
	svc, err := New(cfg, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	_, ok := svc.(auditsink.Sink)
	assert.True(t, ok)
}

func TestFormat(t *testing.T) {
	t.Parallel()

	username := "foo"
	event := &auditv1.RequestEvent{
		ServiceName: "service",
		MethodName:  "method",
		Type:        apiv1.ActionType_READ,
		Resources: []*auditv1.Resource{{
			TypeUrl: "clutch.aws.v1.Instance",
			Id:      "i-01234567890abcdef0",
		}},
	}
	expected := "`foo` performed `method` via `service` using Clutch on resource(s):\n" +
		"- i-01234567890abcdef0 (`clutch.aws.v1.Instance`)"

	actual := formatText(username, event)
	assert.Equal(t, expected, actual)
}
