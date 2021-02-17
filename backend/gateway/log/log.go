package log

import (
	"encoding/json"

	"go.uber.org/zap"
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

// NamedErrorField returns a zap field that logs as the embedded JSON representation of the error.
// If it's a gRPC error, it converts it to a Status and returns
// the unpacked information. Otherwise, it returns the original error.
func NamedErrorField(key string, err error) zap.Field {
	s, ok := status.FromError(err)
	if !ok {
		return zap.NamedError(key, err)
	}
	return ProtoField(key, s.Proto())
}

// ErrorField calls NamedErrorField with the key "error"
func ErrorField(err error) zap.Field {
	return NamedErrorField("error", err)
}
