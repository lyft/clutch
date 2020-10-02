package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertLockIdToAdvisoryLockId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id     string
		input  string
		expect uint32
	}{
		{
			id:     "key with chars",
			input:  "topologycache",
			expect: 1047847070,
		},
		{
			id:     "key with space",
			input:  "topology cache",
			expect: 4086861350,
		},
		{
			id:     "key with special chars",
			input:  "*()#@&!*(#!@",
			expect: 1447352910,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			id, err := convertLockIdToAdvisoryLockId(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, id)
		})
	}
}
