package resolver

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
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
		assert.Len(t, s.Details(), 1)
		d, ok := s.Details()[0].(*apiv1.ErrorDetails)
		assert.True(t, ok)
		assert.Len(t, d.Wrapped, 1)
		assert.EqualValues(t, codes.DataLoss, d.Wrapped[0].Code)
	}

	// Some results, no failures = OK.
	{
		resp.Results = append(resp.Results, &anypb.Any{})
		resp.PartialFailures = nil
		err := resp.isError(wanted, schemas)
		assert.NoError(t, err)
	}
}

func TestAllStatusMatch(t *testing.T) {
	cases := []struct {
		name   string
		code   codes.Code
		expect bool
		status []*statuspb.Status
	}{
		{
			name:   "empty",
			code:   codes.NotFound,
			expect: false,
		},
		{
			name:   "singleMatch",
			code:   codes.Unavailable,
			status: []*statuspb.Status{{Code: int32(codes.Unavailable)}},
			expect: true,
		},
		{
			name:   "multipleMatch",
			code:   codes.Unauthenticated,
			status: []*statuspb.Status{{Code: int32(codes.Unauthenticated)}, {Code: int32(codes.Unauthenticated)}},
			expect: true,
		},
		{
			name:   "singleNoMatch",
			code:   codes.InvalidArgument,
			status: []*statuspb.Status{{Code: int32(codes.Unauthenticated)}},
			expect: false,
		},
		{
			name:   "multipleNoMatch",
			code:   codes.InvalidArgument,
			status: []*statuspb.Status{{Code: int32(codes.Unauthenticated)}, {Code: int32(codes.Unimplemented)}},
			expect: false,
		},
		{
			name:   "multipleSomeMatch",
			code:   codes.InvalidArgument,
			status: []*statuspb.Status{{Code: int32(codes.InvalidArgument)}, {Code: int32(codes.Unimplemented)}},
			expect: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expect, allStatusMatch(c.code, c.status))
		})
	}
}
