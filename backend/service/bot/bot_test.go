package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	botv1 "github.com/lyft/clutch/backend/api/config/service/bot/v1"
)

func TestNew(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	cfg, _ := anypb.New(&botv1.Config{BotProvider: botv1.Bot_SLACK})

	svc, err := New(cfg, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	_, ok := svc.(Service)
	assert.True(t, ok)
}

func TestGetBotProvider(t *testing.T) {
	valid := &botv1.Config{BotProvider: botv1.Bot_SLACK}
	b, err := getBotProvider(valid)
	assert.NoError(t, err)
	assert.Equal(t, botv1.Bot_SLACK, b)

	invalid := &botv1.Config{BotProvider: botv1.Bot_UNSPECIFIED}
	b, err = getBotProvider(invalid)
	assert.Error(t, err)
	assert.Zero(t, b)
}

func TestMatchCommand(t *testing.T) {
	testCases := []struct {
		svc      *svc
		command  string
		expected string
	}{
		{svc: &svc{botProvider: botv1.Bot_SLACK}, command: "hello", expected: "Hello, World!"},
		{svc: &svc{botProvider: botv1.Bot_SLACK}, command: "foo", expected: defaultHelp()},
		{svc: &svc{}, command: "hello", expected: "bot type 'UNSPECIFIED' not implemented"},
	}

	for _, test := range testCases {
		match := test.svc.MatchCommand(test.command)
		assert.Equal(t, test.expected, match)
	}
}
