package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
)

func TestUpdateSchemaError(t *testing.T) {
	empty := &resolverv1.Schema{}
	updateSchemaError(empty)
	assert.Nil(t, empty.Error)

	optsSchema := &resolverv1.Schema{
		Fields: []*resolverv1.Field{
			{
				Name: "foo",
				Metadata: &resolverv1.FieldMetadata{Type: &resolverv1.FieldMetadata_OptionField{
					OptionField: &resolverv1.OptionField{Options: []*resolverv1.Option{}},
				}},
			},
		},
	}
	updateSchemaError(optsSchema)
	assert.Nil(t, optsSchema.Error)

	optsSchema.Fields[0].Metadata.Required = true
	updateSchemaError(optsSchema)
	assert.NotNil(t, optsSchema.Error)
	assert.Contains(t, optsSchema.Error.Message, "missing required options")
}

func TestShouldAutoCompleteBeEnabled(t *testing.T) {
	no := shouldAutoCompleteBeEnabled()
	assert.False(t, no)

	service.Registry["clutch.service.topology"] = topologymock.New()
	yes := shouldAutoCompleteBeEnabled()
	assert.True(t, yes)
}
