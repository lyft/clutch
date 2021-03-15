// Package statuserr provides helper functions for dealing with errors and gRPC status objects.
package statuserr

import (
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func AllProtoCodesMatch(c codes.Code, spb []*statuspb.Status) bool {
	if len(spb) == 0 {
		return false
	}

	for _, s := range spb {
		if s.Code != int32(c) {
			return false
		}
	}
	return true
}
