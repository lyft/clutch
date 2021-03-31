package topology

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
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

func TestPrepareBulkCacheInsert(t *testing.T) {
	t.Parallel()

	topology := &client{}
	tests := []struct {
		id                    string
		input                 []*topologyv1.Resource
		expectedQueryParams   string
		expectedQueryArgsSize int
	}{
		{
			id:                    "one item",
			input:                 generatePrepareBulkCacheInsertInput(1),
			expectedQueryParams:   "($1, $2, $3, $4)",
			expectedQueryArgsSize: 4,
		},
		{
			id:                    "two items",
			input:                 generatePrepareBulkCacheInsertInput(2),
			expectedQueryParams:   "($1, $2, $3, $4),($5, $6, $7, $8)",
			expectedQueryArgsSize: 8,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			queryArgs, queryParams := topology.prepareBulkCacheInsert(tt.input)
			assert.Equal(t, tt.expectedQueryParams, queryParams)
			assert.Equal(t, tt.expectedQueryArgsSize, len(queryArgs))
		})
	}
}

func BenchmarkPrepareBulkCacheInsert(b *testing.B) {
	tenk := generatePrepareBulkCacheInsertInput(10000)
	topology := &client{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		topology.prepareBulkCacheInsert(tenk)
	}
}

func generatePrepareBulkCacheInsertInput(count int) []*topologyv1.Resource {
	results := []*topologyv1.Resource{}

	for i := 0; i < count; i++ {
		instance := ec2v1.Instance{
			InstanceId: fmt.Sprintf("%d", i),
		}

		instancePb, err := anypb.New(&instance)
		if err != nil {
			panic(err)
		}

		results = append(results, &topologyv1.Resource{
			Id:       instance.InstanceId,
			Pb:       instancePb,
			Metadata: make(map[string]*structpb.Value),
		})
	}

	return results
}
