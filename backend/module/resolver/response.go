package resolver

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lyft/clutch/backend/resolver"
)

func newResponse() *response {
	return &response{
		Results:         []*any.Any{},
		PartialFailures: []*statuspb.Status{},
	}
}

// Generic object to handle common operations for SearchResponse and ResolveResponse.
type response struct {
	Results         []*any.Any
	PartialFailures []*statuspb.Status
}

// Get rid of any extraneous results or partial failures based on the limit provided in the request.
func (r *response) truncate(limit uint32) {
	if limit == 0 {
		return
	}

	if len(r.Results) > int(limit) {
		// Truncate in case it wasn't done earlier.
		r.Results = r.Results[:limit]
	}

	if len(r.Results) == int(limit) {
		// If we fulfilled our limit then errors are not relevant.
		r.PartialFailures = nil
	}
}

func (r *response) marshalResults(results *resolver.Results) error {
	for _, result := range results.Messages {
		asAny, err := ptypes.MarshalAny(result)
		if err != nil {
			return err
		}
		r.Results = append(r.Results, asAny)
	}

	for _, failure := range results.PartialFailures {
		r.PartialFailures = append(r.PartialFailures, failure.Proto())
	}

	return nil
}

func (r *response) isError(wanted string, searchedSchemas []string) error {
	if len(r.Results) > 0 {
		return nil
	}

	msg := fmt.Sprintf("Did not find a '%s', looked by: %s", wanted, strings.Join(searchedSchemas, ", "))

	if len(r.PartialFailures) > 0 {
		s := status.New(codes.FailedPrecondition, msg)
		s, _ = s.WithDetails(resolver.MessageSlice(r.PartialFailures)...)
		return s.Err()
	}

	return status.Error(codes.NotFound, msg)
}
