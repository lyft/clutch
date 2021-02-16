package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	resolverv1 "github.com/lyft/clutch/backend/api/resolver/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/mock/service/awsmock"
	"github.com/lyft/clutch/backend/mock/service/topologymock"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/topology"
)

func TestNewAwsResolver(t *testing.T) {
	service.Registry["clutch.service.aws"] = awsmock.New()
	service.Registry["clutch.service.topology"] = topologymock.New()

	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	aws, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, aws)

	// Test successful construction without the topology service
	delete(service.Registry, "clutch.service.topology")
	awsNoTopology, err2 := New(nil, log, scope)
	assert.NoError(t, err2)
	assert.NotNil(t, awsNoTopology)
}

func TestAutoCompleteErrorHandling(t *testing.T) {
	service.Registry["clutch.service.aws"] = awsmock.New()
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)

	aws, err := New(nil, log, scope)
	assert.NoError(t, err)
	assert.NotNil(t, aws)

	// Test error handling for topology service not found
	_, err2 := aws.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.Error(t, err2)

	// Test error handling for a topology search failure
	awsResolver := res{
		topology: &mockTopologySearch{
			autoCompleteError: fmt.Errorf("error"),
		},
	}
	_, err3 := awsResolver.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.Error(t, err3)
}

func TestAutoCompleteResults(t *testing.T) {
	awsResolver := res{
		topology: &mockTopologySearch{
			autoCompleteResults: []*topologyv1.Resource{
				{
					Id: "meow",
				},
				{
					Id: "cat",
				},
				{
					Id: "yawn",
				},
			},
		},
	}

	expect := []*resolverv1.AutocompleteResult{
		{
			Id: "meow",
		},
		{
			Id: "cat",
		},
		{
			Id: "yawn",
		},
	}

	results, err := awsResolver.Autocomplete(context.Background(), "type_url", "search", 0)
	assert.NoError(t, err)
	assert.Equal(t, expect, results)
}

type mockTopologySearch struct {
	topology.Service
	autoCompleteError   error
	autoCompleteResults []*topologyv1.Resource
}

func (m *mockTopologySearch) Search(ctx context.Context, search *topologyv1.SearchRequest) ([]*topologyv1.Resource, string, error) {
	if m.autoCompleteError != nil {
		return nil, "", m.autoCompleteError
	}

	return m.autoCompleteResults, "0", nil
}
