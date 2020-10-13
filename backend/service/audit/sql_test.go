package audit

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

var protoAny = &any.Any{TypeUrl: "type.googleapis.com/clutch.audit.v1.RequestEvent"}

func TestConvertAPIBody(t *testing.T) {
	b := convertAPIBody(protoAny)
	assert.NotNil(t, b)
}

func TestAPIBodyProto(t *testing.T) {
	b := convertAPIBody(protoAny)
	assert.NotNil(t, b)

	a := apiBodyProto(b)
	assert.NotNil(t, a)
}
