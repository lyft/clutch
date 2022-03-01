package shortlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortlinkHash(t *testing.T) {
	tests := []struct {
		id           string
		input        int
		expectLength int
		shouldError  bool
	}{
		{
			id:           "len of 10",
			input:        10,
			expectLength: 10,
			shouldError:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			hash, err := generateShortlinkHash(test.input)
			if test.shouldError {
				assert.Error(t, err)
				assert.Nil(t, hash)
			} else {
				assert.NoError(t, err)
				assert.Len(t, hash, test.expectLength)
			}
		})
	}
}
