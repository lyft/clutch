package sql

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
)

func TestConvertAPIBody(t *testing.T) {
	// set up for TestConvertAPIBody
	a1 := (*anypb.Any)(nil)

	p1 := &ec2v1.Instance{InstanceId: "i-123456789abcdef0"}
	a2, _ := anypb.New(p1)

	p2 := &k8sapiv1.ResizeHPAResponse{}
	a3, _ := anypb.New(p2)

	tests := []struct {
		input *anypb.Any
	}{
		{input: nil},
		// case: input is a typed nil
		{input: a1},
		// case: input is typed with non-nil value
		{input: a2},
		// case: input is typed with non-nil value
		{input: a3},
	}

	for _, test := range tests {
		b, err := convertAPIBody(test.input)
		assert.NotNil(t, b)
		assert.NoError(t, err)
	}
}

func TestAPIBodyProto(t *testing.T) {
	var nilJSON json.RawMessage

	tests := []struct {
		input     json.RawMessage
		expectNil bool
	}{
		{input: nil, expectNil: true},
		{input: nilJSON, expectNil: true},
		{input: []byte(`{}`)},
		{input: []byte(`{"@type":"type.googleapis.com/clutch.audit.v1.RequestEvent"}`)},
	}

	for _, test := range tests {
		a, err := apiBodyProto(test.input)
		if test.expectNil {
			assert.Nil(t, a)
		} else {
			assert.NotNil(t, a)
		}
		assert.NoError(t, err)
	}
}
