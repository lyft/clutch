package resolver

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	"github.com/lyft/clutch/backend/resolver"
)

// Simple truncate tests with 0 or 1 or 2 objects.
func TestTruncateSimple(t *testing.T) {
	resp := newResponse()
	resp.truncate(0)
	assert.Len(t, resp.Results, 0)
	assert.Len(t, resp.PartialFailures, 0)

	resp.truncate(1)
	assert.Len(t, resp.Results, 0)
	assert.Len(t, resp.PartialFailures, 0)

	resp.Results = []*anypb.Any{{}}
	resp.truncate(0)
	assert.Len(t, resp.Results, 1)

	resp.truncate(1)
	assert.Len(t, resp.Results, 1)

	resp.truncate(2)
	assert.Len(t, resp.Results, 1)

	resp.Results = []*anypb.Any{{}, {}}
	resp.truncate(1)
	assert.Len(t, resp.Results, 1)
}

func TestTruncateLimitMetRemovesFailures(t *testing.T) {
	resp := newResponse()
	resp.Results = []*anypb.Any{{}, {}}
	resp.PartialFailures = []*statuspb.Status{{}, {}, {}}

	resp.truncate(2)
	assert.Len(t, resp.PartialFailures, 0)
}

func TestMarshalResults(t *testing.T) {
	resp := newResponse()
	err := resp.marshalResults(&resolver.Results{
		Messages:        []proto.Message{&healthcheckv1.HealthcheckRequest{}, &healthcheckv1.HealthcheckResponse{}},
		PartialFailures: []*status.Status{status.New(codes.Aborted, ""), status.New(codes.DataLoss, "")},
	})

	assert.NoError(t, err)

	assert.Len(t, resp.Results, 2)
	assert.Equal(t, resp.Results[0].TypeUrl, "type.googleapis.com/clutch.healthcheck.v1.HealthcheckRequest")
	assert.Equal(t, resp.Results[1].TypeUrl, "type.googleapis.com/clutch.healthcheck.v1.HealthcheckResponse")

	assert.Len(t, resp.PartialFailures, 2)
	assert.EqualValues(t, resp.PartialFailures[0].Code, codes.Aborted)
	assert.EqualValues(t, resp.PartialFailures[1].Code, codes.DataLoss)
}

func TestIsError(t *testing.T) {
	wanted := "type.googleapis.com/clutch.foo"
	schemas := []string{"foo", "bar"}

	resp := newResponse()

	// No results, no failures = NotFound.
	{
		err := resp.isError(wanted, schemas)
		assert.Error(t, err)
		s, _ := status.FromError(err)
		assert.Equal(t, codes.NotFound, s.Code())
	}

	// No results, some failures = FailedPrecondition.
	{
		resp.PartialFailures = append(resp.PartialFailures, status.New(codes.DataLoss, "").Proto())
		err := resp.isError(wanted, schemas)
		assert.Error(t, err)
		s, _ := status.FromError(err)
		assert.Equal(t, codes.FailedPrecondition, s.Code())
	}

	// Some results, no failures = OK.
	{
		resp.Results = append(resp.Results, &anypb.Any{})
		resp.PartialFailures = nil
		err := resp.isError(wanted, schemas)
		assert.NoError(t, err)
	}
}
