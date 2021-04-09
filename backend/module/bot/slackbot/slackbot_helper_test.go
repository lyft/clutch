package slackbot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimUserid(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{input: "<@UVWXYZ123> describe pod", output: "describe pod"},
		// only the first user id should be removed
		{input: "<@UVWXYZ123> page oncall <@UVWABC123>", output: "page oncall <@UVWABC123>"},
		// no user id so nothing to remove
		{input: "describe pod", output: "describe pod"},
	}

	for _, test := range testCases {
		result := TrimUserId(test.input)
		assert.Equal(t, test.output, result)
	}
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
		result := TrimRedundantSpaces(test.input)
		assert.Equal(t, test.output, result)
	}
}

func TestDefaultHelp(t *testing.T) {
	expected := fmt.Sprintf("%s\n%s", HelpIntro, HelpDetails)
	result := DefaultHelp()
	assert.Equal(t, expected, result)
}
