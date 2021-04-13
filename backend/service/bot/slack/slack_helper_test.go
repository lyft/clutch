package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimSlackUserId(t *testing.T) {
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
		result := TrimSlackUserId(test.input)
		assert.Equal(t, test.output, result)
	}
}
