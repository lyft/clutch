package topology

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/mock/service/dbmock"
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

func TestProcessTopologyObjectChannelSingleItem(t *testing.T) {
	m := dbmock.NewMockDB()
	topology := &client{
		db:               m.DB(),
		log:              zaptest.NewLogger(t),
		scope:            tally.NewTestScope("", nil),
		batchInsertSize:  1,
		batchInsertFlush: time.Millisecond * 1,
	}

	updateCacheChan := make(chan *topologyv1.UpdateCacheRequest, 1)
	updateCacheChan <- &topologyv1.UpdateCacheRequest{
		Resource: generatePrepareBulkCacheInsertInput(1)[0],
		Action:   topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
	}

	go func() {
		time.Sleep(time.Second * 1)

		// Close the channel so processTopologyObjectChannel exits
		// asserts that a single item can make it through even with a batch insert size of 10
		m.MustMeetExpectations()
		close(updateCacheChan)
	}()

	topology.processTopologyObjectChannel(context.Background(), updateCacheChan, "testing")
	m.Mock.ExpectExec("INSERT INTO topology_cache .*")
}

func TestProcessTopologyObjectChannelBatchInsert(t *testing.T) {
	m := dbmock.NewMockDB()
	topology := &client{
		db:               m.DB(),
		scope:            tally.NewTestScope("", nil),
		log:              zaptest.NewLogger(t),
		batchInsertSize:  2,
		batchInsertFlush: time.Second * 10,
	}

	updateCacheChan := make(chan *topologyv1.UpdateCacheRequest)
	go func() {
		for i := 0; i < 4; i++ {
			updateCacheChan <- &topologyv1.UpdateCacheRequest{
				Resource: generatePrepareBulkCacheInsertInput(1)[0],
				Action:   topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
			}
		}

		time.Sleep(time.Second * 1)
		// Close the channel so processTopologyObjectChannel exits
		m.MustMeetExpectations()
		close(updateCacheChan)
	}()

	topology.processTopologyObjectChannel(context.Background(), updateCacheChan, "testing")

	// We expect two batching of inserts of two items each
	m.Mock.ExpectExec("INSERT INTO topology_cache .*")
	m.Mock.ExpectExec("INSERT INTO topology_cache .*")
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

func TestConvertBatchInsertToSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id             string
		input          []*topologyv1.Resource
		expectedLength int
	}{
		{
			id:             "1",
			input:          generatePrepareBulkCacheInsertInput(1),
			expectedLength: 1,
		},
		{
			id:             "10",
			input:          generatePrepareBulkCacheInsertInput(10),
			expectedLength: 10,
		},
		{
			id:             "100",
			input:          generatePrepareBulkCacheInsertInput(100),
			expectedLength: 100,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			resources := make(map[string]*topologyv1.Resource)
			for _, v := range tt.input {
				resources[v.Id] = v
			}

			actaul := convertBatchInsertToSlice(resources)
			assert.Equal(t, tt.expectedLength, len(actaul))
		})
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
