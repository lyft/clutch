package resolver

import (
	"fmt"
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

func TestAppendAutocompleteResultsToLimit(t *testing.T) {
	results := &[]*resolverv1.AutocompleteResult{}
	limit := 50

	first := generateAutocompleteResults(25)
	second := generateAutocompleteResults(25)
	third := generateAutocompleteResults(1)

	appendAutocompleteResultsToLimit(results, first, limit)
	assert.Equal(t, 25, len(*results))

	appendAutocompleteResultsToLimit(results, second, limit)
	assert.Equal(t, 50, len(*results))

	appendAutocompleteResultsToLimit(results, third, limit)
	assert.Equal(t, 50, len(*results))
}

func generateAutocompleteResults(limit int) []*resolverv1.AutocompleteResult {
	results := []*resolverv1.AutocompleteResult{}
	for i := 0; i < limit; i++ {
		results = append(results, &resolverv1.AutocompleteResult{
			Id:    fmt.Sprint(i),
			Label: fmt.Sprint(i),
		})
	}
	return results
}
