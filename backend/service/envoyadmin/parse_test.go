package envoyadmin

import (
	"testing"

	envoytriagev1 "github.com/lyft/clutch/backend/api/envoytriage/v1"
	"github.com/stretchr/testify/assert"
)

func TestNodeMetadataFromResponse(t *testing.T) {
	testCases := []struct {
		resp     []byte
		metadata *envoytriagev1.NodeMetadata
		hasErr   bool
	}{
		// resp has no command line options
		{
			resp:     []byte("{}"),
			metadata: &envoytriagev1.NodeMetadata{},
			hasErr:   false,
		},
	}

	for _, test := range testCases {
		metadata, err := nodeMetadataFromResponse(test.resp)
		if test.hasErr {
			assert.Error(t, err)
			assert.Nil(t, metadata)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, metadata)
			assert.Equal(t, test.metadata, metadata)
		}
	}
}
