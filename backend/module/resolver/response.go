package resolver

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
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

func allStatusMatch(c codes.Code, sl []*statuspb.Status) bool {
	for _, s := range sl {
		if s.Code != int32(c) {
			return false
		}
	}
	return true
}

func (r *response) isError(wanted string, searchedSchemas []string) error {
	// If results, errors will be returned as partial failures.
	if len(r.Results) > 0 {
		return nil
	}

	wanted = strings.TrimPrefix(wanted, resolver.TypePrefix)

	msg := fmt.Sprintf("search for '%s' returned no results", wanted)

	if len(r.PartialFailures) > 0 {
		// Determine code.
		code := codes.NotFound
		if !allStatusMatch(codes.NotFound, r.PartialFailures) {
			code = codes.FailedPrecondition
			msg = fmt.Sprintf("one or more errors were encountered searching for '%s'", wanted)
		}

		s := status.New(code, msg)
		s, _ = s.WithDetails(&apiv1.ErrorMetadata{
			Metadata: map[string]*structpb.Value{
				"searchedSchemas": stringSliceAsValue(searchedSchemas),
			},
			Wrapped: r.PartialFailures,
		})

		return s.Err()
	}

	return status.Error(codes.NotFound, msg)
}

func stringSliceAsValue(v []string) *structpb.Value {
	iface := make([]interface{}, len(v))
	for idx, vv := range v {
		iface[idx] = vv
	}
	l, _ := structpb.NewValue(iface)
	return l
}
