package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
