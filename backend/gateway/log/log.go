package log

import (
	"encoding/json"

	"go.uber.org/zap"
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
