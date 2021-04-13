package bot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	botv1 "github.com/lyft/clutch/backend/api/config/service/bot/v1"
)

func TestSanatize(t *testing.T) {
	testCases := []struct {
		botProvider botv1.Bot
		text        string
		expected    string
		expectedErr bool
	}{
		{botProvider: botv1.Bot_SLACK, text: "hello", expected: "hello"},
		{botProvider: botv1.Bot_UNSPECIFIED, text: "hello", expected: "", expectedErr: true},
		{botProvider: botv1.Bot_SLACK, text: "<@UVWXYZ123> describe  Pod", expected: "describe pod"},
	}

	for _, test := range testCases {
		result, err := sanatize(test.botProvider, test.text)
		if test.expectedErr {
			assert.Error(t, err)
			assert.Empty(t, result)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}

func TestDefaultHelp(t *testing.T) {
	expected := fmt.Sprintf("%s\n%s", HelpIntro, HelpDetails)
	result := defaultHelp()
	assert.Equal(t, expected, result)
}

func TestTrimRedundantSpaces(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		// extra whitespace in the beginning
		{input: " hello world", output: "hello world"},
		// extra whitesapce in the middle
		{input: "hello  world", output: "hello world"},
		// extra whitespace at the end
		{input: "hello world  ", output: "hello world"},
		// no extra whitespace
		{input: "hello world", output: "hello world"},
	}

	for _, test := range testCases {
		result := trimRedundantSpaces(test.input)
		assert.Equal(t, test.output, result)
	}
}
