package resolver

import (
	"testing"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
)

// Check that our TypeURL function matches proto's interpretation of the URL.
func TestTypeURL(t *testing.T) {
	u := TypeURL((*resolverv1.Schema)(nil))

	s := &resolverv1.Schema{}
	a, _ := ptypes.MarshalAny(s)
	assert.Equal(t, u, a.TypeUrl)
	assert.Equal(t, u, "type.googleapis.com/clutch.resolver.v1.Schema")
}

func TestMarshalProtoSliceToAny(t *testing.T) {
	// nil
	v, err := MarshalProtoSlice(nil)
	assert.NoError(t, err)
	assert.Nil(t, v)

	// non-slice value
	_, err = MarshalProtoSlice("foo")
	assert.Error(t, err)

	// slice of strings
	_, err = MarshalProtoSlice([]string{"foo", "bar"})
	assert.Error(t, err)

	// slice of protos
	pbs := []*ec2v1.Instance{
		{
			InstanceId: "123",
		},
		{
			InstanceId: "456",
		},
	}
	v, err = MarshalProtoSlice(pbs)
	assert.NoError(t, err)
	assert.Equal(t, len(pbs), len(v))
	var instance ec2v1.Instance
	err = ptypes.UnmarshalAny(v[0], &instance)
	assert.NoError(t, err)
	assert.Equal(t, "123", instance.InstanceId)
	err = ptypes.UnmarshalAny(v[1], &instance)
	assert.NoError(t, err)
	assert.Equal(t, "456", instance.InstanceId)
}
