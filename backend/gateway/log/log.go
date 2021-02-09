package log

import (
	"encoding/json"

	"go.uber.org/zap"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtoField returns a zap field that logs as the embedded JSON representation of the Protobuf message.
func ProtoField(key string, m proto.Message) zap.Field {
	b, err := protojson.Marshal(m)
	if err != nil {
		return zap.Any(key, m)
	}
	return zap.Any(key, json.RawMessage(b))
}

// GrpcStatusField takes in an error and returns a zap field that logs as the embedded JSON
// representation of a gRPC Status message.
func GrpcStatusField(key string, err error) zap.Field {
	var m *spb.Status
	// the second return val is a bool reporting whether it was able to
	// convert the input error to a Status. If false, FromError returns a new Status with
	// codes.Unknown and the original input error message.
	// https: //github.com/grpc/grpc-go/blob/master/status/status.go#L81
	s, _ := status.FromError(err)

	m = s.Proto()

	b, err := protojson.Marshal(m)
	if err != nil {
		return zap.Any(key, m)
	}

	return zap.Any(key, json.RawMessage(b))
}
