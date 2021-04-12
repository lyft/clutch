package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	"github.com/lyft/clutch/backend/service/bot/slack"
)

func TestNew(t *testing.T) {
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	svc, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	_, ok := svc.(Service)
	assert.True(t, ok)
}

func TestMatchSlackCommand(t *testing.T) {
	s := &svc{}

	testCases := []struct {
		command  string
		expected string
	}{
		{command: "<@UVWXYZ123> hello", expected: "Hello, World!"},
		{command: "hello", expected: "Hello, World!"},
		{command: "foo", expected: slack.DefaultHelp()},
	}

	for _, test := range testCases {
		match := s.MatchSlackCommand(test.command)
		assert.Equal(t, test.expected, match)
	}
}

func TestMatchCommand(t *testing.T) {
	s := &svc{}

	testCases := []struct {
		command    string
		expected   string
		expectedOk bool
	}{
		{command: "hello ", expected: "Hello, World!", expectedOk: true},
		{command: " hello", expected: "Hello, World!", expectedOk: true},
		{command: "Hello", expected: "Hello, World!", expectedOk: true},
		{command: "foo", expected: "", expectedOk: false},
	}

	for _, test := range testCases {
		match, ok := s.MatchCommand(test.command)
		assert.Equal(t, test.expected, match)
		assert.Equal(t, test.expectedOk, ok)
	}
}
