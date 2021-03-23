package topology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
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
			expect: 1953460335,
		},
		{
			id:     "key with special chars",
			input:  "*()#@&!*(#!@",
			expect: 707275043,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			id := convertLockIdToAdvisoryLockId(tt.input)
			assert.Equal(t, tt.expect, id)
		})
	}
}

func TestPrepareSetCacheBulkValues(t *testing.T) {
	t.Parallel()

	topology := &client{}
	tests := []struct {
		id                string
		input             []*topologyv1.Resource
		expectQueryString string
	}{
		{
			id: "one item",
			input: []*topologyv1.Resource{
				{
					Id: "one",
					Pb: &anypb.Any{},
				},
			},
			expectQueryString: "($1, $2, $3, $4)",
		},
		{
			id: "two items",
			input: []*topologyv1.Resource{
				{
					Id: "one",
					Pb: &anypb.Any{},
				},
				{
					Id: "two",
					Pb: &anypb.Any{},
				},
			},
			expectQueryString: "($1, $2, $3, $4),($5, $6, $7, $8)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			_, queryString := topology.prepareBulkCacheInsert(tt.input)
			assert.Equal(t, tt.expectQueryString, queryString)
		})
	}
}
