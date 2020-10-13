package audit

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestConvertAPIBody(t *testing.T) {
	b := convertAPIBody(&any.Any{})
	assert.NotNil(t, b)
}

func TestAPIBodyProto(t *testing.T) {
	proto := &any.Any{TypeUrl: "type.googleapis.com/clutch.audit.v1.RequestEvent"}
	b := convertAPIBody(proto)
	assert.NotNil(t, b)

	a := apiBodyProto(b)
	assert.NotNil(t, a)
}
