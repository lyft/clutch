package sql

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

var protoAny = &any.Any{TypeUrl: "type.googleapis.com/clutch.audit.v1.RequestEvent"}

func TestConvertAPIBody(t *testing.T) {
	b, err := convertAPIBody(protoAny)
	assert.NotNil(t, b)
	assert.NoError(t, err)
}

func TestAPIBodyProto(t *testing.T) {
	b, err := convertAPIBody(protoAny)
	assert.NotNil(t, b)
	assert.NoError(t, err)

	a, err := apiBodyProto(b)
	assert.NotNil(t, a)
	assert.NoError(t, err)
}
