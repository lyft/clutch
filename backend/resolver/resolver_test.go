package resolver

import (
	"testing"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	k8sv1resolver "github.com/lyft/clutch/backend/api/resolver/k8s/v1"
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

func TestInputsToSchema(t *testing.T) {
	tp := "type.googleapis.com/foo.v1.Bar"
	m, err := InputsToSchemas(map[string][]descriptor.Message{
		tp: {
			(*k8sv1resolver.PodID)(nil),
		},
	})
	assert.NoError(t, err)
	assert.Len(t, m, 1)
	assert.Len(t, m[tp], 1)
	assert.Equal(t, TypeURL((*k8sv1resolver.PodID)(nil)), m[tp][0].TypeUrl)
	assert.NotEmpty(t, m[tp][0].Metadata.DisplayName)
	assert.NotEmpty(t, m[tp][0].Fields)
}

func TestHydrateDynamicOptions(t *testing.T) {
	schema := &resolverv1.Schema{
		TypeUrl: "aaa",
		Metadata: &resolverv1.SchemaMetadata{
			DisplayName: "AAA",
			Searchable:  false,
		},
		Fields: []*resolverv1.Field{
			{
				Name: "myOptions",
				Metadata: &resolverv1.FieldMetadata{
					DisplayName: "My Options",
					Type: &resolverv1.FieldMetadata_OptionField{
						OptionField: &resolverv1.OptionField{
							IncludeDynamicOptions: []string{"foo"},
							Options:               nil,
						},
					},
				},
			},
		},
	}

	m := TypeURLToSchemasMap{
		"bar": []*resolverv1.Schema{schema},
	}

	HydrateDynamicOptions(m, map[string][]*resolverv1.Option{
		"foo": {
			{
				DisplayName: "Option 1",
				Value:       &resolverv1.Option_StringValue{StringValue: "option_1"},
			},
			{
				DisplayName: "Option 2",
				Value:       &resolverv1.Option_StringValue{StringValue: "option_2"},
			},
		},
	})

	assert.Len(t, m["bar"][0].Fields[0].Metadata.GetOptionField().Options, 2)
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
