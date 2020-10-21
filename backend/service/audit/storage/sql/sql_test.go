package sql

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	"github.com/stretchr/testify/assert"
)

func TestConvertAPIBody(t *testing.T) {
	// set up
	p1 := (*any.Any)(nil)

	// case: input type Any that has a nil value
	b, err := convertAPIBody(p1)
	assert.Nil(t, p1)
	assert.NoError(t, err)

	// set up
	p2 := &ec2v1.Instance{InstanceId: "i-123456789abcdef0"}
	a, err := ptypes.MarshalAny(p2)
	assert.NoError(t, err)

	// case: input type Any that is a non-nil value
	b, err = convertAPIBody(a)
	assert.NotNil(t, b)
	assert.NoError(t, err)
}

func TestAPIBodyProto(t *testing.T) {
	// set up
	var nilJSON json.RawMessage

	// case: input that is a nil value
	result, err := apiBodyProto(nilJSON)
	assert.Nil(t, result)
	assert.NoError(t, err)

	// set up
	details := `{"@type":"type.googleapis.com/clutch.audit.v1.RequestEvent"}`

	// case: input is not nil value
	a, err := apiBodyProto([]byte(details))
	assert.NotNil(t, a)
	assert.NoError(t, err)
}
